package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"bufio"
)

type Server struct {
	port   string
	ln     net.Listener
	conns  []net.Conn
	quitch chan struct{}
	stdinScanner  *bufio.Scanner
}

func NewServer(port string) *Server {
	return &Server{
		port:   port,
		quitch: make(chan struct{}),
		conns:  make([]net.Conn, 0),
		stdinScanner:  bufio.NewScanner(os.Stdin),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	defer ln.Close()
	s.ln = ln

	fmt.Println("Listening on", ln.Addr().String())

	go s.AcceptConnections()
	go s.HandleInput()

	<-s.quitch

	return nil
}

func (s *Server) AcceptConnections() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		s.conns = append(s.conns, conn)

		fmt.Println("Connection received from", conn.RemoteAddr())

		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error:", err) 
			continue
		}

		msg := string(buf[:n])
		fmt.Println("Message from ", conn.RemoteAddr())
		fmt.Println(msg)
	}
}

func (s *Server) HandleInput() {
	for s.stdinScanner.Scan() {
		message := s.stdinScanner.Text()
		s.BroadcastMessage(message)
	}
}

func (s *Server) BroadcastMessage(message string) {
	for _, c := range s.conns {
		// TODO modify this to send the message using the protocol for CHAT
		_, err := c.Write([]byte(message))
		if err != nil {
			fmt.Println("write error:", err)
		}
	}
}

func main() {
	s := NewServer(":3000")
	log.Fatal(s.Start())
}
