package model

import (
	"regexp"

	"github.com/n-creativesystem/rbns/ncsfw/logger"
)

type Tenant struct {
	ID          string
	Name        string
	Description string
}

var (
	tenantNameRegexp *regexp.Regexp
)

func init() {
	expr := `^[a-zA-Z0-9_+-]+(.[a-zA-Z0-9_+-]+)*$`
	var err error
	tenantNameRegexp, err = regexp.Compile(expr)
	if err != nil {
		logger.Panic(err, "tenant name regexp", "expr", expr)
	}
}

type TenantName string

func (t TenantName) String() string {
	return string(t)
}

func (t TenantName) Valid() error {
	if t.String() == "" {
		return ErrRequired
	}
	if !tenantNameRegexp.MatchString(t.String()) {
		return ErrRequired
	}

	return nil
}

type AddTenantCommand struct {
	Name        TenantName
	Description string
	OwnerEmail  string

	Result *Tenant
}

func (cmd *AddTenantCommand) Valid() error {
	if cmd.Name == "" {
		return ErrRequired
	}
	if cmd.OwnerEmail == "" {
		return ErrRequired
	}

	return nil
}

type DeleteTenantCommand struct {
	ID string
}

type GetTenantQuery struct {
	Result []*Tenant
}

type GetTenantByNameQuery struct {
	Name string

	Result *Tenant
}

type GetTenantByIdQuery struct {
	ID string

	Result *Tenant
}
