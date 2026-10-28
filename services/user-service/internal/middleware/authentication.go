package middleware

import (
	"errors"
	"user-service/internal/common/response"
	"user-service/internal/constants"
	errConstant "user-service/internal/constants/error"
	"user-service/internal/logger"
	sessionrepository "user-service/internal/repository/session"
	jwtservice "user-service/internal/services/jwt"
	"uuid"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Identity struct {
	UUID          uuid.UUID
	Email         string
	RoleID        []int64
	CSRFTokenHash string
}

type Authentication struct {
	jwtService        jwtservice.Service
	sessionRepository sessionrepository.Repository
}

func NewAuthentication(jwtService jwtservice.Service, sessionRepository sessionrepository.Repository) *Authentication {
	return &Authentication{
		jwtService:        jwtService,
		sessionRepository: sessionRepository,
	}
}

func (a *Authentication) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie(constants.AccessTokenName)
		if err != nil {
			response.Unauthorized(c, errConstant.ErrInvalidToken)
			c.Abort()
			return
		}

		claims, err := a.jwtService.ValidateAccessToken(accessToken)
		if err != nil {
			response.Unauthorized(c, errConstant.ErrInvalidToken)
			c.Abort()
			return
		}

		session, err := a.sessionRepository.Get(c.Request.Context(), claims.SessionID)
		if err != nil {
			if errors.Is(err, errConstant.ErrNotFound) {
				response.Unauthorized(c, errConstant.ErrInvalidToken)
				c.Abort()
				return
			}

			logger.WithContext(c.Request.Context()).WithError(err).Error("failed to get session")
			response.InternalServerError(c)
			c.Abort()
			return
		}

		if session.UUID != claims.UserUUID {
			logger.WithContext(c.Request.Context()).WithFields(logrus.Fields{
				"jwt_user_uuid":     claims.UserUUID,
				"session_user_uuid": session.UUID,
			}).Error("authentication identity mismatch")
			response.Unauthorized(c, errConstant.ErrInvalidToken)
			c.Abort()
			return
		}

		c.Set(constants.AuthenticatedIdentityKey, Identity{
			UUID:          session.UUID,
			Email:         session.Email,
			RoleID:        session.RoleID,
			CSRFTokenHash: session.CSRFTokenHash,
		})

		c.Next()
	}
}

func RequireIdentity(c *gin.Context) (Identity, bool) {
	value, exists := c.Get(constants.AuthenticatedIdentityKey)
	if !exists {
		logger.WithContext(c.Request.Context()).Error("authenticated identity missing from context")
		response.InternalServerError(c)
		c.Abort()
		return Identity{}, false
	}

	identity, ok := value.(Identity)
	if !ok {
		logger.WithContext(c.Request.Context()).Error("invalid authenticated identity type")
		response.InternalServerError(c)
		c.Abort()
		return Identity{}, false
	}

	return identity, true
}
