package repository

import "context"

type Repository interface {
	NewConnection() Reader
}

type Reader interface {
	Permission(ctx context.Context) Permission
	Role(ctx context.Context) Role
	Organization(ctx context.Context) Organization
	User(ctx context.Context) User
	Resource(ctx context.Context) Resource
	// Transaction(ctx context.Context) Tx
}

type Writer interface {
	Do(ctx context.Context, fn func(tx Transaction) error) error
}

type Transaction interface {
	Permission() PermissionCommand
	Role() RoleCommand
	Organization() OrganizationCommand
	User() UserCommand
	Resource() ResourceCommand
}
