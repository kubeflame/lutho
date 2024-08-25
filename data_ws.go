package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/yaml"
)

type DataStreamSession struct {
	id    string
	bound chan error
	ws    *websocket.Conn
	close chan struct{}
}

type DataStreamMessage struct {
	Op         DataStreamOp `json:"op"`
	Data       string       `json:"data"`
	SessionID  string       `json:"sessionId"`
	Error      string       `json:"error"`
	StatusCode uint         `json:"statusCode"`
}

type DataStreamOp struct {
	Type    string          `json:"type"`
	Request FrontendRequest `json:"request"`
	OpID    string          `json:"opID"`
}

type DataSessionMap struct {
	Sessions map[string]DataStreamSession
	Lock     sync.RWMutex
}

func (dsm *DataSessionMap) Get(sessionId string) DataStreamSession {
	dsm.Lock.RLock()
	defer dsm.Lock.RUnlock()
	return dsm.Sessions[sessionId]
}

func (dsm *DataSessionMap) Set(sessionId string, session DataStreamSession) {
	dsm.Lock.Lock()
	defer dsm.Lock.Unlock()
	dsm.Sessions[sessionId] = session
}

func (dsm *DataSessionMap) Close(sessionId string, status uint, reason string) {
	dsm.Lock.Lock()
	defer dsm.Lock.Unlock()
	ses := dsm.Sessions[sessionId]
	ses.ws.WriteJSON(DataStreamMessage{Op: DataStreamOp{Type: "close"}, Data: reason, StatusCode: status})
	ses.ws.Close()
	delete(dsm.Sessions, sessionId)
}

type DataStream struct {
	client  *dynamic.DynamicClient
	recvMsg DataStreamMessage
	helm    *Helm
}

func (ds DataStream) kubeList() ([]byte, unstructured.UnstructuredList, error) {
	var lr ListResource
	lr.GVR = schema.GroupVersionResource{
		Group:    ds.recvMsg.Op.Request.KubeGVRK.Group,
		Version:  ds.recvMsg.Op.Request.KubeGVRK.Version,
		Resource: ds.recvMsg.Op.Request.KubeGVRK.Resource,
	}
	lr.Namespace = ds.recvMsg.Op.Request.Namespace
	if !ds.recvMsg.Op.Request.KubeGVRK.IsNamespaced {
		lr.Namespace = ""
	}

	sel, err := fields.ParseSelector(ds.recvMsg.Op.Request.ResourceOptions.FieldSelector)
	if err != nil {
		return nil, unstructured.UnstructuredList{}, err
	}

	lr.Options = metav1.ListOptions{
		FieldSelector: sel.String(),
		LabelSelector: ds.recvMsg.Op.Request.ResourceOptions.LabelSelector,
		Limit:         ds.recvMsg.Op.Request.ResourceOptions.Limit,
		Continue:      ds.recvMsg.Op.Request.ResourceOptions.Continue,
	}

	list, err := lr.List(ds.client)
	if err != nil {
		return nil, unstructured.UnstructuredList{}, err
	}

	if ds.recvMsg.Op.Type == WSOpType.listAll {
		return nil, *list, nil
	}

	dataBytes, err := list.MarshalJSON()
	if err != nil {
		return nil, unstructured.UnstructuredList{}, err
	}

	return dataBytes, unstructured.UnstructuredList{}, nil
}

func (ds DataStream) kubeGet() ([]byte, error) {
	var gr GetResource
	gr.Namespace = ds.recvMsg.Op.Request.Namespace
	gr.Name = ds.recvMsg.Op.Request.Name
	gr.Options = metav1.GetOptions{}
	gr.GVR = schema.GroupVersionResource{
		Group:    ds.recvMsg.Op.Request.KubeGVRK.Group,
		Version:  ds.recvMsg.Op.Request.KubeGVRK.Version,
		Resource: ds.recvMsg.Op.Request.KubeGVRK.Resource,
	}

	get, err := gr.Get(ds.client)
	if err != nil {
		return nil, err
	}

	dataBytes, err := get.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return dataBytes, nil
}

