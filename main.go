package main

import (
	"heartbeat/client"
	"time"
)

func main() {
	ports := []int{8080, 8081, 8082}

	time.Sleep(2 * time.Second)
	client.StartClient(ports)
	select {}
}
