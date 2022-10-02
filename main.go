package main

import (
	"flag"
	"log"
	"time"

	"github.com/ductnn/tinylb/internal/loadbalancer"
	"github.com/ductnn/tinylb/internal/utils"
)

var (
	flagURL utils.FlagURL
	port    int
)

func main() {
	flag.Var(&flagURL, "backends", "Load balanced backends, use commas to separate")
	flag.IntVar(&port, "port", 4000, "Port to serve")
	flag.Parse()

	if len(flagURL.URLs) == 0 {
		log.Fatal("Please provide one or more backends to load balance")
	}

	loadBalancer := loadbalancer.New()
	loadBalancer.Register(flagURL.URLs...)

	go loadBalancer.HeathCheck(1 * time.Minute)
	loadBalancer.Listen(port)
}