func (ds DataStream) kubeUpdate() ([]byte, error) {
	var ur UpdateResource
	ur.GetOptions = metav1.GetOptions{}
	ur.UpdateOptions = metav1.UpdateOptions{FieldValidation: "Strict"}
	ur.Namespace = ds.recvMsg.Op.Request.Namespace
	ur.Name = ds.recvMsg.Op.Request.Name
	ur.Data = []byte(ds.recvMsg.Op.Request.Data)
	ur.GVR = schema.GroupVersionResource{
		Group:    ds.recvMsg.Op.Request.KubeGVRK.Group,
		Version:  ds.recvMsg.Op.Request.KubeGVRK.Version,
		Resource: ds.recvMsg.Op.Request.KubeGVRK.Resource,
	}
	ur.GVK = &schema.GroupVersionKind{
		Group:   ds.recvMsg.Op.Request.KubeGVRK.Group,
		Version: ds.recvMsg.Op.Request.KubeGVRK.Version,
		Kind:    ds.recvMsg.Op.Request.KubeGVRK.Kind,
	}

	update, err := ur.Update(ds.client)
	if err != nil {
		return nil, err
	}

	dataBytes, err := update.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return dataBytes, nil
}

func (ds DataStream) kubeDelete() ([]byte, error) {
	var dr DeleteResource
	dr.Options = metav1.DeleteOptions{}
	dr.Namespace = ds.recvMsg.Op.Request.Namespace
	dr.Name = ds.recvMsg.Op.Request.Name
	dr.GVR = schema.GroupVersionResource{
		Group:    ds.recvMsg.Op.Request.KubeGVRK.Group,
		Version:  ds.recvMsg.Op.Request.KubeGVRK.Version,
		Resource: ds.recvMsg.Op.Request.KubeGVRK.Resource,
	}

	err := dr.Delete(ds.client)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds DataStream) helmShowValues(config *rest.Config) (string, error) {
	var h Helm
	h.ActionConfig = ds.helm.ActionConfig
	h.EnvSettings = ds.helm.EnvSettings

	return h.ShowChartValues(config, ds.recvMsg.Op.Request.HelmOptions)
}

func (ds DataStream) helmList(config *rest.Config) ([]byte, error) {
	var h Helm
	h.ActionConfig = ds.helm.ActionConfig
	h.EnvSettings = ds.helm.EnvSettings
	h.ReleaseNamespace = ds.recvMsg.Op.Request.Namespace

	lr, err := h.ListReleases(config)
	if err != nil {
		return nil, err
	}

	var relList []release.Release
	for _, v := range lr {
		relList = append(relList, release.Release{
			Name:      v.Name,
			Namespace: v.Namespace,
			Info:      v.Info,
			Chart: &chart.Chart{
				Metadata: v.Chart.Metadata,
			},
		})
	}

	dataBytes, err := json.Marshal(relList)
	if err != nil {
		return nil, err
	}

	return dataBytes, nil
}

func (ds DataStream) helmGet(config *rest.Config) ([]byte, error) {
	var h Helm
	h.ActionConfig = ds.helm.ActionConfig
	h.EnvSettings = ds.helm.EnvSettings
	h.AllNamespaces = false
	h.ReleaseName = ds.recvMsg.Op.Request.Name
	h.ReleaseNamespace = ds.recvMsg.Op.Request.Namespace
	h.RepoURL = ds.recvMsg.Op.Request.HelmOptions.RepoURL
	h.ChartVersion = ds.recvMsg.Op.Request.HelmOptions.ChartVersion
	h.AllValues = true

	rl, err := h.GetRelease(config)
	if err != nil {
		return nil, err
	}

	var relValues chartutil.Values
	if h.AllValues {
		relValues, err = chartutil.CoalesceValues(rl.Chart, rl.Config)
		if err != nil {
			return nil, err
		}
	} else {
		relValues = rl.Config
	}

	dataBytes, err := json.Marshal(release.Release{
		Name:      rl.Name,
		Namespace: rl.Namespace,
		Chart: &chart.Chart{
			Metadata: rl.Chart.Metadata,
		},
		Config: relValues,
		Info:   rl.Info,
	})
	if err != nil {
		return nil, err
	}

	return dataBytes, nil
}

func (ds DataStream) helmPull(config *rest.Config) (string, error) {
	var h Helm
	h.ActionConfig = ds.helm.ActionConfig
	h.EnvSettings = ds.helm.EnvSettings
	h.ChartName = ds.recvMsg.Op.Request.HelmOptions.ChartName
	h.RepoURL = ds.recvMsg.Op.Request.HelmOptions.RepoURL
	h.ChartVersion = ds.recvMsg.Op.Request.HelmOptions.ChartVersion

	rl, err := h.PullChart(ds.recvMsg.Op.Request.HelmOptions, config)
	if err != nil {
		fmt.Println("err PullChart:", err)
		return "", err
	}

	fmt.Println("chart string:", rl)

	// dataBytes, err := json.Marshal(rl)
	// if err != nil {
	// 	return "", err
	// }

	return rl, nil
}

