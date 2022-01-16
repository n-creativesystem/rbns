package model

type ApiKey struct {
	Id                  string
	Name                string
	Role                RoleLevel
	Key                 string
	ServiceAccountEmail string
}

type AddApiKeyCommand struct {
	Name                string    `json:"name" binding:"required"`
	Role                RoleLevel `json:"role" binding:"required"`
	HashedKey           string    `json:"-"`
	SecondsToLive       int64     `json:"secondsToLive"`
	ServiceAccountEmail string    `json:"-"`

	Result *ApiKey
}

type DeleteAPIKeyCommand struct {
	PrimaryCommand
}

type GetAPIKeyByNameQuery struct {
	KeyName string

	Result *ApiKey
}
