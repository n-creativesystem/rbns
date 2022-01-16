package service

import (
	"context"
	"fmt"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/ncsfw/logger"
)

type Organization interface {
	Create(ctx context.Context, name, description string) (*model.Organization, error)
	FindById(ctx context.Context, strId string) (*model.Organization, error)
	FindByName(ctx context.Context, name string) (*model.Organization, error)
	FindAll(ctx context.Context) ([]model.Organization, error)
	Update(ctx context.Context, strId, name, description string) error
	Delete(ctx context.Context, strId string) error
}

type OrganizationImpl struct {
	log       logger.Logger
	telemetry telemetryFunc
}

var _ Organization = (*OrganizationImpl)(nil)

func NewOrganizationService(userService User) *OrganizationImpl {
	o := &OrganizationImpl{
		log:       logger.New("organization service"),
		telemetry: createSpanWithPrefix("organization service"),
	}
	return o
}

func (svc *OrganizationImpl) Create(ctx context.Context, name, description string) (out *model.Organization, e error) {
	svc.telemetry(ctx, "create", func(ctx context.Context) {
		orgName, err := model.NewName(name)
		if err != nil {
			e = err
			return
		}
		cmd := model.AddOrganizationCommand{
			Name:        orgName,
			Description: description,
		}
		if err := bus.Dispatch(ctx, &cmd); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
			e = err
			return
		}
		out, e = cmd.Result, nil
	})
	return
}

func (svc *OrganizationImpl) FindById(ctx context.Context, strId string) (out *model.Organization, e error) {
	svc.telemetry(ctx, "find by id", func(ctx context.Context) {
		id, err := model.NewID(strId)
		if err != nil {
			e = err
			return
		}
		query := model.GetOrganizationByIDQuery{
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

func (svc *OrganizationImpl) FindByName(ctx context.Context, name string) (out *model.Organization, e error) {
	svc.telemetry(ctx, "find by name", func(ctx context.Context) {
		orgName, err := model.NewName(name)
		if err != nil {
			e = err
			return
		}
		query := model.GetOrganizationByNameQuery{
			Name: orgName,
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

func (svc *OrganizationImpl) FindAll(ctx context.Context) (out []model.Organization, e error) {
	svc.telemetry(ctx, "find all", func(ctx context.Context) {
		query := model.GetOrganizationQuery{}
		if err := bus.Dispatch(ctx, &query); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "query", fmt.Sprintf("%+v", query))
			e = err
			return
		}

		out, e = query.Result, nil
	})
	return
}

func (svc *OrganizationImpl) Update(ctx context.Context, strId, name, description string) (e error) {
	svc.telemetry(ctx, "Update", func(ctx context.Context) {
		orgName, err := model.NewName(name)
		if err != nil {
			e = err
			return
		}
		org, err := svc.FindById(ctx, strId)
		if err != nil {
			e = err
			return
		}
		cmd := model.UpdateOrganizationCommand{
			PrimaryCommand: model.PrimaryCommand{
				ID: org.ID,
			},
			Name:        orgName,
			Description: description,
		}
		if err := bus.Dispatch(ctx, &cmd); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
			e = err
			return
		}
	})
	return
}

func (svc *OrganizationImpl) Delete(ctx context.Context, strId string) (e error) {
	svc.telemetry(ctx, "delete", func(ctx context.Context) {
		org, err := svc.FindById(ctx, strId)
		if err != nil {
			e = err
			return
		}
		cmd := model.DeleteOrganizationCommand{
			PrimaryCommand: model.PrimaryCommand{
				ID: org.ID,
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