func (ds DataStream) helmChartTags() (HelmChartTags, error) {
	var h Helm
	h.ActionConfig = ds.helm.ActionConfig
	h.EnvSettings = ds.helm.EnvSettings
	h.DryRun = ds.recvMsg.Op.Request.HelmOptions.DryRun
	h.RepoURL = ds.recvMsg.Op.Request.HelmOptions.RepoURL
	h.ChartName = ds.recvMsg.Op.Request.HelmOptions.ChartName
	h.ChartVersion = ds.recvMsg.Op.Request.HelmOptions.ChartVersion

	return h.GetChartTags(ds.recvMsg.Op.Request.HelmOptions)
}

func (ds DataStream) helmInstall(config *rest.Config, opts HelmOptions) ([]byte, error) {
	var h Helm
	h.ActionConfig = ds.helm.ActionConfig
	h.EnvSettings = ds.helm.EnvSettings
	h.RepoURL = ds.recvMsg.Op.Request.HelmOptions.RepoURL
	h.ChartName = ds.recvMsg.Op.Request.HelmOptions.ChartName
	h.ChartVersion = ds.recvMsg.Op.Request.HelmOptions.ChartVersion
	h.ReleaseName = ds.recvMsg.Op.Request.Name
	h.ReleaseNamespace = ds.recvMsg.Op.Request.Namespace
	h.DryRun = ds.recvMsg.Op.Request.HelmOptions.DryRun
	h.AllNamespaces = false
	h.AllValues = true

	var vals chartutil.Values
	if err := yaml.Unmarshal([]byte(ds.recvMsg.Op.Request.Data), &vals); err != nil {
		return nil, err
	}

	ir, err := h.InstallRelease(config, opts, vals)
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(release.Release{
		Name:      ir.Name,
		Namespace: ir.Namespace,
		Chart: &chart.Chart{
			Metadata: ir.Chart.Metadata,
		},
		Config: ir.Config,
		Info:   ir.Info,
	})
	if err != nil {
		return nil, err
	}

	return dataBytes, nil
}

func (ds DataStream) helmUpgrade(config *rest.Config, opts HelmOptions) ([]byte, error) {
	var h Helm
	h.ActionConfig = ds.helm.ActionConfig
	h.EnvSettings = ds.helm.EnvSettings
	h.RepoURL = ds.recvMsg.Op.Request.HelmOptions.RepoURL
	h.ChartName = ds.recvMsg.Op.Request.HelmOptions.ChartName
	h.ChartVersion = ds.recvMsg.Op.Request.HelmOptions.ChartVersion
	h.ReleaseName = ds.recvMsg.Op.Request.Name
	h.ReleaseNamespace = ds.recvMsg.Op.Request.Namespace
	h.DryRun = ds.recvMsg.Op.Request.HelmOptions.DryRun
	h.AllNamespaces = false
	h.AllValues = true

	var vals chartutil.Values
	if err := yaml.Unmarshal([]byte(ds.recvMsg.Op.Request.Data), &vals); err != nil {
		return nil, err
	}

	ur, err := h.UpgradeRelease(config, opts, vals)
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(release.Release{
		Name:      ur.Name,
		Namespace: ur.Namespace,
		Chart: &chart.Chart{
			Metadata: ur.Chart.Metadata,
		},
		Config: ur.Config,
		Info:   ur.Info,
	})
	if err != nil {
		return nil, err
	}

	return dataBytes, nil
}

func (ds DataStream) helmUninstall(config *rest.Config) (*release.UninstallReleaseResponse, error) {
	var h Helm
	h.ActionConfig = ds.helm.ActionConfig
	h.EnvSettings = ds.helm.EnvSettings
	h.ReleaseName = ds.recvMsg.Op.Request.Name
	h.ReleaseNamespace = ds.recvMsg.Op.Request.Namespace
	h.DryRun = ds.recvMsg.Op.Request.HelmOptions.DryRun

	urr, err := h.UninstallRelease(config)
	if err != nil {
		return nil, err
	}

	return urr, nil
}

