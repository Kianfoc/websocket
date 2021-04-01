package websocket

import (
    "crypto/x509"
    "encoding/json"
    "github.com/gorilla/websocket"
    "io/ioutil"
    "net/http"
    "net/url"
    "reflect"
    "testing"
    "time"
)

//Does not work
//func TestWs_CreateConnection(t *testing.T) {
//    mux := http.NewServeMux()
//    srv := &http.Server{Addr: ":2758", Handler: mux}
//    go func() {
//        var upgrader = websocket.Upgrader{
//            ReadBufferSize:  1024,
//            WriteBufferSize: 1024,
//        }
//        mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
//            upgrader.CheckOrigin = func(r *http.Request) bool { return true }
//            connection, _ := upgrader.Upgrade(w, r, nil)
//            defer connection.Close()
//        })
//        _ = srv.ListenAndServe()
//    }()
//    req := httptest.NewRequest("GET", "ws://localhost:2758/ws", nil)
//    req.Header.Set("Connection", "Upgrade")
//    req.Header.Set("Upgrade", "websocket")
//    req.Header.Set("Sec-Websocket-Version", "13")
//    req.Header.Set("Sec-Websocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
//    w := httptest.NewRecorder()
//    type fields struct {
//        conn         *websocket.Conn
//        caPool       *x509.CertPool
//        secure       bool
//        url          url.URL
//        reconnect    bool
//        websocketErr bool
//        closeHandler func(int, string) error
//    }
//    type args struct {
//        rW http.ResponseWriter
//        r  *http.Request
//    }
//    tests := []struct {
//        name    string
//        fields  fields
//        args    args
//        wantErr bool
//    }{
//        {"CreateConnection 1", fields{}, args{rW: w, r: req}, false},
//    }
//    for _, tt := range tests {
//        t.Run(tt.name, func(t *testing.T) {
//            w := &Ws{
//                conn:         tt.fields.conn,
//                caPool:       tt.fields.caPool,
//                secure:       tt.fields.secure,
//                url:          tt.fields.url,
//                reconnect:    tt.fields.reconnect,
//                websocketErr: tt.fields.websocketErr,
//                closeHandler: tt.fields.closeHandler,
//            }
//            if err := w.CreateConnection(tt.args.rW, tt.args.r); (err != nil) != tt.wantErr {
//                t.Errorf("CreateConnection() error = %v, wantErr %v", err, tt.wantErr)
//            }
//        })
//    }
//}

