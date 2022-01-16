package service

import "context"

type AccessCheck interface {
	Check(ctx context.Context)
}

type AccessCheckImpl struct {
	tenantService Tenant
	orgAggregate  OrganizationAggregation
}

func NewAccessCheck(tenantService Tenant, orgAggregate OrganizationAggregation) *AccessCheckImpl {
	return &AccessCheckImpl{
		tenantService: tenantService,
		orgAggregate:  orgAggregate,
	}
}
