package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/igm/sockjs-go/v3/sockjs"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

type LogStreamHandler interface {
	io.Reader
	io.Writer
}

type LogStreamSession struct {
	id            string
	bound         chan error
	sockJSSession sockjs.Session
	close         chan struct{}
}

type LogStreamMessage struct {
	Op         string `json:"op"`
	Data       string `json:"data"`
	SessionID  string `json:"sessionId"`
	Error      string `json:"error"`
	StatusCode uint   `json:"statusCode"`
}

type LogStreamSessionMap struct {
	Sessions map[string]LogStreamSession
	Lock     sync.RWMutex
}

func (lssm *LogStreamSessionMap) Get(sessionId string) LogStreamSession {
	lssm.Lock.RLock()
	defer lssm.Lock.RUnlock()
	return lssm.Sessions[sessionId]
}

func (lssm *LogStreamSessionMap) Set(sessionId string, session LogStreamSession) {
	lssm.Lock.Lock()
	defer lssm.Lock.Unlock()
	lssm.Sessions[sessionId] = session
}

func (lssm *LogStreamSessionMap) Close(sessionId string, status uint32, reason string) {
	lssm.Lock.Lock()
	defer lssm.Lock.Unlock()
	ses := lssm.Sessions[sessionId]
	ses.sockJSSession.Close(status, reason)
	delete(lssm.Sessions, sessionId)
}

var logStreamSessions = LogStreamSessionMap{Sessions: make(map[string]LogStreamSession)}

func handleLogStreamSession(session sockjs.Session) {
	var (
		buf string
		err error
		msg LogStreamMessage
		lss LogStreamSession
	)

	if session.GetSessionState() != sockjs.SessionActive {
		return
	}

	if buf, err = session.Recv(); err != nil {
		fmt.Printf("handleLogStreamSession: can't Recv: %v\n", err)
		return
	}

	if err = json.Unmarshal([]byte(buf), &msg); err != nil {
		fmt.Printf("handleLogStreamSession: can't UnMarshal (%v): %s\n", err, buf)
		return
	}

	if msg.Op != WSOpType.bind {
		fmt.Printf("handleLogStreamSession: expected 'bind' message, got: '%s'\n", msg.Op)
		return
	}

	if lss = logStreamSessions.Get(msg.SessionID); lss.id == "" {
		fmt.Printf("handleLogStreamSession: can't find session '%s'\n", msg.SessionID)
		return
	}

	lss.sockJSSession = session
	logStreamSessions.Set(msg.SessionID, lss)
	lss.bound <- nil

	sendMsg, senderr := json.Marshal(LogStreamMessage{
		Op:        WSOpType.bind,
		SessionID: msg.SessionID,
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

func (lss LogStreamSession) SendLogs(client kubernetes.Interface, pld *PodLogsData) error {
	req := client.CoreV1().RESTClient().Get().
		Namespace(pld.Namespace).
		Name(pld.Name).
		Resource(pld.ResourceType).
		SubResource("log").
		Param("follow", pld.Follow).
		Param("container", pld.ContainerName)

	if pld.TailLines != "" {
		req.Param("tailLines", pld.TailLines)
	}

	out, err := req.Stream(context.TODO())
	if err != nil {
		return err
	}
	defer out.Close()

	reader := bufio.NewReader(out)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}
		msg, err := json.Marshal(LogStreamMessage{
			Op:        WSOpType.stdout,
			SessionID: lss.id,
			Data:      string(line),
		})
		if err != nil {
			return err
		}

		lss.sockJSSession.Send(string(msg))
	}
}

type PodLogsData struct {
	Namespace      string
	Name           string
	ResourceType   string
	ContainerName  string
	Follow         string
	TailLines      string
	Options        *v1.PodLogOptions
	ParameterCodec runtime.ParameterCodec
}

func startLogsStream(client kubernetes.Interface, lss LogStreamSession, closeChan chan struct{}, pld *PodLogsData) error {
	m, err := lss.sockJSSession.Recv()
	if lss.sockJSSession.GetSessionState() != sockjs.SessionActive {
		defer close(closeChan)
		return nil
	}
	if err != nil {
		defer close(closeChan)
		return err
	}

	var msg LogStreamMessage
	if err := json.Unmarshal([]byte(m), &msg); err != nil {
		return err
	}

	switch {
	case msg.Op == WSOpType.stdin && lss.id == msg.SessionID:
		if err := lss.SendLogs(client, pld); err != nil {
			return err
		}
	case msg.Op == WSOpType.close && lss.id == msg.SessionID:
		defer close(closeChan)
		return nil
	}

	return nil
}

func (pld *PodLogsData) WaitForLogs(client kubernetes.Interface, sessionId string) {
	select {
	case <-logStreamSessions.Get(sessionId).bound:
		close(logStreamSessions.Get(sessionId).bound)

		err := startLogsStream(client, logStreamSessions.Get(sessionId), logStreamSessions.Get(sessionId).close, pld)
		if err != nil {
			logStreamSessions.Close(sessionId, WSCloseCode.error, err.Error())
			return
		}

		logStreamSessions.Close(sessionId, WSCloseCode.info, "process exited")

	case <-time.After(30 * time.Second):
		close(logStreamSessions.Get(sessionId).bound)
		delete(logStreamSessions.Sessions, sessionId)
		return
	}
}
