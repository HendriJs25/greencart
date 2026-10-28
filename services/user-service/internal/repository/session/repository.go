package session

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json/v2"
	"errors"
	"fmt"
	"time"
	errConstant "user-service/internal/constants/error"
	"user-service/internal/domain/model"
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
	Get(ctx context.Context, ID string) (*model.Session, error)
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

func (r *repository) Get(ctx context.Context, ID string) (*model.Session, error) {
	payload, err := r.client.Get(ctx, sessionKey(ID)).Bytes()
	if err != nil {
		if errors.Is(err, redislib.Nil) {
			return nil, errConstant.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session model.Session
	err = json.Unmarshal(payload, &session)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}
	return &session, nil
}

func sessionKey(ID string) string {
	digest := sha256.Sum256([]byte(ID))

	return sessionKeyPrefix + hex.EncodeToString(digest[:])
}

func userSessionKey(uuid uuid.UUID) string {
	return fmt.Sprintf("%s%s", userSessionKeyPrefix, uuid.String())
}
