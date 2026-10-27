package session

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json/v2"
	"fmt"
	"time"
	"user-service/internal/domain/model"
	"user-service/internal/logger"
	"uuid"

	redislib "github.com/redis/go-redis/v9"
)

const (
	sessionKeyPrefix     = "session:"
	userSessionKeyPrefix = "user_session:"
)

type repository struct {
	client redislib.Cmdable
}

type Repository interface {
	Set(ctx context.Context, ID string, session model.Session, ttl time.Duration) error
}

func NewRepository(client redislib.Cmdable) Repository {
	return &repository{client: client}
}

func (r *repository) Set(ctx context.Context, ID string, session model.Session, ttl time.Duration) error {
	payload, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	sessionKey := sessionKey(ID)
	indexKey := userSessionKey(session.UUID)
	logger.Log.WithField("session_key", sessionKey).Debug()

	_, err = r.client.TxPipelined(ctx, func(pipe redislib.Pipeliner) error {
		pipe.Set(ctx, sessionKey, payload, ttl)
		pipe.SAdd(ctx, indexKey, sessionKey)
		pipe.Expire(ctx, indexKey, ttl)
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to store session: %w", err)
	}
	return nil
}

func sessionKey(ID string) string {
	digest := sha256.Sum256([]byte(ID))

	return sessionKeyPrefix + hex.EncodeToString(digest[:])
}

func userSessionKey(uuid uuid.UUID) string {
	return fmt.Sprintf("%s%s", userSessionKeyPrefix, uuid.String())
}
