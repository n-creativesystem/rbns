package service

import (
	"context"
	"fmt"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/logger"
)

type OrganizationService interface {
	Create(ctx context.Context, name, description string) (*model.Organization, error)
	FindById(ctx context.Context, strId string) (*model.Organization, error)
	FindByName(ctx context.Context, name string) (*model.Organization, error)
	FindAll(ctx context.Context) ([]model.Organization, error)
	Update(ctx context.Context, strId, name, description string) error
	Delete(ctx context.Context, strId string) error
}

type organizationService struct {
	log logger.Logger
}

var _ OrganizationService = (*organizationService)(nil)

func NewOrganizationService(userService UserService) OrganizationService {
	return &organizationService{
		log: logger.New("organization service"),
	}
}

func (svc *organizationService) Create(ctx context.Context, name, description string) (*model.Organization, error) {
	orgName, err := model.NewName(name)
	if err != nil {
		return nil, err
	}
	cmd := model.AddOrganizationCommand{
		Name:        orgName,
		Description: description,
	}
	if err := bus.Dispatch(ctx, &cmd); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
		return nil, err
	}
	return cmd.Result, nil
}

func (svc *organizationService) FindById(ctx context.Context, strId string) (*model.Organization, error) {
	id, err := model.NewID(strId)
	if err != nil {
		return nil, err
	}
	query := model.GetOrganizationByIDQuery{
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

func (svc *organizationService) FindByName(ctx context.Context, name string) (*model.Organization, error) {
	orgName, err := model.NewName(name)
	if err != nil {
		return nil, err
	}
	query := model.GetOrganizationByNameQuery{
		Name: orgName,
	}
	if err := bus.Dispatch(ctx, &query); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "query", fmt.Sprintf("%+v", query))
		return nil, err
	}
	return query.Result, nil
}

func (svc *organizationService) FindAll(ctx context.Context) ([]model.Organization, error) {
	query := model.GetOrganizationQuery{}
	if err := bus.Dispatch(ctx, &query); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "query", fmt.Sprintf("%+v", query))
		return nil, err
	}
	return query.Result, nil
}

func (svc *organizationService) Update(ctx context.Context, strId, name, description string) error {
	orgName, err := model.NewName(name)
	if err != nil {
		return err
	}
	org, err := svc.FindById(ctx, strId)
	if err != nil {
		return err
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
		return err
	}
	return nil
}

func (svc *organizationService) Delete(ctx context.Context, strId string) error {
	org, err := svc.FindById(ctx, strId)
	if err != nil {
		return err
	}
	cmd := model.DeleteOrganizationCommand{
		PrimaryCommand: model.PrimaryCommand{
			ID: org.ID,
		},
	}
	if err := bus.Dispatch(ctx, &cmd); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
		return err
	}
	return nil
}