func TestWs_Connect(t *testing.T) {
    var ws = Websocket
    ws.SetUrl("localhost:2758", "/ws")
    var wsErr = Websocket
    wsErr.SetUrl("localhost:8080", "/ws")

    mux := http.NewServeMux()
    srv := &http.Server{Addr: ":2758", Handler: mux}
    // Starts a server to connect
    go func() {
        var upgrader = websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
        }
        mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
            upgrader.CheckOrigin = func(r *http.Request) bool { return true }
            connection, _ := upgrader.Upgrade(w, r, nil)
            defer connection.Close()
        })
        _ = srv.ListenAndServe()
    }()
    type fields struct {
        conn         *websocket.Conn
        caPool       *x509.CertPool
        secure       bool
        url          url.URL
        reconnect    bool
        websocketErr bool
        closeHandler func(int, string) error
    }
    tests := []struct {
        name    string
        fields  fields
        wantErr bool
    }{
        {"Connect 1", fields{url: ws.url}, false},
        {"Connect 2", fields{url: wsErr.url}, true},
        {"Connect 3", fields{url: ws.url, secure: true}, false},
        {"Connect 4", fields{url: ws.url, reconnect: true}, false},
        {"Connect 5", fields{url: ws.url, reconnect: true, secure: true}, false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := &Ws{
                conn:         tt.fields.conn,
                caPool:       tt.fields.caPool,
                secure:       tt.fields.secure,
                url:          tt.fields.url,
                reconnect:    tt.fields.reconnect,
                websocketErr: tt.fields.websocketErr,
                closeHandler: tt.fields.closeHandler,
            }
            if err := w.Connect(); (err != nil) != tt.wantErr {
                t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestWs_Reconnect(t *testing.T) {
    var ws = Websocket
    type args struct {
        b bool
    }
    tests := []struct {
        name string
        args args
        want bool
    }{
        {"Reconnect 1", args{b: true}, true},
        {"Reconnect 2", args{b: false}, false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ws.Reconnect(tt.args.b)
            if ws.reconnect != tt.want {
                t.Errorf("Reconnect() want = %v, got %v", tt.want, ws.reconnect)
            }
        })
    }
}

func TestWs_autoReconnect(t *testing.T) {
   var ws = Websocket
   ws.SetUrl("localhost:2758", "/ws")

   var wsErr = Websocket
   wsErr.SetUrl("localhost:8080", "/ws")

   mux := http.NewServeMux()
   srv := &http.Server{Addr: ":2758", Handler: mux}
   // Starts a server to connect
   go func() {
       var upgrader = websocket.Upgrader{
           ReadBufferSize:  1024,
           WriteBufferSize: 1024,
       }
       mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
           upgrader.CheckOrigin = func(r *http.Request) bool { return true }
           connection, _ := upgrader.Upgrade(w, r, nil)
           defer connection.Close()
       })
       _ = srv.ListenAndServe()
   }()
   type fields struct {
       conn         *websocket.Conn
       caPool       *x509.CertPool
       secure       bool
       url          url.URL
       reconnect    bool
       websocketErr bool
       closeHandler func(int, string) error
   }
   tests := []struct {
       name   string
       fields fields
   }{
       {"autoReconnect 1", fields{websocketErr: ws.websocketErr, url: ws.url}},
       {"autoReconnect 2", fields{websocketErr: ws.websocketErr, url: wsErr.url}},
       //{"autoReconnect 2", fields{websocketErr: ws.websocketErr, url: wsErr.url}},
   }
   for _, tt := range tests {
       t.Run(tt.name, func(t *testing.T) {
           w := &Ws{
               conn:         tt.fields.conn,
               caPool:       tt.fields.caPool,
               secure:       tt.fields.secure,
               url:          tt.fields.url,
               reconnect:    tt.fields.reconnect,
               websocketErr: tt.fields.websocketErr,
               closeHandler: tt.fields.closeHandler,
           }
           go w.autoReconnect()
           time.Sleep(5 * time.Second)
           w.websocketErr = true
           time.Sleep(5 * time.Second)
           return
       })
   }
}

func TestWs_SetUrl(t *testing.T) {
    type args struct {
        host string
        path string
    }
    tests := []struct {
        name   string
        secure bool
        args   args
    }{
        {"SetUrl 1", false, args{host: "localhost:2758", path: "/ws"}},
        {"SetUrl 2", true, args{host: "localhost:2758", path: "/ws"}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := &Ws{
                secure: tt.secure,
            }
            w.SetUrl(tt.args.host, tt.args.path)
        })
    }
}

func TestWs_SetSecure(t *testing.T) {
    tests := []struct {
        name   string
        secure bool
    }{
        {"SetSecure 1", false},
        {"SetSecure 2", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := &Ws{
                secure: tt.secure,
            }
            w.SetSecure(tt.secure)
        })
    }
}

func TestWs_WriteMessage(t *testing.T) {
    var ws = Websocket
    ws.SetUrl("localhost:2758", "/ws")

    var wsErr = Websocket
    wsErr.SetUrl("localhost:2758", "/ws")

    mux := http.NewServeMux()
    srv := &http.Server{Addr: ":2758", Handler: mux}
    // Starts a server to connect
    go func() {
        var upgrader = websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
        }
        mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
            upgrader.CheckOrigin = func(r *http.Request) bool { return true }
            connection, _ := upgrader.Upgrade(w, r, nil)
            defer connection.Close()
        })
        _ = srv.ListenAndServe()
    }()
    type fields struct {
        url url.URL
    }
    type args struct {
        messageType int
        data        []byte
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        wantErr bool
    }{
        {"WriteMessage 1", fields{url: ws.url}, args{messageType: 1, data: []byte("hallo")}, false},
        {"WriteMessage 2", fields{url: wsErr.url}, args{messageType: 1, data: []byte("hallo")}, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := &Ws{
                url: tt.fields.url,
            }
            w.dial()
            if tt.wantErr {
                w.conn.Close()
            }
            if err := w.WriteMessage(tt.args.messageType, tt.args.data); (err != nil) != tt.wantErr {
                t.Errorf("WriteMessage() error = %v, wantErr %v", err, tt.wantErr)
            }
            w.conn.Close()
        })
    }
}

