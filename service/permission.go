package service

import (
	"context"
	"fmt"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/logger"
)

type PermissionService interface {
	Create(ctx context.Context, names, descriptions []string) ([]model.Permission, error)
	FindById(ctx context.Context, strId string) (*model.Permission, error)
	FindByIds(ctx context.Context, ids []string) ([]model.Permission, error)
	FindAll(ctx context.Context) ([]model.Permission, error)
	Update(ctx context.Context, strId, name, description string) (*model.Permission, error)
	Delete(ctx context.Context, strId string) error
}
type permissionService struct {
	log       logger.Logger
	telemetry telemetryFunc
}

func NewPermissionService() PermissionService {
	return &permissionService{
		log:       logger.New("permission service"),
		telemetry: createSpanWithPrefix("permission service"),
	}
}

// Permission
func (svc *permissionService) Create(ctx context.Context, names, descriptions []string) (out []model.Permission, e error) {
	svc.telemetry(ctx, "Create", func(ctx context.Context) {
		cmd := model.AddPermissionCommands{
			AddPermissions: make([]model.AddPermissionCommand, 0, 100),
		}
		for idx, name := range names {
			n, err := model.NewName(name)
			if err != nil {
				svc.log.ErrorWithContext(ctx, err, "name constructor error", "constructor name", name)
				continue
			}
			cmd.AddPermissions = append(cmd.AddPermissions, model.AddPermissionCommand{
				Name:        n,
				Description: descriptions[idx],
			})
		}
		if err := bus.Dispatch(ctx, &cmd); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
			e = err
			return
		}
		for _, permission := range cmd.AddPermissions {
			out = append(out, *permission.Result)
		}
	})
	return
}

func (svc *permissionService) FindById(ctx context.Context, strId string) (out *model.Permission, e error) {
	svc.telemetry(ctx, "FindById", func(ctx context.Context) {
		id, err := model.NewID(strId)
		if err != nil {
			e = err
			return
		}
		query := model.GetPermissionByIDQuery{
			PrimaryCommand: model.PrimaryCommand{
				ID: id,
			},
		}
		if err := bus.Dispatch(ctx, &query); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "query", fmt.Sprintf("%+v", query))
			e = err
			return
		}
		out = query.Result
	})
	return
}

func (svc *permissionService) FindByIds(ctx context.Context, ids []string) (out []model.Permission, e error) {
	svc.telemetry(ctx, "FindByIds", func(ctx context.Context) {
		query := model.GetPermissionByIDsQuery{
			Query: make([]model.PrimaryCommand, 0, len(ids)),
		}
		for _, strId := range ids {
			id, err := model.NewID(strId)
			if err != nil {
				continue
			}
			query.Query = append(query.Query, model.PrimaryCommand{
				ID: id,
			})
		}
		if err := bus.Dispatch(ctx, &query); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "query", fmt.Sprintf("%+v", query))
			e = err
			return
		}
		out = query.Result
	})
	return
}

func (svc *permissionService) FindAll(ctx context.Context) (out []model.Permission, e error) {
	svc.telemetry(ctx, "FindAll", func(ctx context.Context) {
		query := model.GetPermissionQuery{}
		if err := bus.Dispatch(ctx, &query); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "query", fmt.Sprintf("%+v", query))
			e = err
			return
		}
		out = query.Result
	})
	return
}

func (svc *permissionService) Update(ctx context.Context, strId, name, description string) (out *model.Permission, e error) {
	svc.telemetry(ctx, "Update", func(ctx context.Context) {
		id, err := model.NewID(strId)
		if err != nil {
			e = err
			return
		}
		n, err := model.NewName(name)
		if err != nil {
			e = err
			return
		}
		cmd := model.UpdatePermissionCommand{
			PrimaryCommand: model.PrimaryCommand{
				ID: id,
			},
			Name:        n,
			Description: description,
		}
		if err := bus.Dispatch(ctx, &cmd); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
			e = err
			return
		}
		query := model.GetPermissionByIDQuery{
			PrimaryCommand: model.PrimaryCommand{
				ID: id,
			},
		}
		if err := bus.Dispatch(ctx, &query); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "query", fmt.Sprintf("%+v", query))
			e = err
			return
		}
		out = query.Result
	})
	return
}

func (svc *permissionService) Delete(ctx context.Context, strId string) (e error) {
	svc.telemetry(ctx, "Delete", func(ctx context.Context) {
		id, err := model.NewID(strId)
		if err != nil {
			e = err
			return
		}
		cmd := model.DeletePermissionCommand{
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
