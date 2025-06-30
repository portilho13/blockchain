package conn

import (
	"fmt"
	"log"
	"net"
)

type Connection struct {
}

func (c *Connection) StartServer(ip string, port int) {
	addrString := fmt.Sprintf("%s:%d", ip, port)

	conn, err := net.Listen("tcp", addrString)
	if err != nil {
		log.Fatal(conn)
	}

	defer conn.Close()

	for {
		c, err := conn.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Connection From %s\n", c.RemoteAddr())

	}

}
