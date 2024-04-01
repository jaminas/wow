package server

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"wow/pkg/protocol"
)

type Server struct {
	address string
	guard   Guard
	handler Handler
	// TODO worker_pool - он тут обязан быть дабы мы не жрали память бесконечно
	// TODO buffer limits - ограничения буффера на чтение, дабы нам не слали километровые запросы
}

func NewServer(
	address string,
	guard Guard,
	handler Handler,
) *Server {
	return &Server{
		address: address,
		guard:   guard,
		handler: handler,
	}
}

// Run
func (s *Server) Run(ctx context.Context) error {
	//
	listener, err := net.Listen(CONN_TYPE, s.address)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Println("listening", listener.Addr())

	//
	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("accept connection error: %w", err)
		}
		go s.handle(ctx, conn)
	}
}

// handle
func (s *Server) handle(ctx context.Context, conn net.Conn) {
	log.Println("client:", conn.RemoteAddr())
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		//
		req, err := reader.ReadString('\n')
		if err != nil {
			log.Println("read connection error:", err)
			return
		}

		//
		inputMessage, err := protocol.ParseMessage(req)
		if err != nil {
			log.Println("parse message error:", err)
			return
		}

		//
		success, header, payload, err := s.guard.Protect(
			ctx,
			inputMessage.Header,
			inputMessage.Payload,
			conn.RemoteAddr().String(),
		)
		if err != nil {
			log.Println("failed protection:", err)
			return
		}

		//
		if success {
			payload, err = s.handler.Handle(ctx)
			if err != nil {
				log.Println("handle error:", err)
				return
			}
		}

		//
		outputMessage := protocol.Message{
			Header:  header,
			Payload: payload,
		}
		err = s.sendMessage(outputMessage, conn)
		if err != nil {
			log.Println("send message error:", err)
		}
	}
}

// sendMessage
func (s *Server) sendMessage(msg protocol.Message, conn net.Conn) error {
	msgStr := fmt.Sprintf("%s\n", msg.ToString())
	_, err := conn.Write([]byte(msgStr))
	return err
}
