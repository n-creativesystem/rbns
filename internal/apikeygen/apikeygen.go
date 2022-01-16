package apikeygen

import (
	"encoding/json"
	"errors"

	"github.com/n-creativesystem/rbns/internal/utils"
	"github.com/n-creativesystem/rbns/utilsconv"
)

var ErrInvalidApiKey = errors.New("invalid API key")

type KeyGenResult struct {
	HashedKey    string
	ClientSecret string
}

type ApiKeyJson struct {
	Key    string `json:"k"`
	Name   string `json:"n"`
	Email  string `json:"id"`
	Tenant string `json:"tenant"`
}

func New(email, tenant string, name string) (KeyGenResult, error) {
	result := KeyGenResult{}

	jsonKey := ApiKeyJson{}
	jsonKey.Email = email
	jsonKey.Tenant = tenant
	jsonKey.Name = name
	var err error
	jsonKey.Key, err = utils.RandomString(32)
	if err != nil {
		return result, err
	}

	result.HashedKey, err = utils.EncodePassword(jsonKey.Key, name)
	if err != nil {
		return result, err
	}

	jsonString, err := json.Marshal(jsonKey)
	if err != nil {
		return result, err
	}
	result.ClientSecret = utilsconv.BytesToBase64(jsonString)
	return result, nil
}

func Decode(keyString string) (*ApiKeyJson, error) {
	jsonString, err := utilsconv.Base64ToByte(keyString)
	if err != nil {
		return nil, ErrInvalidApiKey
	}

	var keyObj ApiKeyJson
	err = json.Unmarshal(jsonString, &keyObj)
	if err != nil {
		return nil, ErrInvalidApiKey
	}

	return &keyObj, nil
}

func IsValid(key *ApiKeyJson, hashedKey string) (bool, error) {
	check, err := utils.EncodePassword(key.Key, key.Name)
	if err != nil {
		return false, err
	}
	return check == hashedKey, nil
}
