package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
)

// PtyHandler is what remotecommand expects from a pty
type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}

// TerminalSession implements PtyHandler (using a websocket connection)
type TerminalSession struct {
	id       string
	bound    chan error
	ws       *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
}

// TerminalMessage is the messaging protocol between ShellController and TerminalSession.
//
// OP      DIRECTION  FIELD(S) USED  DESCRIPTION
// ---------------------------------------------------------------------
// bind    fe->be     SessionID      Id sent back from TerminalResponse
// stdin   fe->be     Data           Keystrokes/paste buffer
// resize  fe->be     Rows, Cols     New terminal size
// stdout  be->fe     Data           Output from the process
// toast   be->fe     Data           OOB message to be shown to the user
type TerminalMessage struct {
	Op         string `json:"op"`
	Data       string `json:"data"`
	SessionID  string `json:"sessionId"`
	Rows       uint16 `json:"rows"`
	Cols       uint16 `json:"cols"`
	Error      string `json:"error"`
	StatusCode uint   `json:"statusCode"`
}

// Next handles pty->process resize events
// Called in a loop from remotecommand as long as the process is running
func (t TerminalSession) Next() *remotecommand.TerminalSize {
	size := <-t.sizeChan
	if size.Height == 0 && size.Width == 0 {
		return nil
	}
	return &size
}

// Read handles pty->process messages (stdin, resize)
// Called in a loop from remotecommand as long as the process is running
func (t TerminalSession) Read(p []byte) (int, error) {
	var tm TerminalMessage

	err := t.ws.ReadJSON(&tm)
	if err != nil {
		if (terminalSessions.Get(t.id) == TerminalSession{}) {
			return copy(p, END_OF_TRANSMISSION), io.EOF
		}
		// Send terminated signal to process to avoid resource leak
		return copy(p, END_OF_TRANSMISSION), err
	}

	switch tm.Op {
	case WSOpType.close:
		return copy(p, END_OF_TRANSMISSION), io.EOF
	case WSOpType.stdin:
		return copy(p, tm.Data), nil
	case WSOpType.resize:
		t.sizeChan <- remotecommand.TerminalSize{Width: tm.Cols, Height: tm.Rows}
		return 0, nil
	default:
		return copy(p, END_OF_TRANSMISSION), fmt.Errorf("unknown message type '%s'", tm.Op)
	}
}

// Write handles process->pty stdout
// Called from remotecommand whenever there is any output
func (t TerminalSession) Write(p []byte) (int, error) {
	tm := TerminalMessage{
		Op:   WSOpType.stdout,
		Data: string(p),
	}

	if err := t.ws.WriteJSON(tm); err != nil {
		return 0, err
	}

	return len(p), nil
}

// Toast can be used to send the user any OOB messages
// hterm puts these in the center of the terminal
func (t TerminalSession) Toast(p string) error {
	tm := TerminalMessage{
		Op:   WSOpType.toast,
		Data: p,
	}

	if err := t.ws.WriteJSON(tm); err != nil {
		return err
	}

	return nil
}

// SessionMap stores a map of all TerminalSession objects and a lock to avoid concurrent conflict
type TerminalSessionMap struct {
	Sessions map[string]TerminalSession
	Lock     sync.RWMutex
}

// Get return a given terminalSession by sessionId
func (tsm *TerminalSessionMap) Get(sessionId string) TerminalSession {
	tsm.Lock.RLock()
	defer tsm.Lock.RUnlock()
	return tsm.Sessions[sessionId]
}

// Set store a TerminalSession to SessionMap
func (tsm *TerminalSessionMap) Set(sessionId string, session TerminalSession) {
	tsm.Lock.Lock()
	defer tsm.Lock.Unlock()
	tsm.Sessions[sessionId] = session
}

