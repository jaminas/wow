package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"wow/internal/client"
)

var (
	host string
	port int
)

func init() {
	flag.StringVar(&host, "host", "localhost", "put host")
	flag.IntVar(&port, "port", 8082, "put port number")
	flag.Parse()
}

func main() {
	address := fmt.Sprintf("%s:%d", host, port)
	log.Println("starting client:", address)
	c := client.NewClient(address)
	err := c.Run(context.Background())
	if err != nil {
		log.Println("client error:", err)
	}
	log.Println("client finished")
}
