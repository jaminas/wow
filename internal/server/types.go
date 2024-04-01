package server

import "context"

const CONN_TYPE = "tcp"

type Handler interface {
	// Handle - Success result
	Handle(ctx context.Context) (string, error)
}

type Guard interface {
	// Protect
	Protect(ctx context.Context, inputHeader int, inputPayload string, clientInfo string) (success bool, outputHeader int, outputPayload string, err error)
}