func TestWs_WriteJSON(t *testing.T) {
    var ws = Websocket
    ws.SetUrl("localhost:2758", "/ws")

    var wsErr = Websocket
    wsErr.SetUrl("localhost:2758", "/ws")

    mux := http.NewServeMux()
    srv := &http.Server{Addr: ":2758", Handler: mux}
    // Starts a server to connect to
    go func() {
        var upgrader = websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
        }
        mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
            upgrader.CheckOrigin = func(r *http.Request) bool { return true }
            connection, _ := upgrader.Upgrade(w, r, nil)
            defer connection.Close()
        })
        _ = srv.ListenAndServe()
    }()
    type test struct {
        TestMsg string `json:"test_msg"`
    }
    data, _ := json.Marshal(test{})
    type fields struct {
        //conn *websocket.Conn
        url url.URL
    }
    type args struct {
        v interface{}
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        wantErr bool
    }{
        {"WriteJSON 1", fields{url: ws.url}, args{data}, false},
        {"WriteJSON 2", fields{url: wsErr.url}, args{data}, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := &Ws{
                url: tt.fields.url,
            }
            w.dial()
            if tt.wantErr {
                w.conn.Close()
            }
            if err := w.WriteJSON(tt.args.v); (err != nil) != tt.wantErr {
                t.Errorf("WriteJSON() error = %v, wantErr %v", err, tt.wantErr)
            }
            w.conn.Close()
        })
    }
}

func TestWs_Read(t *testing.T) {
    var ws = Websocket
    ws.SetUrl("localhost:2758", "/ws")

    var wsErr = Websocket
    wsErr.SetUrl("localhost:2758", "/ws")

    mux := http.NewServeMux()
    srv := &http.Server{Addr: ":2758", Handler: mux}
    // Starts a server to connect to
    go func() {
        var upgrader = websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
        }
        mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
            upgrader.CheckOrigin = func(r *http.Request) bool { return true }
            connection, _ := upgrader.Upgrade(w, r, nil)
            connection.WriteMessage(1, []byte("test"))
            defer connection.Close()
        })
        _ = srv.ListenAndServe()
    }()
    type fields struct {
        url url.URL
    }
    tests := []struct {
        name     string
        fields   fields
        wantCode int
        wantByte []byte
        wantErr  bool
    }{
        //{"Read 1", fields{url: ws.url}, 1, []byte("test"), false},
        {"Read 2", fields{url: wsErr.url}, 0, []byte(""), true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := &Ws{
                url: tt.fields.url,
            }
            w.dial()
            if tt.wantErr {
                w.conn.Close()
            }
            //defer w.conn.Close()

            got, got1, err := w.Read()
            if (err != nil) != tt.wantErr {
                t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.wantCode {
                t.Errorf("Read() got = %v, want %v", got, tt.wantCode)
            }
            if !reflect.DeepEqual(got1, tt.wantByte) {
                t.Errorf("Read() got1 = %v, want %v", got1, tt.wantByte)
            }
            w.conn.Close()
        })
    }
}

