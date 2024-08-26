package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

type LogStreamSession struct {
	id    string
	bound chan error
	ws    *websocket.Conn
	close chan struct{}
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

func (lssm *LogStreamSessionMap) Close(sessionId string, status uint, reason string) {
	lssm.Lock.Lock()
	defer lssm.Lock.Unlock()
	ses := lssm.Sessions[sessionId]
	ses.ws.WriteJSON(LogStreamMessage{Op: "close", Data: reason, StatusCode: status})
	ses.ws.Close()
	delete(lssm.Sessions, sessionId)
}

var logStreamSessions = LogStreamSessionMap{Sessions: make(map[string]LogStreamSession)}

func handleLogStreamSession(ws *websocket.Conn) {
	var (
		lsm LogStreamMessage
		lss LogStreamSession
	)

	if err := ws.ReadJSON(&lsm); err != nil {
		fmt.Printf("handleLogStreamSession: can't Recv: %v\n", err)
		return
	}

	if lsm.Op != WSOpType.bind {
		fmt.Printf("handleLogStreamSession: expected 'bind' message, got: '%s'\n", lsm.Op)
		return
	}

	if lss = logStreamSessions.Get(lsm.SessionID); lss.id == "" {
		fmt.Printf("handleLogStreamSession: can't find session '%s'\n", lsm.SessionID)
		return
	}

	lss.ws = ws
	logStreamSessions.Set(lsm.SessionID, lss)
	lss.bound <- nil

	sendMsg := LogStreamMessage{
		Op:        WSOpType.bind,
		SessionID: lsm.SessionID,
	}

	if senderr := ws.WriteJSON(sendMsg); senderr != nil {
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

	go func() {
		out, err := req.Stream(context.TODO())
		if err != nil {
			return
		}
		defer out.Close()

		reader := bufio.NewReader(out)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					lss.ws.WriteJSON(LogStreamMessage{
						Op:         "close",
						Data:       "Log stream ended",
						StatusCode: WSCloseCode.info,
					})
					lss.ws.Close()
				}
				return
			}
			lsm := LogStreamMessage{
				Op:        WSOpType.stdout,
				SessionID: lss.id,
				Data:      string(line),
			}

			if err := lss.ws.WriteJSON(lsm); err != nil {
				return
			}
		}
	}()

	return nil
}

type PodLogsData struct {
	Namespace      string
	Name           string
	ResourceType   string
	ContainerName  string
	Follow         string
	TailLines      string
	Options        *corev1.PodLogOptions
	ParameterCodec runtime.ParameterCodec
}

func startLogsStream(client kubernetes.Interface, lss LogStreamSession, closeChan chan struct{}, pld *PodLogsData) error {
	for {
		var lsm LogStreamMessage

		err := lss.ws.ReadJSON(&lsm)
		if err != nil {
			defer close(closeChan)
			if (logStreamSessions.Get(lss.id) == LogStreamSession{}) {
				return io.EOF
			}
			return err
		}

		switch {
		case lsm.Op == WSOpType.stdin && lss.id == lsm.SessionID:
			if err := lss.SendLogs(client, pld); err != nil {
				return err
			}
		case lsm.Op == WSOpType.close && lss.id == lsm.SessionID:
			defer close(closeChan)
			return nil
		}
	}
}

func (pld *PodLogsData) WaitForLogs(client kubernetes.Interface, sessionId string) {
	select {
	case <-logStreamSessions.Get(sessionId).bound:
		close(logStreamSessions.Get(sessionId).bound)

		err := startLogsStream(client, logStreamSessions.Get(sessionId), logStreamSessions.Get(sessionId).close, pld)
		if err != nil {
			if err == io.EOF {
				logStreamSessions.Close(sessionId, WSCloseCode.warning, "Log stream ended")
				return
			}
			logStreamSessions.Close(sessionId, WSCloseCode.error, err.Error())
			return
		}

		logStreamSessions.Close(sessionId, WSCloseCode.info, "Log stream ended")

	case <-time.After(30 * time.Second):
		close(logStreamSessions.Get(sessionId).bound)
		delete(logStreamSessions.Sessions, sessionId)
		return
	}
}
