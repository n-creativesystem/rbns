package res

import "github.com/n-creativesystem/rbns/domain/model"

type Organization struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Users       []User `json:"users"`
}

type Organizations struct {
	Organizations []Organization `json:"organizations"`
}

type OrganizationUser struct {
	UserKey                 string `json:"userKey"`
	OrganizationId          string `json:"organizationId"`
	OrganizationName        string `json:"organizationName"`
	OrganizationDescription string `json:"organizationDescription"`
}

func NewOrganization(org model.Organization) Organization {
	return Organization{
		Id:          *org.GetID(),
		Name:        *org.GetName(),
		Description: org.GetDescription(),
		Users:       NewUsers(org.GetUsers()),
	}
}

func NewOrganizations(orgs model.Organizations) Organizations {
	os := make([]Organization, len(orgs))
	for idx, o := range orgs {
		os[idx] = NewOrganization(o)
	}
	return Organizations{
		Organizations: os,
	}
}
