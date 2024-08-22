package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/igm/sockjs-go/v3/sockjs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/urfave/cli/v2"
	authv1 "k8s.io/api/authorization/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
	"k8s.io/kubectl/pkg/scheme"
)

//go:generate yarn --cwd ./frontend install
//go:generate yarn --cwd ./frontend run build
//go:embed all:frontend/dist
var files embed.FS

var homeDir = homedir.HomeDir()

type FrontendRequest struct {
	Namespace       string          `json:"namespace"`
	Name            string          `json:"name"`
	KubeGVRK        KubeGVRK        `json:"kubeGVRK"`
	KubeGVRKList    []KubeGVRK      `json:"kubeGVRKList"`
	Data            string          `json:"data"`
	ResourceOptions ResourceOptions `json:"kubeOptions"`
	HelmOptions     HelmOptions     `json:"helmOptions"`
}

type KubeGVRK struct {
	Group        string `json:"group"`
	Version      string `json:"version"`
	Resource     string `json:"resource"`
	Kind         string `json:"kind"`
	IsNamespaced bool   `json:"isNamespaced"`
}

type HelmOptions struct {
	ChartName    string `json:"chartName"`
	ChartVersion string `json:"chartVersion"`
	EnvPath      string `json:"envPath"`
	RepoURL      string `json:"repoURL"`
	DryRun       bool   `json:"dryRun"`
	IsOCI        bool   `json:"isOCI"`
	ReuseValues  bool   `json:"reuseValues"`
}

type ResourceOptions struct {
	FieldSelector  string `json:"fieldSelector"`
	LabelSelector  string `json:"labelSelector"`
	TimeoutSeconds int64  `json:"timeoutSeconds"`
	Limit          int64  `json:"limit"`
	Continue       string `json:"continue"`
}

type APIResource struct {
	CliContext    *cli.Context
	AuthRequest   *AuthRequest
	AuthState     bool
	SSAR          *authv1.SelfSubjectAccessReview
	Clientset     *kubernetes.Clientset
	DynamicClient *dynamic.DynamicClient
	Error         error
	Config        *rest.Config
	Helm          *Helm
}

type APIResourceMessage struct {
	SessionID  string `json:"sessionId"`
	StatusCode uint   `json:"statusCode"`
	Error      string `json:"error"`
}

var apiRes = APIResource{}

func (ar *APIResource) DataStreamWSHandler(c echo.Context) error {
	if !ar.AuthState && ar.Config == nil {
		return c.JSON(http.StatusUnauthorized, APIResourceMessage{StatusCode: http.StatusUnauthorized})
	}

	sockjs.NewHandler("/srv/data", sockjs.DefaultOptions, handleDataStreamSession).
		ServeHTTP(c.Response(), c.Request())

	return nil
}

func (ar *APIResource) DataStreamWS(c echo.Context) error {
	if !ar.AuthState && ar.Config == nil {
		return c.JSON(http.StatusUnauthorized,
			APIResourceMessage{StatusCode: http.StatusUnauthorized, Error: "Not authenticated"})
	}

	sessionID, err := genSessionId()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResourceMessage{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	dataSessions.Set(sessionID, DataStreamSession{
		id:    sessionID,
		bound: make(chan error),
		close: make(chan struct{}),
	})
	go WaitForDataStream(ar, sessionID)

	return c.JSON(http.StatusOK, APIResourceMessage{SessionID: sessionID, StatusCode: http.StatusOK})
}

func (ar *APIResource) ShellWSHandler(c echo.Context) error {
	if !ar.AuthState && ar.Config == nil {
		return c.JSON(http.StatusUnauthorized, APIResourceMessage{StatusCode: http.StatusUnauthorized})
	}

	sockjs.NewHandler("/srv/shell", sockjs.DefaultOptions, handleTerminalSession).
		ServeHTTP(c.Response(), c.Request())

	return nil
}

func (ar *APIResource) ExecShell(c echo.Context) error {
	if !ar.AuthState && ar.Config == nil {
		return c.JSON(http.StatusUnauthorized,
			APIResourceMessage{StatusCode: http.StatusUnauthorized, Error: "Not authenticated"})
	}

	var ed ExecData
	ed.PodName = c.QueryParam("name")
	ed.PodNamespace = c.QueryParam("namespace")
	ed.ContainerName = c.QueryParam("container")

	shells, err := ed.executeRemoteCommand(ar.Clientset, ar.Config,
		[]string{"grep", "-E", "^/bin", "/etc/shells"})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResourceMessage{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	for _, shell := range strings.Split(shells, "/\n") {
		if strings.Contains(shell, "bash") {
			ed.Command = []string{"bash"}
			ed.Shell = "bash"
			break
		} else {
			ed.Command = []string{"sh"}
			ed.Shell = "sh"
		}
	}

	sessionID, err := genSessionId()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResourceMessage{Error: err.Error()})
	}

	terminalSessions.Set(sessionID, TerminalSession{
		id:       sessionID,
		bound:    make(chan error),
		sizeChan: make(chan remotecommand.TerminalSize),
	})
	go ed.WaitForTerminal(ar.Clientset, ar.Config, sessionID)

	return c.JSON(http.StatusOK, APIResourceMessage{SessionID: sessionID, StatusCode: http.StatusOK})
}