func (dss DataStreamSession) List(recvMsg DataStreamMessage, client *dynamic.DynamicClient) error {
	var ds DataStream
	ds.client = client
	ds.recvMsg = recvMsg

	dsm := DataStreamMessage{
		Op: DataStreamOp{
			OpID: recvMsg.Op.OpID,
			Type: WSOpType.list,
		},
	}

	data, _, kubeErr := ds.kubeList()
	if kubeErr != nil {
		dsm.Error = kubeErr.Error()
	}

	dsm.Data = string(data)

	if err := dss.ws.WriteJSON(dsm); err != nil {
		return err
	}

	return nil
}

func (dss DataStreamSession) ListAll(recvMsg DataStreamMessage, client *dynamic.DynamicClient) error {
	var ds DataStream
	ds.client = client
	ds.recvMsg = recvMsg

	dsm := DataStreamMessage{
		Op: DataStreamOp{
			OpID: recvMsg.Op.OpID,
			Type: WSOpType.listAll,
		},
	}

	var kubeData map[string]unstructured.UnstructuredList = make(map[string]unstructured.UnstructuredList)
	for _, gvr := range ds.recvMsg.Op.Request.KubeGVRKList {
		ds.recvMsg.Op.Request.KubeGVRK = gvr
		_, listData, kubeErr := ds.kubeList()
		if kubeErr != nil && kubeErr.Error() != errCouldNotFindResource.Error() {
			dsm.Error = kubeErr.Error()
		}

		kubeData[gvr.Resource] = listData
	}

	b, err := json.Marshal(kubeData)
	if err != nil {
		return err
	}

	dsm.Data = string(b)

	if err := dss.ws.WriteJSON(dsm); err != nil {
		return err
	}

	return nil
}

func (dss DataStreamSession) Get(recvMsg DataStreamMessage, client *dynamic.DynamicClient) error {
	var ds DataStream
	ds.client = client
	ds.recvMsg = recvMsg

	dsm := DataStreamMessage{
		Op: DataStreamOp{
			OpID: recvMsg.Op.OpID,
			Type: WSOpType.get,
		},
	}

	data, kubeErr := ds.kubeGet()
	if kubeErr != nil {
		dsm.Error = kubeErr.Error()
	}

	dsm.Data = string(data)

	if err := dss.ws.WriteJSON(dsm); err != nil {
		return err
	}

	return nil
}

func (dss DataStreamSession) Update(recvMsg DataStreamMessage, client *dynamic.DynamicClient) error {
	var ds DataStream
	ds.client = client
	ds.recvMsg = recvMsg

	dsm := DataStreamMessage{
		Op: DataStreamOp{
			OpID: recvMsg.Op.OpID,
			Type: WSOpType.update,
		},
	}

	data, kubeErr := ds.kubeUpdate()
	if kubeErr != nil {
		dsm.Error = kubeErr.Error()
	}

	dsm.Data = string(data)

	if err := dss.ws.WriteJSON(dsm); err != nil {
		return err
	}

	return nil
}

func (dss DataStreamSession) Delete(recvMsg DataStreamMessage, client *dynamic.DynamicClient) error {
	var ds DataStream
	ds.client = client
	ds.recvMsg = recvMsg

	dsm := DataStreamMessage{
		Op: DataStreamOp{
			OpID: recvMsg.Op.OpID,
			Type: WSOpType.delete,
		},
	}

	_, kubeErr := ds.kubeDelete()
	if kubeErr != nil {
		dsm.Error = kubeErr.Error()
	}

	if err := dss.ws.WriteJSON(dsm); err != nil {
		return err
	}

	return nil
}

func (dss DataStreamSession) ShowValuesHelm(recvMsg DataStreamMessage, config *rest.Config, h *Helm) error {
	var ds DataStream
	ds.recvMsg = recvMsg
	ds.helm = h

	dsm := DataStreamMessage{
		Op: DataStreamOp{
			OpID: recvMsg.Op.OpID,
			Type: WSOpType.helmShowValues,
		},
	}

	data, helmErr := ds.helmShowValues(config)
	if helmErr != nil {
		dsm.Error = helmErr.Error()
	}

	dsm.Data = data

	if err := dss.ws.WriteJSON(dsm); err != nil {
		return err
	}

	return nil
}

func (dss DataStreamSession) ListHelm(recvMsg DataStreamMessage, config *rest.Config, h *Helm) error {
	var ds DataStream
	ds.recvMsg = recvMsg
	ds.helm = h

	dsm := DataStreamMessage{
		Op: DataStreamOp{
			OpID: recvMsg.Op.OpID,
			Type: WSOpType.helmList,
		},
	}

	data, helmErr := ds.helmList(config)
	if helmErr != nil {
		dsm.Error = helmErr.Error()
	}

	dsm.Data = string(data)

	if err := dss.ws.WriteJSON(dsm); err != nil {
		return err
	}

	return nil
}

