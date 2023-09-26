package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/igm/sockjs-go/v3/sockjs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type DataStreamSession struct {
	id            string
	bound         chan error
	sockJSSession sockjs.Session
	close         chan struct{}
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

func (dsm *DataSessionMap) Close(sessionId string, status uint32, reason string) {
	dsm.Lock.Lock()
	defer dsm.Lock.Unlock()
	ses := dsm.Sessions[sessionId]
	ses.sockJSSession.Close(status, reason)
	delete(dsm.Sessions, sessionId)
}

type DataStream struct {
	client  *dynamic.DynamicClient
	recvMsg DataStreamMessage
	helm    *Helm
}

func (ds DataStream) kubeList() ([]byte, error) {
	var lr ListResource
	lr.GVR = schema.GroupVersionResource{
		Group:    ds.recvMsg.Op.Request.Group,
		Version:  ds.recvMsg.Op.Request.Version,
		Resource: ds.recvMsg.Op.Request.Resource,
	}
	lr.Namespace = ds.recvMsg.Op.Request.Namespace

	sel, err := fields.ParseSelector(ds.recvMsg.Op.Request.ResourceOptions.FieldSelector)
	if err != nil {
		return nil, err
	}

	lr.Options = metav1.ListOptions{
		FieldSelector: sel.String(),
		LabelSelector: ds.recvMsg.Op.Request.ResourceOptions.LabelSelector,
		Limit:         ds.recvMsg.Op.Request.ResourceOptions.Limit,
		Continue:      ds.recvMsg.Op.Request.ResourceOptions.Continue,
	}

	list, err := lr.List(ds.client)
	if err != nil {
		return nil, err
	}

	dataBytes, err := list.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return dataBytes, nil
}

func (ds DataStream) kubeGet() ([]byte, error) {
	var gr GetResource
	gr.Namespace = ds.recvMsg.Op.Request.Namespace
	gr.Name = ds.recvMsg.Op.Request.Name
	gr.Options = metav1.GetOptions{}
	gr.GVR = schema.GroupVersionResource{
		Group:    ds.recvMsg.Op.Request.Group,
		Version:  ds.recvMsg.Op.Request.Version,
		Resource: ds.recvMsg.Op.Request.Resource,
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
		Group:    ds.recvMsg.Op.Request.Group,
		Version:  ds.recvMsg.Op.Request.Version,
		Resource: ds.recvMsg.Op.Request.Resource,
	}
	ur.GVK = &schema.GroupVersionKind{
		Group:   ds.recvMsg.Op.Request.Group,
		Version: ds.recvMsg.Op.Request.Version,
		Kind:    ds.recvMsg.Op.Request.Kind,
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
		Group:    ds.recvMsg.Op.Request.Group,
		Version:  ds.recvMsg.Op.Request.Version,
		Resource: ds.recvMsg.Op.Request.Resource,
	}

	err := dr.Delete(ds.client)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds DataStream) helmList() ([]byte, error) {
	var h Helm
	h.ActionConfig = ds.helm.ActionConfig
	h.EnvSettings = ds.helm.EnvSettings
	h.ReleaseNamespace = ds.recvMsg.Op.Request.Namespace

	rl, err := h.ListReleases()
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(rl)
	if err != nil {
		return nil, err
	}

	return dataBytes, nil
}

func (ds DataStream) helmGet() ([]byte, error) {
	var h Helm
	h.ActionConfig = ds.helm.ActionConfig
	h.EnvSettings = ds.helm.EnvSettings
	h.AllNamespaces = false
	h.ReleaseName = ds.recvMsg.Op.Request.Name
	h.ReleaseNamespace = ds.recvMsg.Op.Request.Namespace
	h.AllValues = true

	rl, err := h.GetValues()
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(rl)
	if err != nil {
		return nil, err
	}

	return dataBytes, nil
}

func (dss DataStreamSession) List(recvMsg DataStreamMessage, client *dynamic.DynamicClient) error {
	var ds DataStream
	ds.client = client
	ds.recvMsg = recvMsg

	dsm := DataStreamMessage{
		Op: DataStreamOp{
			Type: WSOpType.list,
		},
	}

	data, kubeErr := ds.kubeList()
	if kubeErr != nil {
		dsm.Error = kubeErr.Error()
	}

	dsm.Data = string(data)

	msg, err := json.Marshal(dsm)
	if err != nil {
		return err
	}

	if err := dss.sockJSSession.Send(string(msg)); err != nil {
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
			Type: WSOpType.get,
		},
	}

	data, kubeErr := ds.kubeGet()
	if kubeErr != nil {
		dsm.Error = kubeErr.Error()
	}

	dsm.Data = string(data)

	msg, err := json.Marshal(dsm)
	if err != nil {
		return err
	}

	if err := dss.sockJSSession.Send(string(msg)); err != nil {
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
			Type: WSOpType.update,
		},
	}

	data, kubeErr := ds.kubeUpdate()
	if kubeErr != nil {
		dsm.Error = kubeErr.Error()
	}

	dsm.Data = string(data)

	msg, err := json.Marshal(dsm)
	if err != nil {
		return err
	}

	if err := dss.sockJSSession.Send(string(msg)); err != nil {
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
			Type: WSOpType.delete,
		},
	}

	_, kubeErr := ds.kubeDelete()
	if kubeErr != nil {
		dsm.Error = kubeErr.Error()
	}

	msg, err := json.Marshal(dsm)
	if err != nil {
		return err
	}

	if err := dss.sockJSSession.Send(string(msg)); err != nil {
		return err
	}

	return nil
}

