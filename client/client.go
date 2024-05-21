package client

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

type HeartBeat interface {
	HeartBeat()
}

const (
	DefaultMaxFailed = 5
	PingInterval     = 3
	RangeInterval    = 2
)

type Server struct {
	Conn         net.Conn
	Connected    bool
	Fail         uint8
	Mutex        sync.Mutex
	Address      string
	MaxFailed    uint8
	FailCallBack Handle
}
type Handle func(server *Server)

func StartClient(ports []int) {
	var servers []*Server
	for _, port := range ports {
		server := &Server{
			Address:   fmt.Sprintf(":%d", port),
			Mutex:     sync.Mutex{},
			Fail:      0,
			Conn:      nil,
			MaxFailed: DefaultMaxFailed,
		}
		servers = append(servers, server)
	}
	ticker := time.NewTicker(RangeInterval * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, server := range servers {
			var heartBeat HeartBeat
			fn := func(server2 *Server) {
				fmt.Printf("Server %s is down\n", server2.Address)
			}
			server.FailCallBack = fn
			heartBeat.HeartBeat()
		}
	}

}
func (s *Server) HeartBeat() {
	if s.Fail >= s.MaxFailed {
		if s.FailCallBack != nil {
			s.FailCallBack(s)
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
	ticker := time.NewTicker(PingInterval * time.Second)
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
		s.incrementFail()
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
			s.resetFail()
			s.Connected = true
		}
	}
}

//increment fail
func (s *Server) incrementFail() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Fail++
}

//reset fail
func (s *Server) resetFail() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Fail = 0
}
