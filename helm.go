package main

import (
	"os"

	"github.com/labstack/gommon/log"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

func initHelm() (*action.Configuration, *cli.EnvSettings, error) {
	actionConfig := new(action.Configuration)
	settings := cli.New()

	return actionConfig, settings, nil
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
	Version          int
}

func (h *Helm) ListReleases() ([]*release.Release, error) {
	err := h.ActionConfig.Init(h.EnvSettings.RESTClientGetter(), h.ReleaseNamespace, os.Getenv("HELM_DRIVER"), log.Printf)
	if err != nil {
		return nil, err
	}

	// h.EnvSettings.SetNamespace(h.ReleaseNamespace)

	nl := action.NewList(h.ActionConfig)
	nl.AllNamespaces = h.AllNamespaces
	nl.ByDate = h.ByDate
	nl.Limit = h.Limit

	return nl.Run()
}

func (h *Helm) GetValues() (map[string]interface{}, error) {
	// h.EnvSettings.SetNamespace(h.ReleaseNamespace)
	err := h.ActionConfig.Init(h.EnvSettings.RESTClientGetter(), h.ReleaseNamespace, os.Getenv("HELM_DRIVER"), log.Printf)
	if err != nil {
		return nil, err
	}

	gv := action.NewGetValues(h.ActionConfig)
	gv.AllValues = h.AllValues
	gv.Version = h.Version

	return gv.Run(h.ReleaseName)
}

func (h *Helm) UpgradeRelease() (*release.Release, error) {
	// h.EnvSettings.SetNamespace(h.ReleaseNamespace)
	err := h.ActionConfig.Init(h.EnvSettings.RESTClientGetter(), h.ReleaseNamespace, os.Getenv("HELM_DRIVER"), log.Printf)
	if err != nil {
		return nil, err
	}

	// gv := action.NewGetValues(h.ActionConfig)
	u := action.NewUpgrade(h.ActionConfig)

	u.ChartPathOptions = action.ChartPathOptions{RepoURL: ""}
	// gv.AllValues = h.AllValues
	// gv.Version = h.Version

	return u.Run(h.ReleaseName, &chart.Chart{}, make(map[string]interface{}))
}

func (h *Helm) DeleteRelease() (map[string]interface{}, error) {
	// h.EnvSettings.SetNamespace(h.ReleaseNamespace)
	err := h.ActionConfig.Init(h.EnvSettings.RESTClientGetter(), h.ReleaseNamespace, os.Getenv("HELM_DRIVER"), log.Printf)
	if err != nil {
		return nil, err
	}

	gv := action.NewGetValues(h.ActionConfig)
	gv.AllValues = h.AllValues
	gv.Version = h.Version

	return gv.Run(h.ReleaseName)
}