func (dss DataStreamSession) GetHelm(recvMsg DataStreamMessage, config *rest.Config, h *Helm) error {
	var ds DataStream
	ds.recvMsg = recvMsg
	ds.helm = h

	dsm := DataStreamMessage{
		Op: DataStreamOp{
			OpID: recvMsg.Op.OpID,
			Type: WSOpType.helmGet,
		},
	}

	data, helmErr := ds.helmGet(config)
	if helmErr != nil {
		dsm.Error = helmErr.Error()
	}

	dsm.Data = string(data)

	if err := dss.ws.WriteJSON(dsm); err != nil {
		return err
	}

	return nil
}

func (dss DataStreamSession) PullHelm(recvMsg DataStreamMessage, config *rest.Config, h *Helm) error {
	var ds DataStream
	ds.recvMsg = recvMsg
	ds.helm = h

	dsm := DataStreamMessage{
		Op: DataStreamOp{
			OpID: recvMsg.Op.OpID,
			Type: WSOpType.helmPull,
		},
	}

	data, helmErr := ds.helmPull(config)
	if helmErr != nil {
		dsm.Error = helmErr.Error()
	}

	dsm.Data = string(data)

	if err := dss.ws.WriteJSON(dsm); err != nil {
		return err
	}

	return nil
}

