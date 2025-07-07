package conn

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/portilho13/blockchain/block"
	"github.com/portilho13/blockchain/models"
)

type Connection struct {
	BlockChain       *block.Blockchain
	ServerConnection map[string]*net.Conn
	ClientConnection map[string]*net.Conn
	NodesIps         []string
	ServerAddr       string
}

const DOMAIN_IP = "localhost:8000"

func (c *Connection) StartServer(ip string) {

	l, err := net.Listen("tcp", ip)
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		buffer := make([]byte, 1024)

		n, err := conn.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}

		addr := string(buffer[:n])

		fmt.Printf("Connection From %s\n", addr)
		if _, exists := c.ServerConnection[addr]; !exists && c.ClientConnection[addr] == nil {
			c.ServerConnection[addr] = &conn
			go c.ConnectToSingleClient(addr) // Connect to client server as client
			go c.HandleBlockBroadcast(&conn) // Handle received client broadcasted blocks
		}
	}
}

func (c *Connection) StartClient() {
	for _, ip := range c.NodesIps {
		c.ConnectToSingleClient(ip)
	}
}

func (c *Connection) ConnectToSingleClient(ip string) {
	if _, exists := c.ClientConnection[ip]; exists { // Already connected, skip
		return
	}

	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write([]byte(c.ServerAddr))
	if err != nil {
		log.Fatal(err)
	}

	c.ClientConnection[ip] = &conn
}

func (c *Connection) AddIpToList(ip string) {

	if ip == c.ServerAddr { // Dont add own server ip
		return
	}

	for _, existingIP := range c.NodesIps {
		if existingIP == ip {
			return
		}
	}

	c.NodesIps = append(c.NodesIps, ip)
}

func (c *Connection) ResolveHosts(ip string) []string {
	addr, err := net.ResolveUDPAddr("udp4", ip)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	_, err = conn.Write([]byte(c.ServerAddr))
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 1024)

	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Fatal(err)
	}

	var data []string

	err = json.Unmarshal(buffer[:n], &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func (c *Connection) Start(ip string) {
	c.ServerAddr = ip

	c.ServerConnection = make(map[string]*net.Conn)
	c.ClientConnection = make(map[string]*net.Conn)

	ips := c.ResolveHosts(DOMAIN_IP)

	for _, ip := range ips {
		c.AddIpToList(ip)
	}

	fmt.Println(c.NodesIps)

	go c.StartServer(ip)
	if len(c.NodesIps) != 0 {
		go c.StartClient()
	}

}

func (c *Connection) PrintConnectionMap() {
	fmt.Println("ServerConnection map:")
	for key, val := range c.ServerConnection {
		if val != nil {
			fmt.Printf("Key: %s, Value: %v\n", key, *val)
		} else {
			fmt.Printf("Key: %s, Value: nil\n", key)
		}
	}

	fmt.Println("ClientConnection map:")
	for key, val := range c.ClientConnection {
		if val != nil {
			fmt.Printf("Key: %s, Value: %v\n", key, *val)
		} else {
			fmt.Printf("Key: %s, Value: nil\n", key)
		}
	}
}

func (c *Connection) BroadcastBlockToAll(b models.Block) error {
	for sv, conn := range c.ClientConnection {
		log.Printf("Sending block info %s to %s", b.BlockHeader.Hash, sv)

		blockAsBytes, err := json.Marshal(b)
		if err != nil {
			return err
		}

		_, err = (*conn).Write(blockAsBytes) // Send block as bytes to all clients
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func (c *Connection) HandleBlockBroadcast(conn *net.Conn) error {
	buffer := make([]byte, 1024)

	n, err := (*conn).Read(buffer)
	if err != nil {
		return err
	}

	var b models.Block

	fmt.Println("BlockInfo: ", b)

	err = json.Unmarshal(buffer[:n], &b)
	if err != nil {
		return err
	}

	log.Printf("Received block info %s", b.BlockHeader.Hash)
	c.BlockChain.Addblock(b)
	c.BlockChain.PrintBlockchain()
	return nil
}