// Close shuts down the websocket connection and sends the status code and reason to the client
// Can happen if the process exits or if there is an error starting up the process
// For now the status code is unused and reason is shown to the user (unless "")
func (tsm *TerminalSessionMap) Close(sessionId string, status uint, reason string) {
	tsm.Lock.Lock()
	defer tsm.Lock.Unlock()
	ses := tsm.Sessions[sessionId]
	ses.ws.WriteJSON(TerminalMessage{Op: "close", Data: reason, StatusCode: status})
	ses.ws.Close()
	close(ses.sizeChan)
	delete(tsm.Sessions, sessionId)
}

var terminalSessions = TerminalSessionMap{Sessions: make(map[string]TerminalSession)}

// handleTerminalSession is Called by net/http for any new /srv/shell connections
func handleTerminalSession(ws *websocket.Conn) {
	var (
		tm TerminalMessage
		ts TerminalSession
	)

	if err := ws.ReadJSON(&tm); err != nil {
		fmt.Printf("handleTerminalSession: can't Recv: %v\n", err)
		return
	}

	if tm.Op != WSOpType.bind {
		fmt.Printf("handleTerminalSession: expected 'bind' message, got: '%s'\n", tm.Op)
		return
	}

	if ts = terminalSessions.Get(tm.SessionID); ts.id == "" {
		ws.Close()
		return
	}

	ts.ws = ws
	terminalSessions.Set(tm.SessionID, ts)
	ts.bound <- nil

	sendMsg := TerminalMessage{
		Op:        WSOpType.bind,
		SessionID: tm.SessionID,
	}

	if senderr := ts.ws.WriteJSON(sendMsg); senderr != nil {
		fmt.Println("handleTerminalSession senderr:", senderr)
		return
	}
}

// Executed cmd in the container specified in request and connects it up with the ptyHandler (a session)
func startProcess(k8sClient kubernetes.Interface, cfg *rest.Config, ed *ExecData, ptyHandler PtyHandler) error {
	req := k8sClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(ed.PodName).
		Namespace(ed.PodNamespace).
		SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Container: ed.ContainerName,
		Command:   ed.Command,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
	if err != nil {
		return err
	}

	if err := exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdin:             ptyHandler,
		Stdout:            ptyHandler,
		Stderr:            ptyHandler,
		TerminalSizeQueue: ptyHandler,
		Tty:               true,
	}); err != nil {
		return err
	}

	return nil
}

type ExecData struct {
	Command       []string
	PodName       string
	PodNamespace  string
	ContainerName string
	Shell         string
	GVR           schema.GroupVersionResource
	GVK           *schema.GroupVersionKind
	remotecommand.TerminalSizeQueue
}

// execute a single command without attaching to a TTY
func (ed *ExecData) executeRemoteCommand(k8sClient kubernetes.Interface, cfg *rest.Config, command []string) (string, error) {
	request := k8sClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(ed.PodName).
		Namespace(ed.PodNamespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: ed.ContainerName,
			Command:   command,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(cfg, "POST", request.URL())
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdout: buf,
		Stderr: errBuf,
	})

	switch {
	case errBuf.String() != "":
		return "", fmt.Errorf("cannot execute command '%s': '%s'", command, errBuf.String())
	case err != nil:
		return "", fmt.Errorf("cannot execute command '%s': %s", command, err)
	}

	return buf.String(), nil
}

// WaitForTerminal is called from apihandler.handleAttach as a goroutine
// Waits for the websocket connection to be opened by the client the session to be bound in handleTerminalSession
func (ed *ExecData) WaitForTerminal(k8sClient kubernetes.Interface, cfg *rest.Config, sessionId string) {
	select {
	case <-terminalSessions.Get(sessionId).bound:
		close(terminalSessions.Get(sessionId).bound)

		err := startProcess(k8sClient, cfg, ed, terminalSessions.Get(sessionId))
		if err != nil {
			terminalSessions.Close(sessionId, WSCloseCode.error, err.Error())
			return
		}

		terminalSessions.Close(sessionId, WSCloseCode.info, "Process exited")

	case <-time.After(30 * time.Second):
		close(terminalSessions.Get(sessionId).bound)
		delete(terminalSessions.Sessions, sessionId)
		return
	}
}
