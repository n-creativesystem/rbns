package service

import (
	"context"

	"github.com/n-creativesystem/rbns/di"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/domain/repository"
)

func init() {
	di.MustRegister(newOrganizationService)
}

type OrganizationService interface {
	Create(ctx context.Context, name, description string) (*model.Organization, error)
	FindById(ctx context.Context, strId string) (*model.Organization, error)
	FindAll(ctx context.Context) (model.Organizations, error)
	Update(ctx context.Context, strId, name, description string) error
	Delete(ctx context.Context, strId string) error
}

type organizationService struct {
	repo repository.Repository
}

var _ OrganizationService = (*organizationService)(nil)

func newOrganizationService(repo repository.Repository) OrganizationService {
	return &organizationService{repo: repo}
}

// Organization
func (srv *organizationService) Create(ctx context.Context, name, description string) (*model.Organization, error) {
	var out model.Organization
	orgName, err := model.NewName(name)
	if err != nil {
		return nil, err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	err = tx.Do(func(tx repository.Transaction) error {
		orgRepo := tx.Organization()
		org, err := orgRepo.Create(orgName, description)
		if err != nil {
			return err
		}
		out = *org
		return nil
	})
	return &out, err
}

func (srv *organizationService) FindById(ctx context.Context, strId string) (*model.Organization, error) {
	orgRepo := srv.repo.NewConnection().Organization(ctx)
	id, err := model.NewID(strId)
	if err != nil {
		return nil, err
	}
	organization, err := orgRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return organization, nil
}

func (srv *organizationService) FindAll(ctx context.Context) (model.Organizations, error) {
	orgRepo := srv.repo.NewConnection().Organization(ctx)
	organizations, err := orgRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return organizations, nil
}

func (srv *organizationService) Update(ctx context.Context, strId, name, description string) error {
	mOrg, err := model.NewOrganization(strId, name, description)
	if err != nil {
		return err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	if err := tx.Do(func(tx repository.Transaction) error { return tx.Organization().Update(mOrg) }); err != nil {
		return err
	}
	return nil
}

func (srv *organizationService) Delete(ctx context.Context, strId string) error {
	id, err := model.NewID(strId)
	if err != nil {
		return err
	}
	tx := srv.repo.NewConnection().Transaction(ctx)
	if err := tx.Do(func(tx repository.Transaction) error { return tx.Organization().Delete(id) }); err != nil {
		return err
	}
	return nil
}
