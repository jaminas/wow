package handler

import (
	"context"
)

type Handler struct {
}

// Handle
func (h *Handler) Handle(ctx context.Context) (string, error) {
	return "", nil
}
