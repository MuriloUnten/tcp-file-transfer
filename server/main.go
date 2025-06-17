package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/MuriloUnten/tcp-file-transfer/protocol"
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

		// the server should just ignore any message that is not a Request
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
			s.handleFetchRequest(conn, req)
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
	sse := protocol.NewSSE(message)
	out, err := protocol.EncodeMessage(sse)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range s.conns {
		_, err := c.Write(out)
		if err != nil {
			fmt.Println("write error:", err)
		}
	}
}

func (s *Server) handleFetchRequest(conn net.Conn, req *protocol.Request) {
	fileName := req.Body
	if len(fileName) == 0 {
		s.WriteResponse(conn, protocol.BadRequest, "invalid empty file name")
	}

	file, err := os.Open(s.fileDirectory + "/" + fileName)
	if err != nil {
		if os.IsNotExist(err) {
			s.WriteResponse(conn, protocol.NotFound, err.Error())
			return
		}
		s.WriteResponse(conn, protocol.InternalError, err.Error())
		return
	}
	defer file.Close()
	
	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}
	fileHash := string(h.Sum(nil))
	file.Seek(0, 0) // return file pointer to the beginning

	s.WriteResponse(conn, protocol.Ok, fileHash)

	buf := make([]byte, 4096)
	for {
		n, err := file.Read(buf)
		fmt.Printf("\n\n\t%d\n\n", n)
		if err != nil {
			if err != io.EOF {
				s.WriteResponse(conn, protocol.InternalError, err.Error())
				return
			}

			// handle EOF
			stream := protocol.NewStream(0, "EOF")
			err = s.sendMessage(conn, stream)
			if err != nil {
				s.WriteResponse(conn, protocol.InternalError, err.Error())
				return
			}
			return
		}

		stream := protocol.NewStream(n, string(buf))
		err = s.sendMessage(conn, stream)
		if err != nil {
			s.WriteResponse(conn, protocol.InternalError, err.Error())
			return
		}
	}
}

func (s *Server) sendMessage(conn net.Conn, msg protocol.Message) error {
	out, err := protocol.EncodeMessage(msg)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = conn.Write(out)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func main() {
	s := NewServer(":3000", "files")
	log.Fatal(s.Start())
}
