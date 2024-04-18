package gomessage

import (
	"net"
	"strings"

	"github.com/jaxoncarelos/GoMessage/helper"
)

type GoMessage struct {
	Server      net.Listener
	Commands    map[string]CommandFunc
	Connections []net.Conn
	OnConnect   func(net.Conn)
	OnMessage   func(string, net.Conn)
}

func (g *GoMessage) Start(port string) {
	server, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	g.Server = server
	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}
		g.Connections = append(g.Connections, conn)
		go g.OnConnect(conn)
		go g.ReadMessage(conn)
	}
}

func (g *GoMessage) ReadMessage(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		g.Disconnect(conn)
		return
	}
	parts := strings.Split(string(buf[:n]), ":")
	if len(parts) != 3 {
		g.Disconnect(conn)
		return
	}
	if f, ok := g.Commands[parts[0]]; ok {
		err := f(parts[1] + ":" + parts[2])
		if err != nil {
			conn.Write([]byte(err.Error()))
		}
	}
	g.OnMessage(string(buf[:n]), conn)
}

func (g *GoMessage) Disconnect(conn net.Conn) {
	conn.Close()
	g.Connections = helper.Filter(g.Connections, func(c net.Conn) bool { return c != conn })
}

func (g *GoMessage) Stop() {
	g.Server.Close()
}

func NewGoMessage() *GoMessage {
	return &GoMessage{Commands: make(map[string]CommandFunc)}
}

func (g *GoMessage) AddCommand(alias string, f CommandFunc) {
	g.Commands[alias] = f
}

// command func should take in a string and return an errora and return an error
type CommandFunc func(string) error
