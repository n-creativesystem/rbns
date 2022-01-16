package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/internal/apikeygen"
	"github.com/n-creativesystem/rbns/ncsfw/logger"
)

var (
	InValidAPIKey = errors.New("invalid API key")
)

type APIKeyError struct {
	Code    int
	Message string
	Error   error
}

type APIKey interface {
	Generate(email, tenant, strId string) (string, *APIKeyError)
	Decode(ctx context.Context, clientSecret string) (*model.LoginUser, *APIKeyError)
	GetLoginUser(ctx context.Context, email string) (*model.LoginUser, *APIKeyError)
	Add(ctx context.Context, tenant string, cmd *model.AddApiKeyCommand) (*apikeygen.KeyGenResult, *APIKeyError)
	Delete(ctx context.Context, id string) error
}

type APIKeyImpl struct {
	log logger.Logger
}

var _ APIKey = (*APIKeyImpl)(nil)

func NewAPIKey() *APIKeyImpl {
	return &APIKeyImpl{
		log: logger.New("apiKey service implements"),
	}
}

func (s *APIKeyImpl) Generate(email, tenant, name string) (string, *APIKeyError) {
	newKeyInfo, err := apikeygen.New(email, tenant, name)
	if err != nil {
		return "", &APIKeyError{
			Code:    500,
			Message: "Generating API key failed",
			Error:   err,
		}
	}
	return newKeyInfo.ClientSecret, nil
}

func (s *APIKeyImpl) Decode(ctx context.Context, clientSecret string) (*model.LoginUser, *APIKeyError) {
	decoded, err := apikeygen.Decode(clientSecret)
	if err != nil {
		return nil, &APIKeyError{
			Code:    401,
			Message: "APIKey decode",
			Error:   err,
		}
	}

	queryAPIKey := model.GetAPIKeyByNameQuery{
		KeyName: decoded.Name,
	}

	if err := bus.Dispatch(ctx, &queryAPIKey); err != nil {
		return nil, &APIKeyError{
			Code:    500,
			Message: "Get APIKey by name",
			Error:   err,
		}
	}

	apiKey := queryAPIKey.Result

	isValid, err := apikeygen.IsValid(decoded, apiKey.Key)
	if err != nil {
		s.log.ErrorWithContext(ctx, err, "Validating API key failed", "API key name", apiKey.Name)
		return nil, &APIKeyError{
			Code:    500,
			Message: "Validating API key failed",
			Error:   err,
		}
	}
	if !isValid {
		s.log.ErrorWithContext(ctx, InValidAPIKey, InValidAPIKey.Error(), "API key name", apiKey.Name)
		return nil, &APIKeyError{
			Code:    500,
			Message: InValidAPIKey.Error(),
			Error:   fmt.Errorf("API key name: %s", apiKey.Name),
		}
	}

	return s.GetLoginUser(ctx, apiKey.ServiceAccountEmail)
}

func (s *APIKeyImpl) GetLoginUser(ctx context.Context, email string) (*model.LoginUser, *APIKeyError) {
	queryLoginUser := model.GetLoginUserByEmailQuery{
		Email: email,
	}
	if err := bus.Dispatch(ctx, &queryLoginUser); err != nil {
		return nil, &APIKeyError{
			Code:    500,
			Message: "Get Login user",
			Error:   err,
		}
	}
	return queryLoginUser.Result, nil
}

func (s *APIKeyImpl) Add(ctx context.Context, tenant string, cmd *model.AddApiKeyCommand) (*apikeygen.KeyGenResult, *APIKeyError) {
	newKeyInfo, err := apikeygen.New(cmd.ServiceAccountEmail, tenant, cmd.Name)
	if err != nil {
		return nil, &APIKeyError{
			Code:    500,
			Message: "Generating API key failed",
			Error:   err,
		}
	}
	cmd.HashedKey = newKeyInfo.HashedKey
	if err := bus.Dispatch(ctx, &cmd); err != nil {
		if errors.Is(err, model.ErrAlreadyExists) {
			return nil, &APIKeyError{
				Code:    409,
				Message: "Conflict",
				Error:   err,
			}
		}
		return nil, &APIKeyError{
			Code:    500,
			Message: "Failed to add API Key",
			Error:   err,
		}
	}
	return &newKeyInfo, nil
}

func (s *APIKeyImpl) Delete(ctx context.Context, strId string) error {
	id, err := model.NewID(strId)
	if err != nil {
		return err
	}
	cmd := model.DeleteAPIKeyCommand{
		PrimaryCommand: model.PrimaryCommand{
			ID: id,
		},
	}
	if err := bus.Dispatch(ctx, &cmd); err != nil {
		return err
	}
	return nil
}
