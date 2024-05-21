package server

import (
	"testing"
)

func TestStartServer(t *testing.T) {
	go StartServer(8080)
	go StartServer(8081)
	StartServer(8082)
}
