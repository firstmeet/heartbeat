package client

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

type HeartBeat interface {
	HeartBeat(fns ...Handle)
}

type Server struct {
	Conn      net.Conn
	Connected bool
	Fail      uint8
	Mutex     sync.Mutex
	Address   string
}
type Handle func(server *Server)

func StartClient(ports []int) {
	var servers []*Server
	for _, port := range ports {
		server := &Server{
			Address: fmt.Sprintf(":%d", port),
			Mutex:   sync.Mutex{},
			Fail:    0,
			Conn:    nil,
		}
		servers = append(servers, server)
	}
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, server := range servers {
			var heartBeat HeartBeat
			fn := func(server2 *Server) {
				fmt.Printf("Server %s is down\n", server2.Address)
			}
			heartBeat = server
			heartBeat.HeartBeat(fn)
		}
	}

}
func (s *Server) HeartBeat(handles ...Handle) {
	if s.Fail >= 5 {
		if len(handles) > 0 {
			for _, handle := range handles {
				handle(s)
			}
		}
	}
	if !s.Connected {
		s.connect()
		if s.Connected {
			go s.sendMsg()
			go s.receive()
		}
	}
}
func (s *Server) connect() {
	fmt.Println("Connecting to", s.Address)
	conn, err := net.Dial("tcp", s.Address)
	if err != nil {
		fmt.Println(err)
		s.Fail++
		return
	}
	fmt.Println("Connected to", s.Address)
	s.Conn = conn
	s.Connected = true
}
func (s *Server) disconnect() {
	s.Conn.Close()
	s.Connected = false
}
func (s *Server) sendMsg() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		err := s.ping()
		if err != nil {
			s.disconnect()
		}
	}
}

//ping
func (s *Server) ping() error {
	_, err := s.Conn.Write([]byte("ping\n"))
	if err != nil {
		s.Fail++
		return err
	}
	return nil
}

//receive
func (s *Server) receive() {
	reader := bufio.NewReader(s.Conn)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		if string(line) == "pong" {
			fmt.Printf("%s received pong\n", s.Address)
			s.Fail = 0
			s.Connected = true
		}
	}
}
