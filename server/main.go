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
	port          string
	ln            net.Listener
	conns         []net.Conn
	quitch        chan struct{}
	stdinScanner  *bufio.Scanner
	fileDirectory string
}

func NewServer(port, fileDirectory string) *Server {
	if _, err := os.ReadDir(fileDirectory); err != nil {
		log.Fatal(err)
	}

	return &Server{
		port:          port,
		quitch:        make(chan struct{}),
		conns:         make([]net.Conn, 0),
		stdinScanner:  bufio.NewScanner(os.Stdin),
		fileDirectory: fileDirectory,
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
	buf := make([]byte, 4096 + 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error:", err) 
			continue
		}

		msgBytes := buf[:n]
		msg, err := protocol.DecodeMessage(msgBytes)
		if err != nil {
			fmt.Println("error decoding request:", err)
			s.WriteResponse(conn, protocol.BadRequest, err.Error())
			continue
		}
		req, ok := msg.(*protocol.Request)
		if !ok {
			continue
		}

		fmt.Printf("Request from %s | %s\n", conn.RemoteAddr(), req.Method)
		fmt.Println(req.Body)

		switch req.Method {
		case protocol.Chat:
			s.WriteResponse(conn, protocol.Ok, "")
		case protocol.Fetch:
			fileName := req.Body
			if len(fileName) == 0 {
				s.WriteResponse(conn, protocol.BadRequest, "invalid empty file name")
			}
			// fileBytes, err := os.ReadFile(s.fileDirectory + "/" + fileName)
			// if err != nil {
			// 	if os.IsNotExist(err) {
			// 		s.WriteResponse(conn, protocol.NotFound, err.Error())
			// 	}
			// }
		case protocol.Quit:

		default:
			log.Fatal("ITS COOKED!! THIS SHOULD NEVER HAPPEN")
		}
	}
}

func (s *Server) WriteResponse(conn net.Conn, status protocol.StatusCode, body string) error {
	response := protocol.NewResponse(status, body)
	out, err := protocol.EncodeMessage(response)
	if err != nil {
		log.Fatal("error encoding response:", err)
	}

	_, err = conn.Write(out)
	return err
}

func (s *Server) HandleInput() {
	for s.stdinScanner.Scan() {
		message := s.stdinScanner.Text()
		s.BroadcastMessage(message)
	}
}

func (s *Server) BroadcastMessage(message string) {
	sse := protocol.NewSSE()
	sse.Body = message
	out, err := sse.Encode()
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range s.conns {
		// TODO modify this to send the message using the protocol for CHAT
		_, err := c.Write(out)
		if err != nil {
			fmt.Println("write error:", err)
		}
	}
}

func main() {
	s := NewServer(":3000", "files")
	log.Fatal(s.Start())
}
