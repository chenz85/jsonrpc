package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func StartHttpServer(host string, port uint16, path string) {
	var addr string = fmt.Sprintf("%s:%d", host, port)
	if path != "" && path[0] != '/' {
		path = "/" + path
	}
	http.HandleFunc(path, rpc)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println("rpc server exit with err:", err)
	} else {
		log.Println("prc server exit")
	}
}

func rpc(w http.ResponseWriter, r *http.Request) {
	if data, err := ioutil.ReadAll(r.Body); err != nil {
		// TODO: parse error
		log.Println("parse error:", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else if resp, err := ProcessRequest(string(data)); err != nil {
		w.Write([]byte(err.Json()))
	} else {
		w.Write([]byte(resp.Json()))
	}
}
