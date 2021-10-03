package storage

import (
	"io"

	"github.com/n-creativesystem/rbns/domain/repository"
	"github.com/sirupsen/logrus"
)

type Factory interface {
	Initialize(settings map[string]interface{}, logger *logrus.Logger) error
	Reader() repository.Reader
	Writer() repository.Writer
	io.Closer
}
