package auth

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
	"user-service/internal/common/helper"
	"user-service/internal/config"
	"user-service/internal/constants"
	errConstant "user-service/internal/constants/error"
	"user-service/internal/domain/event"
	"user-service/internal/domain/model"
	"user-service/internal/logger"
	"user-service/internal/messaging/kafka"
	"user-service/internal/repository"
	sessionrepository "user-service/internal/repository/session"
	userrepository "user-service/internal/repository/user"
	verificationtokenrepository "user-service/internal/repository/verificationtoken"
	jwtservice "user-service/internal/services/jwt"
	"uuid"
)

type service struct {
	jwtService                  jwtservice.Service
	userRepository              userrepository.Repository
	sessionRepository           sessionrepository.Repository
	verificationTokenRepository verificationtokenrepository.Repository
	transactionManager          repository.TransactionManager
	kafkaProducer               kafka.Producer
	configURL                   *config.URL
}

type Service interface {
	SignIn(ctx context.Context, input SignInInput) (*SignInResult, error)
	SignUp(ctx context.Context, input SignUpInput) error
	SignOut(ctx context.Context, sessionID string) error
	ForgotPassword(ctx context.Context, email string) error
	VerifyEmail(ctx context.Context, token string) error
	ResetPassword(ctx context.Context, input ResetPasswordInput) error
}

func NewService(jwtService jwtservice.Service, userRepository userrepository.Repository, sessionRepository sessionrepository.Repository, verificationTokenRepository verificationtokenrepository.Repository,
	transactionManager repository.TransactionManager, kafkaProducer kafka.Producer, configURL *config.URL) Service {
	return &service{
		jwtService:                  jwtService,
		userRepository:              userRepository,
		sessionRepository:           sessionRepository,
		verificationTokenRepository: verificationTokenRepository,
		transactionManager:          transactionManager,
		kafkaProducer:               kafkaProducer,
		configURL:                   configURL,
	}
}

func (s *service) SignIn(ctx context.Context, input SignInInput) (*SignInResult, error) {
	normalizedEmail := helper.NormalizeEmail(input.Email)
	user, err := s.userRepository.FindByEmail(ctx, normalizedEmail)
	if err != nil {
		if errors.Is(err, errConstant.ErrNotFound) {
			return nil, errConstant.ErrUserNotFound
		}
		return nil, err
	}

	isVerify := helper.VerifyPassword(input.Password, user.PasswordHash)
	if !isVerify {
		return nil, errConstant.ErrPasswordIncorrect
	}

	csrfToken, err := helper.GenerateCSRFToken()
	if err != nil {
		return nil, err
	}

	hashedCSRFToken := helper.HashToken(csrfToken)

	jti := uuid.NewV4()

	accessToken, err := s.jwtService.GenerateAccessToken(user.UUID, jti)
	if err != nil {
		return nil, err
	}

	var userRoles []string
	for _, role := range user.Roles {
		userRoles = append(userRoles, role.Name)
	}

	sessionCreatedAt := time.Now().UTC()

	sessionTTL := accessToken.ExpiresAt.Sub(sessionCreatedAt)

	session := model.Session{
		UUID:          user.UUID,
		CSRFTokenHash: hashedCSRFToken,
		Email:         user.Email,
		RoleName:      userRoles,
		CreatedAt:     sessionCreatedAt,
	}

	err = s.sessionRepository.Set(ctx, jti.String(), session, sessionTTL)
	if err != nil {
		return nil, err
	}

	return ToSignInResult(user, accessToken, csrfToken), nil
}

func (s *service) SignUp(ctx context.Context, input SignUpInput) error {
	normalizedEmail := helper.NormalizeEmail(input.Email)
	exist, err := s.userRepository.IsEmailExist(ctx, normalizedEmail)
	if err != nil {
		return err
	}

	if exist {
		return errConstant.ErrEmailAlreadyExists
	}

	hashedPassword, err := helper.HashPassword(input.Password)
	if err != nil {
		return err
	}

	verificationToken, err := helper.GenerateVerificationToken()
	if err != nil {
		return err
	}

	verifyURL, err := url.Parse(s.configURL.VerifyBaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse verify URL: %w", err)
	}
	query := verifyURL.Query()
	query.Set("token", verificationToken)
	verifyURL.RawQuery = query.Encode()

	var payload []byte
	var userID int64

	hashedVerificationToken := helper.HashToken(verificationToken)
	expiresAt := time.Now().UTC().Add(time.Hour)

	err = s.transactionManager.WithinTransaction(ctx, func(repositories *repository.Registry) error {
		customerRole, err := repositories.RoleRepository.FindByName(ctx, constants.Customer)
		if err != nil {
			if errors.Is(err, errConstant.ErrNotFound) {
				return fmt.Errorf("misconfigured: customer role not found")
			}
			return err
		}

		user := &model.User{
			UUID:         uuid.NewV7(),
			Name:         input.Name,
			Email:        normalizedEmail,
			PasswordHash: hashedPassword,
			IsVerified:   false,
			Roles:        []model.Role{*customerRole},
		}

		err = repositories.UserRepository.Create(ctx, user)
		if err != nil {
			return err
		}

		token := &model.VerificationToken{
			UserID:    user.ID,
			TokenHash: hashedVerificationToken,
			TokenType: constants.TokenTypeEmailVerification,
			ExpiresAt: expiresAt,
		}

		err = repositories.VerificationTokenRepository.Create(ctx, token)
		if err != nil {
			return err
		}

		notification := event.EmailVerificationRequested{
			EventID:          uuid.NewV4().String(),
			ReceiverID:       user.ID,
			ReceiverEmail:    normalizedEmail,
			Subject:          "【GreenCart】メールアドレス認証のお願い",
			Message:          "以下のリンクをクリックして、アカウント認証を完了してください: " + verifyURL.String(),
			NotificationType: "EMAIL",
		}

		payload, err = json.Marshal(notification)
		if err != nil {
			return fmt.Errorf("failed to marshal notification: %w", err)
		}
		userID = user.ID

		return nil
	})
	if err != nil {
		return err
	}

	key := strconv.FormatInt(userID, 10)
	err = s.kafkaProducer.Publish(ctx, constants.KafkaTopicEmailVerification, key, payload)
	if err != nil {
		logger.WithContext(ctx).WithError(err).WithField("user_id", userID).Info("failed to publish email verification notification")
	}

	return nil
}