func TestWs_ReadJSON(t *testing.T) {
    var ws = Websocket
    ws.SetUrl("localhost:2758", "/ws")

    var wsErr = Websocket
    wsErr.SetUrl("localhost:2758", "/ws")

    mux := http.NewServeMux()
    srv := &http.Server{Addr: ":2758", Handler: mux}

    type test struct {
        TestMsg string `json:"test_msg"`
    }

    var testStruct test
    // Starts a server to connect
    go func() {
        var connection *websocket.Conn
        var upgrader = websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
        }
        mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
            upgrader.CheckOrigin = func(r *http.Request) bool { return true }
            connection, _ = upgrader.Upgrade(w, r, nil)
            connection.ReadJSON(&testStruct)
            connection.WriteJSON(&testStruct)
            defer connection.Close()
        })
        _ = srv.ListenAndServe()
    }()
    type fields struct {
        url url.URL
    }
    type args struct {
        v interface{}
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        wantErr bool
    }{
        //{"ReadJson 1", fields{url: ws.url}, args{&testStruct}, false},
        {"ReadJson 2", fields{url: wsErr.url}, args{&testStruct}, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := &Ws{
                url: tt.fields.url,
            }
            w.dial()
            if tt.wantErr {
              w.conn.Close()
            }else if !tt.wantErr{
                w.WriteJSON(&testStruct)
            }
            if err := w.ReadJSON(tt.args.v); (err != nil) != tt.wantErr {
                t.Errorf("ReadJSON() error = %v, wantErr %v", err, tt.wantErr)
            }
            w.conn.Close()
        })
    }
}

func TestWs_CloseConnection(t *testing.T) {
    var ws = Websocket
    ws.SetUrl("localhost:2758", "/ws")

    var wsErr = Websocket
    wsErr.SetUrl("localhost:2758", "/ws")

    mux := http.NewServeMux()
    srv := &http.Server{Addr: ":2758", Handler: mux}
    // Starts a server to connect to
    go func() {
        var upgrader = websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
        }
        mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
            upgrader.CheckOrigin = func(r *http.Request) bool { return true }
            connection, _ := upgrader.Upgrade(w, r, nil)
            defer connection.Close()
        })
        _ = srv.ListenAndServe()
    }()

    type fields struct {
        url url.URL
    }
    type args struct {
        code int
        text string
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        wantErr bool
    }{
        {"CloseConnection 1", fields{url: ws.url}, args{1000, "Closed connection"}, false},
        {"CloseConnection 1", fields{url: wsErr.url}, args{1000, "Closed connection"}, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := &Ws{
                url: tt.fields.url,
            }
            w.dial()
            if tt.wantErr {
                w.conn.Close()
            }
            if err := w.CloseConnection(tt.args.code, tt.args.text); (err != nil) != tt.wantErr {
                t.Errorf("CloseConnection() error = %v, wantErr %v", err, tt.wantErr)
            }
            w.conn.Close()
        })
    }
}

func TestWs_AppendCert(t *testing.T) {
    serverCert, _ := ioutil.ReadFile("server.crt")

    type fields struct {
        caPool *x509.CertPool
    }
    type args struct {
        cert []byte
    }
    tests := []struct {
        name   string
        fields fields
        args   args
    }{
        {"AppendCert 1", fields{}, args{serverCert}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := &Ws{
                caPool: tt.fields.caPool,
            }
            w.SetSecure(true)
            w.AppendCert(serverCert)
        })
    }
}

func TestWs_CheckConnection(t *testing.T) {
   var ws = Websocket
   ws.SetUrl("localhost:2758", "/ws")

   var wsErr = Websocket
   wsErr.SetUrl("localhost:2758", "/ws")

   mux := http.NewServeMux()
   srv := &http.Server{Addr: ":2758", Handler: mux}
   // Starts a server to connect to
   go func() {
       var upgrader = websocket.Upgrader{
           ReadBufferSize:  1024,
           WriteBufferSize: 1024,
       }
       mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
           upgrader.CheckOrigin = func(r *http.Request) bool { return true }
           connection, _ := upgrader.Upgrade(w, r, nil)
           defer connection.Close()
       })
       _ = srv.ListenAndServe()
   }()
   type fields struct {
       url url.URL
   }
   tests := []struct {
       name   string
       fields fields
   }{
       {"checkConnection 1", fields{url: ws.url}},
   }
   for _, tt := range tests {
       t.Run(tt.name, func(t *testing.T) {
           w := &Ws{
               url: tt.fields.url,
           }
           w.dial()
           go w.CheckConnection()
           time.Sleep(5 * time.Second)
           w.conn.Close()
           time.Sleep(60*time.Second)
           return
       })
   }
}
