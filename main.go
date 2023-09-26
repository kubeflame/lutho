package main

import (
	"embed"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/igm/sockjs-go/v3/sockjs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
)

//go:generate yarn --cwd ./frontend install
//go:generate yarn --cwd ./frontend run build
//go:embed all:frontend/dist
var files embed.FS

type FrontendRequest struct {
	Namespace       string          `json:"namespace"`
	Name            string          `json:"name"`
	Group           string          `json:"group"`
	Version         string          `json:"version"`
	Resource        string          `json:"resource"`
	Kind            string          `json:"kind"`
	Data            string          `json:"data"`
	ResourceOptions ResourceOptions `json:"options"`
}

type ResourceOptions struct {
	FieldSelector  string `json:"fieldSelector"`
	LabelSelector  string `json:"labelSelector"`
	TimeoutSeconds int64  `json:"timeoutSeconds"`
	Limit          int64  `json:"limit"`
	Continue       string `json:"_continue"`
}

type APIResource struct {
	AuthState     bool
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
	if !ar.AuthState {
		return c.JSON(http.StatusUnauthorized, APIResourceMessage{StatusCode: http.StatusUnauthorized})
	}

	sockjs.NewHandler("/srv/data", sockjs.DefaultOptions, handleDataStreamSession).
		ServeHTTP(c.Response(), c.Request())

	return nil
}

func (ar *APIResource) DataStreamWS(c echo.Context) error {
	if !ar.AuthState {
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
	if !ar.AuthState {
		return c.JSON(http.StatusUnauthorized, APIResourceMessage{StatusCode: http.StatusUnauthorized})
	}

	sockjs.NewHandler("/srv/shell", sockjs.DefaultOptions, handleTerminalSession).
		ServeHTTP(c.Response(), c.Request())

	return nil
}

func (ar *APIResource) ExecShell(c echo.Context) error {
	if !ar.AuthState {
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
	if !ar.AuthState {
		return c.JSON(http.StatusUnauthorized, APIResourceMessage{StatusCode: http.StatusUnauthorized})
	}

	sockjs.NewHandler("/srv/logs", sockjs.DefaultOptions,
		handleLogStreamSession).ServeHTTP(c.Response(), c.Request())

	return nil
}

func (ar *APIResource) StreamLogs(c echo.Context) error {
	if !ar.AuthState {
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
	Error error `json:"error"`
	State bool  `json:"state"`
}

func (ar *APIResource) Auth(c echo.Context) (err error) {
	cfg, err := initKubeconfig(&KubeconfigInit{
		TLSInsecure: true,
		Type:        Config,
	})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, AuthResponse{
			Error: fmt.Errorf("could not initialize kubeconfig: %s", err),
		})
	}
	ar.Config = cfg

	cs, err := initClientset(ar.Config)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, AuthResponse{
			Error: fmt.Errorf("could not initialize kubernetes clientset: %s", err),
		})
	}
	ar.Clientset = cs

	dc, err := initDynamicClient(ar.Config)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, AuthResponse{
			Error: fmt.Errorf("could not initialize kubernetes dynamic client: %s", err),
		})
	}
	ar.DynamicClient = dc

	hac, hes, err := initHelm()
	if err != nil {
		return c.JSON(http.StatusUnauthorized, AuthResponse{
			Error: fmt.Errorf("could not initialize helm client: %s", err),
		})
	}
	ar.Helm = &Helm{ActionConfig: hac, EnvSettings: hes}

	ar.AuthState = true

	return c.JSON(http.StatusOK, AuthResponse{State: true})
}

func main() {
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

	fs := echo.MustSubFS(files, "frontend/dist")
	e.StaticFS("/", fs)

	e.GET("/auth", apiRes.Auth)

	e.GET("/srv/data*", apiRes.DataStreamWSHandler)
	e.GET("/srv/data/ws", apiRes.DataStreamWS)

	e.GET("/srv/shell*", apiRes.ShellWSHandler)
	e.GET("/srv/shell/exec", apiRes.ExecShell)

	e.GET("/srv/logs*", apiRes.LogsWSHandler)
	e.GET("/srv/logs/stream", apiRes.StreamLogs)

	e.Logger.Fatal(e.Start(":3001"))
}
