package auth

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/url"
	"strconv"
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
	jwtservice "user-service/internal/services/jwt"
	"uuid"
)

type service struct {
	jwtService         jwtservice.Service
	userRepository     userrepository.Repository
	sessionRepository  sessionrepository.Repository
	transactionManager repository.TransactionManager
	kafkaProducer      kafka.Producer
	config             *config.Config
}

type Service interface {
	SignIn(ctx context.Context, input SignInInput) (*SignInResult, error)
	SignUp(ctx context.Context, input SignUpInput) error
}

func NewService(jwtService jwtservice.Service, userRepository userrepository.Repository, sessionRepository sessionrepository.Repository, transactionManager repository.TransactionManager, kafkaProducer kafka.Producer, config *config.Config) Service {
	return &service{
		jwtService:         jwtService,
		userRepository:     userRepository,
		sessionRepository:  sessionRepository,
		transactionManager: transactionManager,
		kafkaProducer:      kafkaProducer,
		config:             config,
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

	if !user.IsVerified {
		return nil, errConstant.ErrEmailNotVerified
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

	var userRoles []int64
	for _, role := range user.Roles {
		userRoles = append(userRoles, role.ID)
	}

	sessionCreatedAt := time.Now().UTC()

	sessionTTL := accessToken.ExpiresAt.Sub(sessionCreatedAt)

	session := model.Session{
		UUID:          user.UUID,
		CSRFTokenHash: hashedCSRFToken,
		Email:         user.Email,
		RoleID:        userRoles,
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

	verifyURL, err := url.Parse(s.config.VerifyBaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse verify URL: %w", err)
	}
	query := verifyURL.Query()
	query.Set("token", verificationToken)
	verifyURL.RawQuery = query.Encode()

	var notification event.EmailVerificationRequested
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

		notification = event.EmailVerificationRequested{
			EventID:          uuid.NewV4().String(),
			ReceiverID:       user.ID,
			ReceiverEmail:    normalizedEmail,
			Subject:          "メールアドレス認証のお願い",
			Message:          "以下のリンクをクリックして、アカウント認証を完了してください: " + verifyURL.String(),
			NotificationType: "EMAIL",
		}

		userID = user.ID

		return nil
	})
	if err != nil {
		return err
	}

	payload, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	key := strconv.FormatInt(userID, 10)

	err = s.kafkaProducer.Publish(ctx, constants.KafkaTopicEmailVerification, key, payload)
	if err != nil {
		logger.WithContext(ctx).Info("failed to publish email verification notification")
	}

	return nil
}
