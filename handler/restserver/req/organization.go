package req

type OrganizationCreateBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OrganizationUpdateBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
