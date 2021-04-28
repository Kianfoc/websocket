/*
Websocket package for Resato
*/
package websocket

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"time"
)

//todo
//Error log
//Check error json

type Ws struct {
	// the websocket connection
	conn *websocket.Conn

	// certificate pool used for secure connections
	caPool *x509.CertPool

	// set to true to use a certificate
	secure bool

	// url contains the url to connect to
	url url.URL

	// set to true to automatically reconnect
	reconnect bool

	// set to true when there is a error
	websocketErr bool

	// Array for messages that didn't send because of a error
	messages []historyMessage

	// close handler is called when a connection ends
	closeHandler func(int, string) error
}

type historyMessage struct {
	messageByte []byte
	messageJson interface{}
	messageType string
	time        int64
}

var Websocket Ws

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Create websocket connection
func (w *Ws) CreateConnection(rW http.ResponseWriter, r *http.Request) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	conn, err := upgrader.Upgrade(rW, r, nil)
	if err != nil {
		return err
	}
	w.conn = conn
	return nil
}

// Tries to connect to a websocket server
func (w *Ws) Connect() error {
	err := w.dial()
	if err != nil {
		return err
	}
	if w.reconnect {
		go w.CheckConnection()
		go w.autoReconnect()
	}
	return nil
}

func (w *Ws) dial() error {
	var d websocket.Dialer
	if w.secure {
		config := tls.Config{RootCAs: w.caPool}
		d = websocket.Dialer{TLSClientConfig: &config}
	}
	conn, _, err := d.Dial(w.url.String(), nil)
	if err != nil {
		return err
	}
	w.conn = conn
	return nil
}

func (w *Ws) Reconnect(b bool) {
	w.reconnect = b
}

func (w *Ws) autoReconnect() {
	for {
		time.Sleep(10 * time.Second)
		if !w.websocketErr {
			continue
		}
		fmt.Println("hallo1")
		err := w.dial()
		if err != nil {
			fmt.Println(err)
			//todo read error
			continue
		}
		//Connection succeeded
		//w.WriteHistory()
		fmt.Println("hallo2")
		w.websocketErr = false
	}
}

// Sends every 60 seconds a ping message to check if there is a connection
func (w *Ws) CheckConnection() {
	for range time.Tick(60 * time.Second) {
		err := w.conn.WriteMessage(websocket.PingMessage, []byte{})
		if err != nil {
			w.websocketErr = true
		}
	}
}

// Appends a cert to pool
func (w *Ws) AppendCert(cert []byte) {
	w.caPool.AppendCertsFromPEM(cert)
}

// Set the url to connect to
// Auto choose scheme
func (w *Ws) SetUrl(host, path string) {
	var scheme string
	if w.secure {
		scheme = "wss"
	} else {
		scheme = "ws"
	}
	w.url = url.URL{Scheme: scheme, Host: host, Path: path}
}

// set secure websocket connection
// create a new certificate pool
func (w *Ws) SetSecure(b bool) {
	w.secure = b
	if w.secure {
		w.caPool = x509.NewCertPool()
	}
}

// Write a message
func (w *Ws) WriteMessage(messageType int, data []byte) error {
	err := w.conn.WriteMessage(messageType, data)
	if err != nil {
		w.websocketErr = true
		w.messages = append(w.messages, historyMessage{
			messageByte: data,
			messageType: "msg",
			time:        time.Now().Unix(),
		})
		return err
	}
	return err
}

//todo fix errorcheck

//Write a message in json format
func (w *Ws) WriteJSON(data interface{}) error {
	err := w.conn.WriteJSON(data)
	if err != nil {
		w.websocketErr = true
		w.messages = append(w.messages, historyMessage{
			messageJson: data,
			messageType: "json",
			time:        time.Now().Unix(),
		})
		return err
	}
	return err
}

func (w *Ws) WriteHistory() {
	for k, v := range w.messages {
		if time.Now().Unix()-v.time >= 5 {
			w.messages = append(w.messages[:k], w.messages[k+1:]...)
			continue
		}
		if v.messageType == "msg" {
			w.WriteMessage(1, v.messageByte)
		} else if v.messageType == "json" {
			w.WriteJSON(v.messageJson)
		}
	}
	w.messages = []historyMessage{}
}

// Read a websocket message
func (w *Ws) Read() (int, []byte, error) {
	t, d, err := w.conn.ReadMessage()
	if err != nil {
		w.websocketErr = true
		return 0, []byte{}, err
	}
	return t, d, nil
}

//todo fix errorcheck

// Read a websocket message in json format
func (w *Ws) ReadJSON(v interface{}) error {
	err := w.conn.ReadJSON(v)
	if err != nil {
		w.websocketErr = true
		return err
	}
	return err
}

// Close websocket connection with a message
func (w *Ws) CloseConnection(code int, text string) error {
	err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, text))
	if err != nil {
		return err
	}
	return nil
}

//todo

//Channel for errors
func (w *Ws) Error(e chan error) {
	err := errors.New("testError")
	e <- err
}
