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

func (lss LogStreamSession) Read(p []byte) (int, error) {
	m, err := lss.sockJSSession.Recv()
	if err != nil {
		return copy(p, END_OF_TRANSMISSION), err
	}

	var msg LogStreamMessage
	if err := json.Unmarshal([]byte(m), &msg); err != nil {
		return copy(p, END_OF_TRANSMISSION), err
	}

	switch msg.Op {
	case WSOpType.stdin:
		return copy(p, msg.Data), nil
	case WSOpType.close:
		close(lss.close)
		return copy(p, END_OF_TRANSMISSION), nil
	default:
		return copy(p, END_OF_TRANSMISSION), fmt.Errorf("unknown message type '%s'", msg.Op)
	}
}

func (lss LogStreamSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(LogStreamMessage{
		Op:   WSOpType.stdout,
		Data: string(p),
	})
	if err != nil {
		return 0, err
	}

	if err := lss.sockJSSession.Send(string(msg)); err != nil {
		return 0, err
	}
	return len(p), nil
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
		buf              string
		err              error
		msg              LogStreamMessage
		logStreamSession LogStreamSession
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

	if logStreamSession = logStreamSessions.Get(msg.SessionID); logStreamSession.id == "" {
		fmt.Printf("handleLogStreamSession: can't find session '%s'\n", msg.SessionID)
		return
	}

	logStreamSession.sockJSSession = session
	logStreamSessions.Set(msg.SessionID, logStreamSession)
	logStreamSession.bound <- nil
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

func readLines(response io.ReadCloser, done chan struct{}, lsh LogStreamHandler) error {
	cancel := make(chan struct{})
	go func() {
		if _, err := lsh.Read([]byte{}); err != nil {
			return
		}

		select {
		case <-done:
			response.Close()
		case <-cancel:
			return
		}
	}()

	defer response.Close()
	defer close(cancel) // ensure that goroutine exits

	reader := bufio.NewReader(response)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}
		lsh.Write(line)
	}

}

func startLogsStream(client kubernetes.Interface, lsh LogStreamHandler, close chan struct{}, pld *PodLogsData) error {
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

	errRead := readLines(out, close, lsh)
	if errRead != nil {
		return errRead
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

		logStreamSessions.Close(sessionId, WSCloseCode.info, "Process exited")

	case <-time.After(30 * time.Second):
		close(logStreamSessions.Get(sessionId).bound)
		delete(logStreamSessions.Sessions, sessionId)
		return
	}
}
