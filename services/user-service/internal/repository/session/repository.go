package session

import (
	"context"
	"time"
	"user-service/internal/domain/model"

	"github.com/redis/go-redis/v9"
)

const sessionKeyPrefix = "session:"

type repository struct {
	client redis.Cmdable
}

type Repository interface {
	Set(ctx context.Context, accessToken string, session model.Session, ttl time.Duration) error
}

func NewRepository(client redis.Cmdable) Repository {
	return &repository{client: client}
}
