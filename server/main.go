package main

import (
	"github.com/MuriloUnten/tcp-file-transfer/protocol"
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

		msg := buf[:n]
		req := protocol.Request{}
		err = req.Decode(msg)
		if err != nil {
			fmt.Println("error decoding request:", err)
			s.WriteResponse(conn, protocol.BadRequest, err.Error())
			continue
		}

		fmt.Printf("Request from %s | %s\n", conn.RemoteAddr(), protocol.TranslateMethod(req.Method))
		fmt.Println(req.Body)

		switch req.Method {
		case protocol.Chat:
			s.WriteResponse(conn, protocol.Ok, "")
		case protocol.Fetch:

		case protocol.Quit:

		default:
			log.Fatal("ITS COOKED!! THIS SHOULD NEVER HAPPEN")
		}
	}
}

func (s *Server) HandleInput() {
	for s.stdinScanner.Scan() {
		message := s.stdinScanner.Text()
		s.BroadcastMessage(message)
	}
}

func (s *Server) WriteResponse(conn net.Conn, status protocol.StatusCode, body string) error {
	response := protocol.Response{
		StatusCode: status,
		Body: body,
	}
	out, err := response.Encode()
	if err != nil {
		log.Fatal("error encoding response:", err)
	}

	_, err = conn.Write(out)
	return err
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
