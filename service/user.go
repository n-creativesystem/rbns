package service

import (
	"context"
	"fmt"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/ncsfw/logger"
)

type User interface {
	Create(ctx context.Context, userID string, name string) error
	Delete(ctx context.Context, userID string) error
	FindById(ctx context.Context, userID string) (*model.User, error)
	FindByIds(ctx context.Context, userIDs []string) ([]model.User, error)
}

type UserImpl struct {
	log       logger.Logger
	telemetry telemetryFunc
}

var _ User = (*UserImpl)(nil)

func NewUserService() *UserImpl {
	return &UserImpl{
		log:       logger.New("user service"),
		telemetry: createSpanWithPrefix("user service"),
	}
}

// User
func (svc *UserImpl) Create(ctx context.Context, userID string, name string) (e error) {
	svc.telemetry(ctx, "create", func(ctx context.Context) {
		cmd := model.AddUserCommand{
			ID:   userID,
			Name: name,
		}
		if err := bus.Dispatch(ctx, &cmd); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
			e = err
			return
		}
	})
	return
}

func (svc *UserImpl) Delete(ctx context.Context, userID string) (e error) {
	svc.telemetry(ctx, "delete", func(ctx context.Context) {
		id, err := model.NewID(userID)
		if err != nil {
			e = err
			return
		}
		cmd := model.DeleteUserCommand{
			PrimaryCommand: model.PrimaryCommand{
				ID: id,
			},
		}
		if err := bus.Dispatch(ctx, &cmd); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
			e = err
			return
		}
	})
	return
}

func (svc *UserImpl) FindById(ctx context.Context, userID string) (out *model.User, e error) {
	svc.telemetry(ctx, "find by id", func(ctx context.Context) {
		id, err := model.NewID(userID)
		if err != nil {
			e = err
			return
		}
		query := model.GetUserByIDQuery{
			PrimaryQuery: model.PrimaryQuery{
				ID: id,
			},
		}
		if err := bus.Dispatch(ctx, &query); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "query", fmt.Sprintf("%+v", query))
			e = err
			return
		}
		out, e = query.Result, nil
	})
	return
}

func (svc *UserImpl) FindByIds(ctx context.Context, userIDs []string) (out []model.User, e error) {
	svc.telemetry(ctx, "find by ids", func(ctx context.Context) {
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
			e = err
			return
		}
		out, e = query.Result, nil
	})
	return
}