func (s *service) SignOut(ctx context.Context, sessionID string) error {
	return s.sessionRepository.Del(ctx, sessionID)
}

func (s *service) ForgotPassword(ctx context.Context, email string) error {
	normalizedEmail := helper.NormalizeEmail(email)
	user, err := s.userRepository.FindByEmail(ctx, normalizedEmail)
	if err != nil {
		if errors.Is(err, errConstant.ErrNotFound) {
			return errConstant.ErrUserNotFound
		}
		return err
	}

	verificationToken, err := helper.GenerateVerificationToken()
	if err != nil {
		return err
	}

	forgotPasswordBaseURL, err := url.Parse(s.configURL.ForgotPasswordBaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse forgot password URL: %w", err)
	}

	query := forgotPasswordBaseURL.Query()
	query.Set("token", verificationToken)
	forgotPasswordBaseURL.RawQuery = query.Encode()

	notification := event.ForgotPasswordRequested{
		EventID:          uuid.NewV4().String(),
		ReceiverID:       user.ID,
		ReceiverEmail:    normalizedEmail,
		Subject:          "パスワード再設定のご案内",
		Message:          "以下のリンクから新しいパスワードを設定してください: " + forgotPasswordBaseURL.String(),
		NotificationType: "EMAIL",
	}

	payload, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	hashedVerificationToken := helper.HashToken(verificationToken)
	expiresAt := time.Now().UTC().Add(time.Hour)

	token := &model.VerificationToken{
		UserID:    user.ID,
		TokenHash: hashedVerificationToken,
		TokenType: constants.TokenTypeForgotPassword,
		ExpiresAt: expiresAt,
	}

	if err := s.verificationTokenRepository.Create(ctx, token); err != nil {
		return err
	}

	key := strconv.FormatInt(user.ID, 10)
	if err := s.kafkaProducer.Publish(ctx, constants.KafkaTopicForgotPassword, key, payload); err != nil {
		logger.WithContext(ctx).WithError(err).WithField("user_id", user.ID).Info("failed to publish forgot password notification")
	}

	return nil
}

func (s *service) VerifyEmail(ctx context.Context, token string) error {
	trimmedToken := strings.TrimSpace(token)
	if trimmedToken == "" {
		return errConstant.ErrInvalidToken
	}

	hashedVerificationToken := helper.HashToken(trimmedToken)

	verificationToken, err := s.verificationTokenRepository.FindByTokenHash(ctx, hashedVerificationToken)
	if err != nil {
		if errors.Is(err, errConstant.ErrNotFound) {
			return errConstant.ErrInvalidToken
		}
		return err
	}

	if verificationToken.TokenType != constants.TokenTypeEmailVerification {
		return errConstant.ErrInvalidToken
	}

	now := time.Now().UTC()

	if !now.Before(verificationToken.ExpiresAt) {
		return errConstant.ErrInvalidToken
	}

	err = s.transactionManager.WithinTransaction(ctx, func(repositories *repository.Registry) error {
		if err := repositories.UserRepository.MarkAsVerified(ctx, verificationToken.UserID); err != nil {
			return err
		}

		if err := repositories.VerificationTokenRepository.DeleteByID(ctx, verificationToken.ID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil

}

func (s *service) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	hashedVerificationToken := helper.HashToken(input.Token)
	verificationToken, err := s.verificationTokenRepository.FindByTokenHash(ctx, hashedVerificationToken)
	if err != nil {
		if errors.Is(err, errConstant.ErrNotFound) {
			return errConstant.ErrInvalidToken
		}
		return err
	}

	if verificationToken.TokenType != constants.TokenTypeForgotPassword {
		return errConstant.ErrInvalidToken
	}

	now := time.Now().UTC()
	if !now.Before(verificationToken.ExpiresAt) {
		return errConstant.ErrInvalidToken
	}

	user, err := s.userRepository.FindByID(ctx, verificationToken.UserID)
	if err != nil {
		return err
	}

	if err := s.sessionRepository.DeleteByUserUUID(ctx, user.UUID, true, nil); err != nil {
		return err
	}

	hashedPassword, err := helper.HashPassword(input.NewPassword)
	if err != nil {
		return err
	}

	err = s.transactionManager.WithinTransaction(ctx, func(repositories *repository.Registry) error {
		if err := repositories.VerificationTokenRepository.DeleteByID(ctx, verificationToken.ID); err != nil {
			return err
		}

		if err := repositories.UserRepository.UpdateHashedPasswordByID(ctx, user.ID, hashedPassword); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
