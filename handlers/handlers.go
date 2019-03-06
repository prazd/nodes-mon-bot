package handlers

import (
	"github.com/anvie/port-scanner"
	"time"
)

func NodesAlive(address string, port int, c chan bool) {
	ps := portscanner.NewPortScanner(address, 2*time.Second, 5)
	c <- ps.IsOpen(port)
}
