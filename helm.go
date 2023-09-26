package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/gommon/log"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"sigs.k8s.io/yaml"
)

const (
	defaultHelmRegistry   = "/.lutho/helm/registry"
	defaultHelmRepository = "/.lutho/helm/repository"
)

func initHelm(ar *APIResource) (*action.Configuration, *cli.EnvSettings, error) {
	if err := os.MkdirAll(homeDir+defaultHelmRegistry, os.ModePerm); err != nil {
		return nil, nil, err
	}
	if err := os.MkdirAll(homeDir+defaultHelmRepository, os.ModePerm); err != nil {
		return nil, nil, err
	}

	actionConfig := new(action.Configuration)

	settings := &cli.EnvSettings{}
	settings.KubeAPIServer = ar.Config.Host
	settings.KubeToken = ar.Config.BearerToken
	settings.KubeCaFile = ar.Config.CAFile
	settings.KubeInsecureSkipTLSVerify = ar.Config.Insecure
	settings.RegistryConfig = filepath.Join(homeDir, defaultHelmRegistry, "config.json")
	settings.RepositoryCache = homeDir + defaultHelmRepository

	return actionConfig, settings, nil
}

func GetActionConfig(ctx context.Context, namespace string, config *rest.Config) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)
	cliConfig := genericclioptions.NewConfigFlags(false)
	cliConfig.APIServer = &config.Host
	cliConfig.BearerToken = &config.BearerToken
	cliConfig.Namespace = &namespace

	// Drop their rest.Config and just return inject own
	wrapper := func(*rest.Config) *rest.Config {
		return config
	}
	cliConfig.WithWrapConfigFn(wrapper)
	if err := actionConfig.Init(cliConfig, namespace, "secret", log.Debugf); err != nil {
		return nil, err
	}
	return actionConfig, nil
}

type Helm struct {
	ActionConfig     *action.Configuration
	EnvSettings      *cli.EnvSettings
	ReleaseName      string
	ReleaseNamespace string
	AllValues        bool
	AllNamespaces    bool
	ByDate           bool
	Limit            int
	RepoURL          string
	ChartName        string
	ChartVersion     string
	ReleaseVersion   int
	DryRun           bool
}

func (h *Helm) ShowChartValues(config *rest.Config, opts HelmOptions) (vals string, err error) {
	h.ActionConfig, err = GetActionConfig(context.Background(), h.ReleaseNamespace, config)
	if err != nil {
		return "", err
	}

	sc := action.NewShow(action.ShowOutputFormat("values"))
	sc.Version = opts.ChartVersion
	sc.ChartPathOptions.RepoURL = opts.RepoURL

	chartRef := opts.ChartName
	if opts.IsOCI {
		rc, errReg := newRegistryClient(h.EnvSettings, true)
		if errReg != nil {
			return "", fmt.Errorf("failed to create registry: %w", errReg)
		}
		sc.SetRegistryClient(rc)
		chartRef = opts.RepoURL
		sc.ChartPathOptions.RepoURL = ""
	}

	cp, err := sc.ChartPathOptions.LocateChart(chartRef, h.EnvSettings)
	if err != nil {
		return "", fmt.Errorf("failed to locate chart: %w", err)
	}

	return sc.Run(cp)
}

func (h *Helm) ListReleases(config *rest.Config) (rel []*release.Release, err error) {
	h.ActionConfig, err = GetActionConfig(context.Background(), h.ReleaseNamespace, config)
	if err != nil {
		return nil, err
	}

	nl := action.NewList(h.ActionConfig)
	nl.AllNamespaces = h.AllNamespaces
	nl.ByDate = h.ByDate
	nl.Limit = h.Limit

	return nl.Run()
}

func (h *Helm) GetRelease(config *rest.Config) (rel *release.Release, err error) {
	h.ActionConfig, err = GetActionConfig(context.Background(), h.ReleaseNamespace, config)
	if err != nil {
		return nil, err
	}

	gr := action.NewGet(h.ActionConfig)
	gr.Version = h.ReleaseVersion

	return gr.Run(h.ReleaseName)
}

func newRegistryClient(settings *cli.EnvSettings, plainHTTP bool) (*registry.Client, error) {
	opts := []registry.ClientOption{
		registry.ClientOptDebug(settings.Debug),
		registry.ClientOptEnableCache(true),
		registry.ClientOptWriter(os.Stderr),
		registry.ClientOptCredentialsFile(settings.RegistryConfig),
	}
	if plainHTTP {
		opts = append(opts, registry.ClientOptPlainHTTP())
	}

	// Create a new registry client
	registryClient, err := registry.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return registryClient, nil
}

