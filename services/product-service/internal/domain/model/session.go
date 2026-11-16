package model

import (
	"time"
	"uuid"
)

type Session struct {
	UUID          uuid.UUID `json:"uuid"`
	CSRFTokenHash string    `json:"csrf_token_hash"`
	Email         string    `json:"email"`
	RoleName      []string  `json:"role_name"`
	CreatedAt     time.Time `json:"created_at"`
}
