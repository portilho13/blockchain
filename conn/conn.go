package conn

import (
	"fmt"
	"log"
	"net"
)

type Connection struct {
	ServerConnection map[string]*net.Conn
	ClientConnection map[string]*net.Conn
	NodesIps         []string
}

func (c *Connection) StartServer(ip string, port int) {
	addrString := net.JoinHostPort(ip, fmt.Sprintf("%d", port))

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

		remoteAddr := conn.RemoteAddr()

		fmt.Printf("Connection From %s\n", conn.RemoteAddr())
		c.ServerConnection[conn.RemoteAddr().String()] = &conn

		go c.ConnectToSingleClient(remoteAddr.String()) // Connect as a client to the node entering the network
	}
}

func (c *Connection) StartClient() {
	for _, ip := range c.NodesIps {
		conn, err := net.Dial("tcp", ip)
		if err != nil {
			log.Fatal(err)
		}

		c.ClientConnection[conn.LocalAddr().String()] = &conn
	}

}

func (c *Connection) ConnectToSingleClient(ip string) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Fatal(err)
	}

	c.ClientConnection[conn.LocalAddr().String()] = &conn
}

func (c *Connection) AddIpToList(ip string) {

	for _, existingIP := range c.NodesIps {
		if existingIP == ip {
			return
		}
	}

	c.NodesIps = append(c.NodesIps, ip)
}

func (c *Connection) ResolveHosts(ip string) {
	addr, err := net.ResolveUDPAddr("udp4", ip)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	_, err = conn.Write([]byte("0"))
	if err != nil {
		log.Fatal(err)
	}
}