func (h *Helm) InstallRelease(config *rest.Config, opts HelmOptions, values map[string]interface{}) (rel *release.Release, err error) {
	h.ActionConfig, err = GetActionConfig(context.Background(), h.ReleaseNamespace, config)
	if err != nil {
		return nil, err
	}

	i := action.NewInstall(h.ActionConfig)
	i.Version = opts.ChartVersion
	i.ReleaseName = h.ReleaseName
	i.Namespace = h.ReleaseNamespace
	i.DryRun = h.DryRun
	i.ChartPathOptions.RepoURL = opts.RepoURL

	chartRef := opts.ChartName
	if opts.IsOCI {
		rc, errReg := newRegistryClient(h.EnvSettings, true)
		if errReg != nil {
			return nil, fmt.Errorf("failed to create registry: %w", errReg)
		}
		i.SetRegistryClient(rc)
		chartRef = opts.RepoURL
		i.ChartPathOptions.RepoURL = ""
	}

	cp, err := i.ChartPathOptions.LocateChart(chartRef, h.EnvSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to locate chart: %w", err)
	}

	chart, err := loader.Load(cp)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart: %w", err)
	}

	providers := getter.All(h.EnvSettings)

	if dep := chart.Metadata.Dependencies; dep != nil {
		if err := action.CheckDependencies(chart, dep); err != nil {
			err = fmt.Errorf("failed to check chart dependencies: %w", err)
			if !i.DependencyUpdate {
				return nil, err
			}

			man := &downloader.Manager{
				Out:              os.Stdout,
				ChartPath:        cp,
				Keyring:          i.ChartPathOptions.Keyring,
				SkipUpdate:       false,
				Getters:          providers,
				RepositoryConfig: h.EnvSettings.RepositoryConfig,
				RepositoryCache:  h.EnvSettings.RepositoryCache,
				Debug:            h.EnvSettings.Debug,
			}
			if err := man.Update(); err != nil {
				return nil, err
			}
			// Reload the chart with the updated Chart.lock file.
			if chart, err = loader.Load(cp); err != nil {
				return nil, fmt.Errorf("failed to reload chart after repo update: %w", err)
			}
		}
	}

	return i.RunWithContext(context.TODO(), chart, values)
}

func (h *Helm) UpgradeRelease(config *rest.Config, opts HelmOptions, values map[string]interface{}) (rel *release.Release, err error) {
	h.ActionConfig, err = GetActionConfig(context.Background(), h.ReleaseNamespace, config)
	if err != nil {
		return nil, err
	}

	u := action.NewUpgrade(h.ActionConfig)
	u.ChartPathOptions.Version = opts.ChartVersion
	u.Namespace = h.ReleaseNamespace
	u.DryRun = h.DryRun
	u.ReuseValues = opts.ReuseValues
	u.ChartPathOptions.RepoURL = opts.RepoURL

	chartRef := opts.ChartName
	if opts.IsOCI {
		rc, errReg := newRegistryClient(h.EnvSettings, true)
		if errReg != nil {
			return nil, fmt.Errorf("failed to create registry: %w", errReg)
		}
		u.SetRegistryClient(rc)
		chartRef = opts.RepoURL
		u.ChartPathOptions.RepoURL = ""
	}

	cp, err := u.ChartPathOptions.LocateChart(chartRef, h.EnvSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to locate chart: %w", err)
	}

	chart, err := loader.Load(cp)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart: %w", err)
	}

	providers := getter.All(h.EnvSettings)

	if dep := chart.Metadata.Dependencies; dep != nil {
		if err := action.CheckDependencies(chart, dep); err != nil {
			err = fmt.Errorf("failed to check chart dependencies: %w", err)
			if !u.DependencyUpdate {
				return nil, err
			}

			man := &downloader.Manager{
				Out:              os.Stdout,
				ChartPath:        cp,
				Keyring:          u.ChartPathOptions.Keyring,
				SkipUpdate:       false,
				Getters:          providers,
				RepositoryConfig: h.EnvSettings.RepositoryConfig,
				RepositoryCache:  h.EnvSettings.RepositoryCache,
				Debug:            h.EnvSettings.Debug,
			}
			if err := man.Update(); err != nil {
				return nil, err
			}
			// Reload the chart with the updated Chart.lock file.
			if chart, err = loader.Load(cp); err != nil {
				return nil, fmt.Errorf("failed to reload chart after repo update: %w", err)
			}
		}
	}

	return u.RunWithContext(context.TODO(), h.ReleaseName, chart, values)
}

func (h *Helm) UninstallRelease(config *rest.Config) (urr *release.UninstallReleaseResponse, err error) {
	h.ActionConfig, err = GetActionConfig(context.Background(), h.ReleaseNamespace, config)
	if err != nil {
		return nil, err
	}

	u := action.NewUninstall(h.ActionConfig)
	u.DryRun = h.DryRun

	return u.Run(h.ReleaseName)
}

