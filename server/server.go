package server

import (
	"bufio"
	"fmt"
	"net"
)

func StartServer(port int) {
	fmt.Println("Starting server on port", port)
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		accept, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConn(accept)
	}
}
func handleConn(conn net.Conn) {
	defer conn.Close()
	ReceiveMsg(conn)
}
func ReceiveMsg(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		if string(line) == "ping" {
			fmt.Printf("%s received ping\n", conn.LocalAddr())
			err = SendMsg(conn)
			if err != nil {
				return
			}
		}
	}
}
func SendMsg(conn net.Conn) error {
	_, err := conn.Write([]byte("pong\n"))
	if err != nil {
		return err
	}
	return nil
}