func (dss DataStreamSession) GetHelmTags(recvMsg DataStreamMessage, config *rest.Config, h *Helm) error {
	var ds DataStream
	ds.recvMsg = recvMsg
	ds.helm = h

	dsm := DataStreamMessage{
		Op: DataStreamOp{
			OpID: recvMsg.Op.OpID,
			Type: WSOpType.helmGetTags,
		},
	}

	data, helmErr := ds.helmChartTags()
	if helmErr != nil {
		dsm.Error = helmErr.Error()
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	dsm.Data = string(b)

	if err := dss.ws.WriteJSON(dsm); err != nil {
		return err
	}

	return nil
}

func (dss DataStreamSession) InstallHelm(recvMsg DataStreamMessage, config *rest.Config, h *Helm) error {
	var ds DataStream
	ds.recvMsg = recvMsg
	ds.helm = h

	dsm := DataStreamMessage{
		Op: DataStreamOp{
			OpID: recvMsg.Op.OpID,
			Type: WSOpType.helmInstall,
		},
	}

	data, helmErr := ds.helmInstall(config, ds.recvMsg.Op.Request.HelmOptions)
	if helmErr != nil {
		dsm.Error = helmErr.Error()
	}

	dsm.Data = string(data)

	if err := dss.ws.WriteJSON(dsm); err != nil {
		return err
	}

	return nil
}

func (dss DataStreamSession) UpgradeHelm(recvMsg DataStreamMessage, config *rest.Config, h *Helm) error {
	var ds DataStream
	ds.recvMsg = recvMsg
	ds.helm = h

	dsm := DataStreamMessage{
		Op: DataStreamOp{
			OpID: recvMsg.Op.OpID,
			Type: WSOpType.helmUpgrade,
		},
	}

	data, helmErr := ds.helmUpgrade(config, ds.recvMsg.Op.Request.HelmOptions)
	if helmErr != nil {
		dsm.Error = helmErr.Error()
	}

	dsm.Data = string(data)

	if err := dss.ws.WriteJSON(dsm); err != nil {
		return err
	}

	return nil
}

func (dss DataStreamSession) UninstallHelm(recvMsg DataStreamMessage, config *rest.Config, h *Helm) error {
	var ds DataStream
	ds.recvMsg = recvMsg
	ds.helm = h

	dsm := DataStreamMessage{
		Op: DataStreamOp{
			OpID: recvMsg.Op.OpID,
			Type: WSOpType.helmUninstall,
		},
	}

	data, err := ds.helmUninstall(config)
	if err != nil {
		dsm.Error = err.Error()
	}

	if dsm.Error == "" {
		dsm.Data = data.Release.Info.Description
	}

	if err := dss.ws.WriteJSON(dsm); err != nil {
		return err
	}

	return nil
}

var dataSessions = DataSessionMap{Sessions: make(map[string]DataStreamSession)}

func handleDataStreamSession(ws *websocket.Conn) {
	var (
		msg DataStreamMessage
		dss DataStreamSession
	)

	if err := ws.ReadJSON(&msg); err != nil {
		fmt.Printf("handleDataStreamSession: can't Recv: %v\n", err)
		return
	}

	if msg.Op.Type != WSOpType.bind {
		fmt.Printf("handleDataStreamSession: expected 'bind' message, got: %s\n", msg.Op.Type)
		return
	}

	if dss = dataSessions.Get(msg.SessionID); dss.id == "" {
		fmt.Printf("handleDataStreamSession: can't find session '%s'\n", msg.SessionID)
		return
	}

	dss.ws = ws
	dataSessions.Set(msg.SessionID, dss)
	dss.bound <- nil

	sendMsg := DataStreamMessage{
		Op: DataStreamOp{
			Type: WSOpType.bind,
		},
	}

	if senderr := dss.ws.WriteJSON(sendMsg); senderr != nil {
		fmt.Println("handleStreamData senderr:", senderr)
		return
	}
}

func startDataStream(ar *APIResource, dss DataStreamSession, closeChan chan struct{}) error {
	for {
		var dsm DataStreamMessage

		err := dss.ws.ReadJSON(&dsm)
		if err != nil {
			defer close(closeChan)
			return err
		}

		switch {
		case dsm.Op.Type == WSOpType.list && dss.id == dsm.SessionID:
			if err := dss.List(dsm, ar.DynamicClient); err != nil {
				return err
			}
		case dsm.Op.Type == WSOpType.listAll && dss.id == dsm.SessionID:
			if err := dss.ListAll(dsm, ar.DynamicClient); err != nil {
				return err
			}
		case dsm.Op.Type == WSOpType.get && dss.id == dsm.SessionID:
			if err := dss.Get(dsm, ar.DynamicClient); err != nil {
				return err
			}
		case dsm.Op.Type == WSOpType.update && dss.id == dsm.SessionID:
			if err := dss.Update(dsm, ar.DynamicClient); err != nil {
				return err
			}
		case dsm.Op.Type == WSOpType.delete && dss.id == dsm.SessionID:
			if err := dss.Delete(dsm, ar.DynamicClient); err != nil {
				return err
			}
		case dsm.Op.Type == WSOpType.helmShowValues && dss.id == dsm.SessionID:
			if err := dss.ShowValuesHelm(dsm, ar.Config, ar.Helm); err != nil {
				return err
			}
		case dsm.Op.Type == WSOpType.helmList && dss.id == dsm.SessionID:
			if err := dss.ListHelm(dsm, ar.Config, ar.Helm); err != nil {
				return err
			}
		case dsm.Op.Type == WSOpType.helmGet && dss.id == dsm.SessionID:
			if err := dss.GetHelm(dsm, ar.Config, ar.Helm); err != nil {
				return err
			}
		case dsm.Op.Type == WSOpType.helmInstall && dss.id == dsm.SessionID:
			if err := dss.InstallHelm(dsm, ar.Config, ar.Helm); err != nil {
				return err
			}
		case dsm.Op.Type == WSOpType.helmUpgrade && dss.id == dsm.SessionID:
			if err := dss.UpgradeHelm(dsm, ar.Config, ar.Helm); err != nil {
				return err
			}
		case dsm.Op.Type == WSOpType.helmPull && dss.id == dsm.SessionID:
			if err := dss.PullHelm(dsm, ar.Config, ar.Helm); err != nil {
				return err
			}
		case dsm.Op.Type == WSOpType.helmGetTags && dss.id == dsm.SessionID:
			if err := dss.GetHelmTags(dsm, ar.Config, ar.Helm); err != nil {
				return err
			}
		case dsm.Op.Type == WSOpType.helmUninstall && dss.id == dsm.SessionID:
			if err := dss.UninstallHelm(dsm, ar.Config, ar.Helm); err != nil {
				return err
			}
		case dsm.Op.Type == WSOpType.close && dss.id == dsm.SessionID:
			defer close(closeChan)
			return nil
		}
	}
}

func WaitForDataStream(ar *APIResource, sessionId string) {
	select {
	case <-dataSessions.Get(sessionId).bound:
		close(dataSessions.Get(sessionId).bound)

		err := startDataStream(ar, dataSessions.Get(sessionId), dataSessions.Get(sessionId).close)
		if err != nil {
			dataSessions.Close(sessionId, WSCloseCode.error, err.Error())
			return
		}

		dataSessions.Close(sessionId, WSCloseCode.info, "Process exited")

	case <-time.After(30 * time.Second):
		close(dataSessions.Get(sessionId).bound)
		delete(dataSessions.Sessions, sessionId)
		return
	}
}
