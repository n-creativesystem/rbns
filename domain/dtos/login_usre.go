package dtos

type CurrentUser struct {
	// ID oauth provider id
	ID string
	// UserName oauth username
	UseName string
	// Email oauth email
	Email string
	// Role user role
	Role   string
	Groups []string

	IsSignedIn bool
}