func (dss DataStreamSession) ListHelm(recvMsg DataStreamMessage, h *Helm) error {
	var ds DataStream
	ds.recvMsg = recvMsg
	ds.helm = h

	dsm := DataStreamMessage{
		Op: DataStreamOp{
			Type: WSOpType.list,
		},
	}

	data, helmErr := ds.helmList()
	if helmErr != nil {
		dsm.Error = helmErr.Error()
	}

	dsm.Data = string(data)

	msg, err := json.Marshal(dsm)
	if err != nil {
		return err
	}

	if err := dss.sockJSSession.Send(string(msg)); err != nil {
		return err
	}

	return nil
}

func (dss DataStreamSession) GetHelm(recvMsg DataStreamMessage, h *Helm) error {
	var ds DataStream
	ds.recvMsg = recvMsg
	ds.helm = h

	dsm := DataStreamMessage{
		Op: DataStreamOp{
			Type: WSOpType.get,
		},
	}

	data, helmErr := ds.helmGet()
	if helmErr != nil {
		dsm.Error = helmErr.Error()
	}

	dsm.Data = string(data)

	msg, err := json.Marshal(dsm)
	if err != nil {
		return err
	}

	if err := dss.sockJSSession.Send(string(msg)); err != nil {
		return err
	}

	return nil
}

var dataSessions = DataSessionMap{Sessions: make(map[string]DataStreamSession)}

func handleDataStreamSession(session sockjs.Session) {
	var (
		buf               string
		err               error
		msg               DataStreamMessage
		dataStreamSession DataStreamSession
	)

	if session.GetSessionState() != sockjs.SessionActive {
		fmt.Println("getsessionstate:", session.GetSessionState())
		return
	}

	if buf, err = session.Recv(); err != nil {
		fmt.Printf("handleDataStreamSession: can't Recv: %v\n", err)
		return
	}

	if err = json.Unmarshal([]byte(buf), &msg); err != nil {
		fmt.Printf("handleDataStreamSession: can't UnMarshal (%v): %s\n", err, buf)
		return
	}

	if msg.Op.Type != WSOpType.bind {
		fmt.Printf("handleDataStreamSession: expected 'bind' message, got: %s\n", buf)
		return
	}

	if dataStreamSession = dataSessions.Get(msg.SessionID); dataStreamSession.id == "" {
		fmt.Printf("handleDataStreamSession: can't find session '%s'\n", msg.SessionID)
		return
	}

	dataStreamSession.sockJSSession = session
	dataSessions.Set(msg.SessionID, dataStreamSession)
	dataStreamSession.bound <- nil

	sendMsg, senderr := json.Marshal(DataStreamMessage{
		Op: DataStreamOp{
			Type: WSOpType.bind,
		},
	})
	if senderr != nil {
		fmt.Println("handleStreamData senderr:", senderr)
		return
	}

	if senderr := session.Send(string(sendMsg)); senderr != nil {
		fmt.Println("handleStreamData senderr:", senderr)
		return
	}
}

func startDataStream(ar *APIResource, dss DataStreamSession, close chan struct{}) error {
	for {
		m, err := dss.sockJSSession.Recv()
		if dss.sockJSSession.GetSessionState() != sockjs.SessionActive {
			return nil
		}
		if err != nil {
			fmt.Println(err)
			return err
		}

		var msg DataStreamMessage
		if err := json.Unmarshal([]byte(m), &msg); err != nil {
			return err
		}

		switch msg.Op.Type {
		case WSOpType.list:
			if err := dss.List(msg, ar.DynamicClient); err != nil {
				return err
			}
		case WSOpType.get:
			if err := dss.Get(msg, ar.DynamicClient); err != nil {
				return err
			}
		case WSOpType.update:
			if err := dss.Update(msg, ar.DynamicClient); err != nil {
				return err
			}
		case WSOpType.delete:
			if err := dss.Delete(msg, ar.DynamicClient); err != nil {
				return err
			}
		case WSOpType.helmList:
			if err := dss.ListHelm(msg, ar.Helm); err != nil {
				return err
			}
		case WSOpType.helmGet:
			if err := dss.GetHelm(msg, ar.Helm); err != nil {
				return err
			}
		case WSOpType.close:
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

	case <-time.After(10 * time.Second):
		close(dataSessions.Get(sessionId).bound)
		delete(dataSessions.Sessions, sessionId)
		return
	}
}
