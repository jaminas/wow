package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"wow/pkg/guard"
	"wow/pkg/pow"
	"wow/pkg/protocol"
)

type Client struct {
	address string
}

func NewClient(address string) *Client {
	return &Client{address: address}
}

// Run
func (c *Client) Run(ctx context.Context) error {

	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Println("connected")

	//
	message, err := c.handle(ctx, conn, conn)
	if err != nil {
		return err
	}
	log.Println("result:", message)
	return nil
}

// handle
// 1. request challenge from server
// 2. compute hashcash to check Proof of Work
// 3. send hashcash solution back to server
// 4. get result from server
func (c *Client) handle(_ context.Context, readerConn io.Reader, writerConn io.Writer) (string, error) {
	reader := bufio.NewReader(readerConn)

	// 1
	err := c.sendMessage(protocol.Message{
		Header: guard.RequestChallenge,
	}, writerConn)
	if err != nil {
		return "", fmt.Errorf("send request error: %w", err)
	}

	//
	msgStr, err := c.readMessage(reader)
	if err != nil {
		return "", fmt.Errorf("read message error: %w", err)
	}

	//
	msg, err := protocol.ParseMessage(msgStr)
	if err != nil {
		return "", fmt.Errorf("parse message error: %w", err)
	}

	//
	var hashcash pow.HashcashData
	err = json.Unmarshal([]byte(msg.Payload), &hashcash)
	if err != nil {
		return "", fmt.Errorf("parse hashcash errot: %w", err)
	}
	log.Println("hashcash:", hashcash)

	// 2
	hashcash, err = hashcash.ComputeHashcash(pow.MaxIterations)
	if err != nil {
		return "", fmt.Errorf("err compute hashcash: %w", err)
	}
	log.Println("hashcash computed:", hashcash)

	//
	byteData, err := json.Marshal(hashcash)
	if err != nil {
		return "", fmt.Errorf("marshal hashcash error: %w", err)
	}

	// 3
	err = c.sendMessage(protocol.Message{
		Header:  guard.RequestResource,
		Payload: string(byteData),
	}, writerConn)
	if err != nil {
		return "", fmt.Errorf("send request error: %w", err)
	}
	log.Println("challenge sent to server")

	// 4
	msgStr, err = c.readMessage(reader)
	if err != nil {
		return "", fmt.Errorf("read message error: %w", err)
	}
	msg, err = protocol.ParseMessage(msgStr)
	if err != nil {
		return "", fmt.Errorf("parse message error: %w", err)
	}
	return msg.Payload, nil
}

// readMessage
func (c *Client) readMessage(reader *bufio.Reader) (string, error) {
	return reader.ReadString('\n')
}

// sendMessage
func (c *Client) sendMessage(msg protocol.Message, conn io.Writer) error {
	msgStr := fmt.Sprintf("%s\n", msg.ToString())
	_, err := conn.Write([]byte(msgStr))
	return err
}
