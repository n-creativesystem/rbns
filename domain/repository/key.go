package repository

type ApiKey interface {
	Generate() (string, error)
	Get() string
}
