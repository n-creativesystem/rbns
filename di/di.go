package di

import (
	"go.uber.org/dig"
)

func MustInvoke(f interface{}) {
	if err := Invoke(f); err != nil {
		panic(err)
	}
}

func Invoke(f interface{}) error {
	return container.Invoke(f)
}

func MustRegister(f interface{}) {
	if err := Register(f); err != nil {
		panic(err)
	}
}

func Register(f interface{}) error {
	return container.Provide(f)
}

var (
	container = dig.New()
)
