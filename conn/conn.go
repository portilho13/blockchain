package conn

import (
	"fmt"
	"log"
	"net"
)

type Connection struct {
	ServerConnection map[string]*net.Conn
	ClientConnection map[string]*net.Conn
}

func (c *Connection) StartServer(ip string, port int) {
	addrString := fmt.Sprintf("%s:%d", ip, port)

	l, err := net.Listen("tcp", addrString)
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Connection From %s\n", conn.RemoteAddr())
		c.ServerConnection[conn.RemoteAddr().String()] = &conn
	}
}

func (c *Connection) StartClient(ip string, port int) {
	addrString := fmt.Sprintf("%s:%d", ip, port)

	conn, err := net.Dial("tcp", addrString)
	if err != nil {
		log.Fatal(err)
	}

	c.ClientConnection[conn.LocalAddr().String()] = &conn
}
