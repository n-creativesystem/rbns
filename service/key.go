package service

import "github.com/n-creativesystem/rbns/domain/repository"

type ApiKey interface {
	Generate() (string, error)
	Get() string
}

type apiKey struct {
	repo repository.ApiKey
}

func NewApiKey(repo repository.ApiKey) ApiKey {
	return &apiKey{repo: repo}
}

func (key *apiKey) Generate() (string, error) {
	return key.repo.Generate()
}

func (key *apiKey) Get() string {
	return key.repo.Get()
}
