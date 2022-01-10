package response

import "github.com/n-creativesystem/rbns/domain/model"

type Tenant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewTenants(tenants []model.Tenant) []Tenant {
	results := make([]Tenant, 0, len(tenants))
	for _, t := range tenants {
		result := Tenant{
			ID:   t.ID,
			Name: t.Name,
		}
		results = append(results, result)
	}
	return results
}

type Tenants []Tenant

func (t Tenants) BodyName() string {
	return "tenants"
}
