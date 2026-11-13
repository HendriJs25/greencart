package customer

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"strconv"
	"user-service/internal/common/helper"
	"user-service/internal/constants"
	errConstant "user-service/internal/constants/error"
	"user-service/internal/domain/event"
	"user-service/internal/domain/model"
	queryinput "user-service/internal/domain/query"
	"user-service/internal/logger"
	"user-service/internal/messaging/kafka"
	"user-service/internal/query"
	"user-service/internal/repository"
	customerrepository "user-service/internal/repository/customer"
	userrepository "user-service/internal/repository/user"
	"uuid"
)

type service struct {
	customerRepository customerrepository.Repository
	userRepository     userrepository.Repository
	transactionManager repository.TransactionManager
	kafkaProducer      kafka.Producer
}

type Service interface {
	Create(ctx context.Context, input CreateCustomerInput) error
	GetAll(ctx context.Context, params queryinput.GetAllCustomerQuery) ([]GetAllCustomerResult, *query.Meta, error)
	GetDetail(ctx context.Context, id int64) (*GetDetailCustomerResult, error)
	UpdateData(ctx context.Context, id int64, input UpdateCustomerInput) error
	Delete(ctx context.Context, id int64) error
}

func NewService(customerRepository customerrepository.Repository, userRepository userrepository.Repository, transactionManager repository.TransactionManager, kafkaProducer kafka.Producer) Service {
	return &service{
		customerRepository: customerRepository,
		userRepository:     userRepository,
		transactionManager: transactionManager,
		kafkaProducer:      kafkaProducer,
	}
}

func (s *service) Create(ctx context.Context, input CreateCustomerInput) error {
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

	var formatedPhone *string
	if input.Phone != nil {
		tmp := fmt.Sprintf("+81%d", *input.Phone)
		formatedPhone = &tmp
	}

	var payload []byte
	var userID int64

	err = s.transactionManager.WithinTransaction(ctx, func(repositories *repository.Registry) error {
		customerRole, err := repositories.RoleRepository.FindByName(ctx, constants.Customer)
		if err != nil {
			if errors.Is(err, errConstant.ErrNotFound) {
				return fmt.Errorf("misconfigured: customer role not found")
			}
			return err
		}

		user := model.User{
			UUID:         uuid.NewV7(),
			Name:         input.Name,
			Email:        normalizedEmail,
			PasswordHash: hashedPassword,
			Phone:        formatedPhone,
			Address:      input.Address,
			Lat:          input.Lat,
			Lng:          input.Lng,
			IsVerified:   true,
			Roles:        []model.Role{*customerRole},
		}

		err = repositories.UserRepository.Create(ctx, &user)
		if err != nil {
			return err
		}

		notification := event.CreateCustomerRequested{
			EventID:          uuid.NewV4().String(),
			ReceiverID:       user.ID,
			ReceiverEmail:    normalizedEmail,
			Subject:          "【GreenCart】アカウント登録完了のお知らせ",
			Message:          fmt.Sprintf("以下のログイン情報をご確認のうえ、GreenCartにログインしてください。\n Email: %s \n Password: %s", normalizedEmail, input.Password),
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
	err = s.kafkaProducer.Publish(ctx, constants.KafkaTopicCreateCustomer, key, payload)
	if err != nil {
		logger.WithContext(ctx).WithError(err).WithField("user_id", userID).Info("failed to publish notification to customer")
	}
	return nil
}

func (s *service) GetAll(ctx context.Context, params queryinput.GetAllCustomerQuery) ([]GetAllCustomerResult, *query.Meta, error) {
	customers, total, err := s.customerRepository.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	metadata := query.Metadata(params.Page, params.Size, total)
	result := ToGetAllCustomerResult(customers)
	return result, metadata, nil
}

func (s *service) GetDetail(ctx context.Context, id int64) (*GetDetailCustomerResult, error) {
	customer, err := s.customerRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := ToGetDetailCustomerResult(customer)
	return result, nil
}

func (s *service) UpdateData(ctx context.Context, id int64, input UpdateCustomerInput) error {
	var formatedPhone *string
	if input.Phone != nil {
		tmp := fmt.Sprintf("+81%d", *input.Phone)
		formatedPhone = &tmp
	}

	updates := map[string]any{
		"name":        input.Name,
		"phone":       formatedPhone,
		"address":     input.Address,
		"lat":         input.Lat,
		"lng":         input.Lng,
		"is_verified": input.IsVerified,
	}

	err := s.customerRepository.UpdateData(ctx, id, updates)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.customerRepository.DeleteByID(ctx, id)
}
