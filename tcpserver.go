package tcpserver

import (
	"crypto/tls"
	"log"
	"net"

	"github.com/brucewangzhihua/logger"
)

// Client holds info about connection
type Client struct {
	Connection net.Conn
	Server     *server
}

// TCP server
type server struct {
	address                  string // Address to open connection: localhost:9999
	config                   *tls.Config
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
}

// listen Read client data from channel
func (c *Client) listen() {
	c.Server.onNewClientCallback(c)
}

// Send text message to client
func (c *Client) Send(message []byte) error {
	logger.Debug("Send", message)
	_, err := c.Connection.Write(message)
	return err
}

// SendBytes Send bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.Connection.Write(b)
	return err
}

// Conn Get connection
func (c *Client) Conn() net.Conn {
	return c.Connection
}

// Close Close server
func (c *Client) Close() error {
	return c.Connection.Close()
}

// Called right after server starts listening new client
func (s *server) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// OnClientConnectionClosed Called right after connection closed
func (s *server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// Listen starts network server
func (s *server) Listen() {
	var listener net.Listener
	var err error
	if s.config == nil {
		listener, err = net.Listen("tcp", s.address)
	} else {
		listener, err = tls.Listen("tcp", s.address, s.config)
	}
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		Connection, _ := listener.Accept()
		client := &Client{
			Connection: Connection,
			Server:     s,
		}
		go client.listen()
	}
}

// New Creates new tcp server instance
func New(address string) *server {
	log.Println("Creating server with address", address)
	server := &server{
		address: address,
		config:  nil,
	}

	server.OnNewClient(func(c *Client) {})
	server.OnClientConnectionClosed(func(c *Client, err error) {})

	return server
}

// NewWithTLS Creates new ssl tcp server instance
func NewWithTLS(address string, certFile string, keyFile string) *server {
	log.Println("Creating server with address", address)
	cert, _ := tls.LoadX509KeyPair(certFile, keyFile)
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	server := &server{
		address: address,
		config:  &config,
	}

	server.OnNewClient(func(c *Client) {})
	server.OnClientConnectionClosed(func(c *Client, err error) {})

	return server
}
