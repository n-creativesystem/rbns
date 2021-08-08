package repository

import "github.com/n-creativesystem/rbns/domain/model"

type Organization interface {
	FindAll() (model.Organizations, error)
	FindByID(id model.ID) (*model.Organization, error)
	FindByName(name model.Name) (*model.Organization, error)
}

type OrganizationCommand interface {
	Organization
	Create(name model.Name, description string) (*model.Organization, error)
	Update(organization *model.Organization) error
	Delete(id model.ID) error
}
