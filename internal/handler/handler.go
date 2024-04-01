package handler

import (
	"context"
	"math/rand"
	"time"
)

type Handler struct {
	rgen *rand.Rand
}

// NewHandler
func NewHandler() *Handler {
	return &Handler{
		rgen: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Handle
func (h *Handler) Handle(ctx context.Context) (string, error) {
	return Answers[rand.Intn(len(Answers))], nil
}
