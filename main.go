package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message),
	}
}

func (s *Server) Start() error {
	fmt.Println("start a new serverr")
	ln, err := net.Listen("tcp", s.listenAddr)

	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.AcceptLoop()

	// TODO: what does the following line mean?
	<-s.quitch
	close(s.msgch)

	return nil
}

func (s *Server) AcceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		go s.ReadLoop(conn)
	}
}

func (s *Server) ReadLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}

		fmt.Println("New connection to the server:", conn.RemoteAddr())

		s.msgch <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}

		conn.Write([]byte("thank you for your message!\n"))
	}
}

func main() {
	server := NewServer(":3000")
	go func() {
		for msg := range server.msgch {
			fmt.Printf("received message from connection (%s): %s\n", msg.from, string(msg.payload))
		}
	}()
	log.Fatal(server.Start())
}