func runPull(h *Helm, config *rest.Config) (result string, err error) {
	h.ActionConfig, err = GetActionConfig(context.Background(), h.ReleaseNamespace, config)
	if err != nil {
		return "", err
	}

	registryClient, err := newRegistryClient(h.EnvSettings, false)
	if err != nil {
		return "", fmt.Errorf("failed to created registry client: %w", err)
	}
	h.ActionConfig.RegistryClient = registryClient

	pullClient := action.NewPullWithOpts(action.WithConfig(h.ActionConfig))
	pullClient.DestDir = h.EnvSettings.RepositoryCache
	pullClient.Settings = h.EnvSettings
	pullClient.Version = h.ChartVersion

	return pullClient.Run(h.RepoURL)
}

func (h *Helm) PullChart(opts HelmOptions, config *rest.Config) (repo string, err error) {
	if opts.EnvPath == "" {
		opts.EnvPath = "/.lutho/helm"
	}

	h.EnvSettings.RegistryConfig = homeDir + opts.EnvPath + "/registry/config.json"
	h.EnvSettings.RepositoryConfig = homeDir + opts.EnvPath + "/repositories.yaml"
	h.EnvSettings.RepositoryCache = homeDir + opts.EnvPath + "/repository"

	result, err := runPull(h, config)
	if err != nil {
		return "", err
	}
	fmt.Println("result helmPull ->", result)

	return "", nil
}

type HelmChartTags struct {
	ChartName string   `json:"chartName"`
	Tags      []string `json:"chartTags"`
}

func (h *Helm) GetChartTags(opts HelmOptions) (ct HelmChartTags, err error) {
	if opts.IsOCI {
		tags, err := h.GetTags()
		if err != nil {
			return ct, err
		}

		ct.ChartName = opts.ChartName
		ct.Tags = tags
	} else {
		f, err := h.GetIndexFile()
		if err != nil {
			return ct, err
		}

		var vers []string
		for _, v := range f.Entries[h.ChartName] {
			vers = append(vers, v.Version)
		}

		ct.ChartName = opts.ChartName
		ct.Tags = vers
	}

	return ct, nil
}

type Entry struct {
	Version string `yaml:"version"`
	// Created time.Time
}

type Entries []Entry

type Index struct {
	Entries map[string]Entries
}

func (h *Helm) GetIndexFile() (*Index, error) {
	cr, crErr := repo.NewChartRepository(&repo.Entry{
		InsecureSkipTLSverify: false,
		PassCredentialsAll:    false,
		Name:                  h.ChartName,
		URL:                   h.RepoURL,
	}, getter.All(h.EnvSettings))
	if crErr != nil {
		return nil, fmt.Errorf("failed to initialize index file repository: %w", crErr)
	}

	cr.CachePath = h.EnvSettings.RepositoryCache

	indexURL, err := repo.ResolveReferenceURL(h.RepoURL, "index.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve index file URL: %w", err)
	}

	resp, err := cr.Client.Get(indexURL,
		getter.WithURL(cr.Config.URL),
		getter.WithInsecureSkipVerifyTLS(cr.Config.InsecureSkipTLSverify),
		getter.WithTLSClientConfig(cr.Config.CertFile, cr.Config.KeyFile, cr.Config.CAFile),
		getter.WithBasicAuth(cr.Config.Username, cr.Config.Password),
		getter.WithPassCredentialsAll(cr.Config.PassCredentialsAll),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to download index file: %w", err)
	}

	idxF, err := io.ReadAll(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read index file: %w", err)
	}

	index := &Index{}
	errUnmarshal := yaml.Unmarshal(idxF, index)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}

	return index, nil
}

func (h *Helm) GetTags() (tags []string, err error) {
	tagsURL, err := url.Parse(h.RepoURL)
	if err != nil {
		return nil, err
	}

	remoteRepo, err := remote.NewRepository(tagsURL.Host + tagsURL.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	client := &http.Client{Transport: &http.Transport{
		Proxy:             http.ProxyFromEnvironment,
		TLSClientConfig:   &tls.Config{},
		DisableKeepAlives: true,
	}}

	remoteRepo.Client = &auth.Client{
		Client: client,
		Cache:  nil,
		Credential: auth.StaticCredential(tagsURL.Host, auth.Credential{
			Username: "",
			Password: "",
		}),
	}

	remoteRepo.PlainHTTP = tagsURL.Scheme == "http"

	errTags := remoteRepo.Tags(context.Background(), "", func(tagsResult []string) error {
		for _, tag := range tagsResult {
			// By convention: Change underscore (_) back to plus (+) to get valid SemVer
			convertedTag := strings.ReplaceAll(tag, "_", "+")
			if !strings.Contains(convertedTag, "sha256") {
				tags = append(tags, convertedTag)
			}

		}

		return nil
	})
	if errTags != nil {
		return nil, fmt.Errorf("failed to get tags: %w", errTags)
	}

	reverseStringSlice(tags)

	return tags, nil
}
