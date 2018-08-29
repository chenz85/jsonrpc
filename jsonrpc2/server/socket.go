package server

import (
	"io"
	"log"
	"net"
)

type socketServer struct {
	BaseServer

	network string
	addr    string
}

func (svr *socketServer) Serve() {
	go svr._Serve()
}

func (svr *socketServer) _Serve() {
	if l, le := net.Listen(svr.network, svr.addr); le != nil {
		log.Println("[RPC] server exit with err:", le)
	} else {
		log.Println("[RPC] server listen on:", l.Addr())
		for {
			c_in, ae := l.Accept()
			if ae != nil {
				log.Println("[RPC] accept incoming conn with err:", ae)
				continue
			}
			log.Println("[RPC] receive incoming conn:", c_in.RemoteAddr())
			go svr.handleConn(c_in)
		}
	}

}

func (svr *socketServer) handleConn(c net.Conn) {
	defer c.Close()
	defer log.Println("[RPC] conn closed:", c.RemoteAddr())

	var buf []byte = make([]byte, 4096)
	for {
		if rn, re := c.Read(buf); re == nil && rn > 0 {
			req := buf[:rn]
			resp := svr.HandleRequest(req)

			var wn int
			for _wn, we := c.Write(resp); we == nil && wn+_wn < rn; {
				resp = resp[_wn:]
				wn += _wn
			}
		} else if re == io.EOF {
			break
		}
	}
}

// create a socket server
func NewSocketServer(network string, addr string) JSONRPCServer {
	var svr = &socketServer{
		network: network,
		addr:    addr,
	}
	return svr
}
