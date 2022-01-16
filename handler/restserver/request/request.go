package request

type TenantPost struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
