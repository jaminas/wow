package handler

import (
	"context"
)

type Guard struct {
}

// Protect
func (p *Guard) Protect(_ context.Context, inputHeader int, inputPayload string, clientInfo string) (bool, int, string, error) {
	return false, 0, "", nil
}
