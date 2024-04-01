package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"wow/internal/handler"
	"wow/internal/server"
	"wow/pkg/cache"
	"wow/pkg/guard"
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
	log.Println("start server")

	//
	address := fmt.Sprintf("%s:%d", host, port)
	c := cache.InitInMemoryCache()
	h := handler.NewHandler()
	g := guard.NewGuard(c)
	s := server.NewServer(address, g, h)

	//
	err := s.Run(context.Background())
	if err != nil {
		log.Println("server error:", err)
	}
}
