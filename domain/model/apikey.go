package model

type ApiKey struct {
	Id               string
	Name             string
	Role             RoleType
	Key              string
	ServiceAccountID string
}

type AddApiKeyCommand struct {
	Name             string   `json:"name" binding:"required"`
	Role             RoleType `json:"role" binding:"required"`
	HashedKey        string   `json:"-"`
	SecondsToLive    int64    `json:"secondsToLive"`
	ServiceAccountID string   `json:"-"`

	Result *ApiKey
}

type DeleteAPIKeyCommand struct {
	PrimaryCommand
}

type GetAPIKeyByNameQuery struct {
	KeyName string

	Result *ApiKey
}
