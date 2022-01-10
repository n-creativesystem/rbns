package service

import (
	"context"
	"fmt"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/logger"
)

type UserService interface {
	Create(ctx context.Context, userID string, name string) error
	Delete(ctx context.Context, userID string) error
	FindById(ctx context.Context, userID string) (*model.User, error)
	FindByIds(ctx context.Context, userIDs []string) ([]model.User, error)
}

type userService struct {
	log logger.Logger
}

var _ UserService = (*userService)(nil)

func NewUserService() UserService {
	return &userService{logger.New("user service")}
}

// User
func (svc *userService) Create(ctx context.Context, userID string, name string) error {
	id, err := model.NewID(userID)
	if err != nil {
		return err
	}
	cmd := model.AddUserCommand{
		PrimaryCommand: model.PrimaryCommand{
			ID: id,
		},
	}
	if err := bus.Dispatch(ctx, &cmd); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
		return err
	}
	return nil
}

func (svc *userService) Delete(ctx context.Context, userID string) error {
	id, err := model.NewID(userID)
	if err != nil {
		return err
	}
	cmd := model.DeleteUserCommand{
		PrimaryCommand: model.PrimaryCommand{
			ID: id,
		},
	}
	if err := bus.Dispatch(ctx, &cmd); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
		return err
	}
	return nil
}

func (svc *userService) FindById(ctx context.Context, userID string) (*model.User, error) {
	id, err := model.NewID(userID)
	if err != nil {
		return nil, err
	}
	query := model.GetUserByIDQuery{
		PrimaryQuery: model.PrimaryQuery{
			ID: id,
		},
	}
	if err := bus.Dispatch(ctx, &query); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "query", fmt.Sprintf("%+v", query))
		return nil, err
	}
	return query.Result, nil
}

func (svc *userService) FindByIds(ctx context.Context, userIDs []string) ([]model.User, error) {
	query := model.GetUserByIDsQuery{
		Query: make([]model.PrimaryQuery, 0, len(userIDs)),
	}
	for _, userID := range userIDs {
		id, err := model.NewID(userID)
		if err != nil {
			continue
		}
		query.Query = append(query.Query, model.PrimaryQuery{
			ID: id,
		})
	}
	if err := bus.Dispatch(ctx, &query); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "query", fmt.Sprintf("%+v", query))
		return nil, err
	}
	return query.Result, nil
}