func (ar *APIResource) LogsWSHandler(c echo.Context) error {
	if !ar.AuthState && ar.Config == nil {
		return c.JSON(http.StatusUnauthorized, APIResourceMessage{StatusCode: http.StatusUnauthorized})
	}

	sockjs.NewHandler("/srv/logs", sockjs.DefaultOptions,
		handleLogStreamSession).ServeHTTP(c.Response(), c.Request())

	return nil
}

func (ar *APIResource) StreamLogs(c echo.Context) error {
	if !ar.AuthState && ar.Config == nil {
		return c.JSON(http.StatusUnauthorized, APIResourceMessage{StatusCode: http.StatusUnauthorized})
	}

	var pld PodLogsData
	pld.ResourceType = "pods"
	pld.Name = c.QueryParam("name")
	pld.Namespace = c.QueryParam("namespace")
	pld.ContainerName = c.QueryParam("container")
	pld.Follow = c.QueryParam("follow")
	pld.TailLines = c.QueryParam("tailLines")
	pld.Options = &v1.PodLogOptions{}
	pld.ParameterCodec = scheme.ParameterCodec

	sessionID, err := genSessionId()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, LogStreamMessage{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	logStreamSessions.Set(sessionID, LogStreamSession{
		id:    sessionID,
		bound: make(chan error),
		close: make(chan struct{}),
	})
	go pld.WaitForLogs(ar.Clientset, sessionID)

	return c.JSON(http.StatusOK, LogStreamMessage{SessionID: sessionID, StatusCode: http.StatusOK})
}

type AuthResponse struct {
	Error        string                          `json:"error"`
	KubeHost     string                          `json:"kubeHost"`
	State        bool                            `json:"state"`
	AccessReview *authv1.SelfSubjectAccessReview `json:"accessReview"`
}

type AuthRequest struct {
	Type           string `json:"type"`
	AccessToken    string `json:"token"`
	KubeconfigPath string `json:"kubeconfigPath"`
	KubeconfigRaw  string `json:"kubeconfigRaw"`
	MasterURL      string `json:"masterURL"`
	TLSInsecure    bool   `json:"tlsInsecure"`
}

func (ar *APIResource) Auth(c echo.Context) (err error) {
	if err := c.Bind(ar.AuthRequest); err != nil {
		return c.JSON(http.StatusUnauthorized, AuthResponse{
			Error: fmt.Sprintf("Auth error: could not initialize cluster config: %s", err),
			State: false,
		})
	}

	ssar, errInit := authInit(ar)
	if errInit != nil {
		return c.JSON(http.StatusUnauthorized, AuthResponse{
			Error:        fmt.Sprintf("Auth error: %s", err),
			State:        false,
			AccessReview: ssar,
		})
	}

	ar.AuthState = true
	ar.SSAR = ssar

	return c.JSON(http.StatusOK, AuthResponse{State: ar.AuthState, KubeHost: ar.Config.Host, AccessReview: ar.SSAR})
}

var portFlagValue string
var portFlag = &cli.StringFlag{
	Name:        "port",
	Usage:       "the port where the UI can be accessed at",
	Value:       "3001",
	Destination: &portFlagValue,
}

var kubeconfigFlagValue string
var kubeconfigFlag = &cli.StringFlag{
	Name:        "kubeconfig",
	Usage:       "kubeconfig file path",
	Value:       filepath.Join(homeDir, ".kube", "config"),
	Destination: &kubeconfigFlagValue,
}

var kubeAccessTokenFlagValue string
var kubeAccessTokenFlag = &cli.StringFlag{
	Name:        "access-token",
	Usage:       "kubernetes cluster access token",
	Destination: &kubeAccessTokenFlagValue,
}

var kubeMasterURLFlagValue string
var kubeMasterURLFlag = &cli.StringFlag{
	Name:        "master-url",
	Usage:       "kubernetes master URL",
	Destination: &kubeMasterURLFlagValue,
}

var kubeInClusterConfigFlagValue bool
var kubeInClusterConfigFlag = &cli.BoolFlag{
	Name:        "in-cluster",
	Usage:       "set this flag if the app is inside the kubernetes cluster",
	Destination: &kubeInClusterConfigFlagValue,
}

var AppVersion = "0.0.0"

func main() {
	app := &cli.App{
		HideHelpCommand: true,
		Name:            "lutho",
		Usage:           "Manage Kubernetes based resources in a different way",
		Version:         AppVersion,
		Compiled:        time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "KubeFlame",
				Email: "https://github.com/kubeflame",
			},
		},
		Commands: []*cli.Command{
			// {
			// 	Name:        "config",
			// 	Usage:       "manage the configuration",
			// 	Description: "",
			// 	Subcommands: []*cli.Command{
			// 		{
			// 			Name: "get",
			// 			Action: func(ctx *cli.Context) error {
			// 				return nil
			// 			},
			// 		},
			// 		{
			// 			Name: "set",
			// 			Action: func(ctx *cli.Context) error {
			// 				return nil
			// 			},
			// 		},
			// 		{
			// 			Name: "list",
			// 			Action: func(ctx *cli.Context) error {
			// 				return nil
			// 			},
			// 		},
			// 	},
			// },
			{
				Name:        "start",
				Usage:       "Starts the application",
				Description: "This command will start the application",
				Flags: []cli.Flag{
					portFlag,
					kubeconfigFlag,
					kubeAccessTokenFlag,
					kubeMasterURLFlag,
					kubeInClusterConfigFlag,
				},
				Action: func(ctx *cli.Context) error {
					apiRes.CliContext = ctx
					apiRes.AuthRequest = &AuthRequest{}
					apiRes.SSAR = &authv1.SelfSubjectAccessReview{}

					switch {
					case apiRes.CliContext.IsSet(kubeconfigFlag.Name):
						apiRes.AuthRequest.Type = KubernetesConfigType.kubeconfigPath
						apiRes.AuthRequest.KubeconfigPath = kubeconfigFlagValue
						ssar, errInit := authInit(&apiRes)
						if errInit != nil {
							return errInit
						}
						apiRes.AuthState = true
						apiRes.SSAR = ssar
					case apiRes.CliContext.IsSet(kubeAccessTokenFlag.Name):
						apiRes.AuthRequest.Type = KubernetesConfigType.accessToken
						apiRes.AuthRequest.MasterURL = kubeMasterURLFlagValue
						apiRes.AuthRequest.AccessToken = kubeAccessTokenFlagValue
						apiRes.AuthRequest.TLSInsecure = true
						ssar, errInit := authInit(&apiRes)
						if errInit != nil {
							return errInit
						}
						apiRes.AuthState = true
						apiRes.SSAR = ssar
					case apiRes.CliContext.IsSet(kubeInClusterConfigFlag.Name):
						apiRes.AuthRequest.Type = KubernetesConfigType.inClusterConfig
						ssar, errInit := authInit(&apiRes)
						if errInit != nil {
							return errInit
						}
						apiRes.AuthState = true
						apiRes.SSAR = ssar
					}

					e := echo.New()
					e.HideBanner = true

					config := middleware.RateLimiterConfig{
						Skipper: middleware.DefaultSkipper,
						Store: middleware.NewRateLimiterMemoryStoreWithConfig(
							middleware.RateLimiterMemoryStoreConfig{
								Rate: 60, Burst: 30, ExpiresIn: time.Minute,
							},
						),
						IdentifierExtractor: func(c echo.Context) (string, error) {
							id := c.RealIP()
							return id, nil
						},
						ErrorHandler: func(c echo.Context, err error) error {
							return c.JSON(http.StatusForbidden,
								APIResourceMessage{StatusCode: http.StatusForbidden})
						},
						DenyHandler: func(c echo.Context, identifier string, err error) error {
							return c.JSON(http.StatusTooManyRequests,
								APIResourceMessage{StatusCode: http.StatusTooManyRequests})
						},
					}

					e.Use(middleware.RateLimiterWithConfig(config))
					e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
						Level: 9,
					}))

					fs := echo.MustSubFS(files, "frontend/dist")
					e.StaticFS("/", fs)

					e.POST("/srv/auth", apiRes.Auth)
					e.GET("/srv/auth/state", func(c echo.Context) error {
						var kubeHost string
						if apiRes.Config != nil {
							kubeHost = apiRes.Config.Host
						}
						return c.JSON(http.StatusOK, AuthResponse{
							State: apiRes.AuthState, KubeHost: kubeHost, AccessReview: apiRes.SSAR,
						},
						)
					})

					e.GET("/srv/data*", apiRes.DataStreamWSHandler)
					e.GET("/srv/data/ws", apiRes.DataStreamWS)

					e.GET("/srv/shell*", apiRes.ShellWSHandler)
					e.GET("/srv/shell/exec", apiRes.ExecShell)

					e.GET("/srv/logs*", apiRes.LogsWSHandler)
					e.GET("/srv/logs/stream", apiRes.StreamLogs)

					e.Logger.Fatal(e.Start(":" + portFlagValue))

					return nil
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	if err := app.Run(os.Args); err != nil {
		log.Fatal("Application could not start: ", err)
	}

}
