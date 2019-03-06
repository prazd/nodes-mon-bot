package handlers

import (
	"github.com/anvie/port-scanner"
	"time"
)

type Result bool

func NodesAlive(address string, port int, c chan bool) {
	ps := portscanner.NewPortScanner(address, 2*time.Second, 5)
	c <- ps.IsOpen(port)
}
