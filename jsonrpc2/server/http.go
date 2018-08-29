package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/czsilence/jsonrpc/jsonrpc2/object"
)

type httpServer struct {
	BaseServer

	server *http.Server
}

func (svr *httpServer) Serve() {
	go svr._Serve()
}

func (svr *httpServer) _Serve() {
	if err := svr.server.ListenAndServe(); err != nil {
		log.Println("rpc server exit with err:", err)
	} else {
		log.Println("prc server exit")
	}
}

func (svr *httpServer) rpc(w http.ResponseWriter, r *http.Request) {
	if data, err := ioutil.ReadAll(r.Body); err != nil {
		log.Println("parse error:", err)
		w.Write(object.ErrParse.JsonObject().ToJson())
	} else {
		w.Write(svr.HandleRequest(data))
	}
}

func NewHttpServer(host string, port uint16, path string) JSONRPCServer {

	var addr string = fmt.Sprintf("%s:%d", host, port)
	if path != "" && path[0] != '/' {
		path = "/" + path
	}

	mux := http.NewServeMux()

	var svr = &httpServer{
		server: &http.Server{Addr: addr, Handler: mux},
	}
	mux.HandleFunc(path, svr.rpc)

	return svr
}
