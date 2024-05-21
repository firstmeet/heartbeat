package client

import "testing"

func TestStartClient(t *testing.T) {
	ports := []int{8080, 8081, 8082}
	StartClient(ports)
	select {}
}
