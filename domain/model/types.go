package model

import (
	"fmt"

	"github.com/n-creativesystem/rbns/infra/entity/plugins"
	"github.com/oklog/ulid/v2"
)

// type String interface {
// 	equals(v String) bool
// 	Err() error
// 	Value() plugins.ID
// 	value() *string
// 	String() string
// }

type ID interface {
	fmt.Stringer
}

type id string

// func (p id) Err() error {
// 	if p.value() == nil {
// 		return ErrRequired
// 	}
// 	return nil
// }

// func (p id) equals(v String) bool {
// 	if p.value() == nil {
// 		return false
// 	}
// 	if v.value() == nil {
// 		return false
// 	}
// 	return *p.value() == *v.value()
// }

func (p id) String() string {
	return string(p)
}

// func (p id) value() *string {
// 	if v := string(p); v == "" {
// 		return nil
// 	} else {
// 		return &v
// 	}
// }

func (p id) Value() plugins.ID {
	if v := string(p); v == "" {
		return plugins.ID("")
	} else {
		return plugins.ID(v)
	}
}

func newID(value string) (*id, error) {
	if value == "" {
		return nil, ErrRequired
	}
	_, err := ulid.Parse(value)
	if err != nil {
		return nil, err
	}
	id := id(value)
	return &id, nil
}

func NewID(value string) (ID, error) {
	return newID(value)
}

type Name interface {
	fmt.Stringer
}

type requiredString string

func (r requiredString) String() string {
	return string(r)
}

func newRequiredString(value string) (*requiredString, error) {
	if value == "" {
		return nil, ErrRequired
	}
	v := requiredString(value)
	return &v, nil
}

func newName(name string) (*requiredString, error) {
	return newRequiredString(name)
}

func newKey(key string) (*requiredString, error) {
	return newRequiredString(key)
}

func NewName(name string) (Name, error) {
	return newName(name)
}
