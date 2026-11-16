package session

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json/v2"
	"errors"
	"fmt"
	errConstant "product-service/internal/constants/error"
	"product-service/internal/domain/model"

	redislib "github.com/redis/go-redis/v9"
)

const (
	sessionKeyPrefix = "session:"
)

type repository struct {
	client redislib.Cmdable
}

type Repository interface {
	Get(ctx context.Context, ID string) (*model.Session, error)
}

func NewRepository(client redislib.Cmdable) Repository {
	return &repository{client: client}
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
