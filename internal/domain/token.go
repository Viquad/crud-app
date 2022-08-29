package domain

import (
	"context"
	"time"
)

type RefreshSession struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
}

type TokenRepository interface {
	Create(ctx context.Context, token RefreshSession) error
	Get(ctx context.Context, token string) (*RefreshSession, error)
}
