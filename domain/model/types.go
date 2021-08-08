package model

import (
	"strings"

	"github.com/oklog/ulid/v2"
)

type String interface {
	equals(v String) bool
	Err() error
	Value() *string
}

type ID interface {
	String
}

type id string

func (p id) Err() error {
	if p.Value() == nil {
		return ErrRequired
	}
	return nil
}

func (p id) equals(v String) bool {
	if p.Value() == nil {
		return false
	}
	if v.Value() == nil {
		return false
	}
	return *p.Value() == *v.Value()
}

func (p id) String() string {
	v := p.Value()
	if v == nil {
		return ""
	}
	return *v
}

func (p id) Value() *string {
	if v := string(p); v == "" {
		return nil
	} else {
		return &v
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
	String
}

type Key interface {
	String
}

type requiredString string

func (r requiredString) Err() error {
	if r.Value() == nil {
		return ErrRequired
	}
	return nil
}

func (r requiredString) equals(v String) bool {
	return strings.EqualFold(*r.Value(), *v.Value())
}

func (r requiredString) Value() *string {
	v := string(r)
	if v == "" {
		return nil
	}
	return &v
}

func newName(name string) (*requiredString, error) {
	if name == "" {
		return nil, ErrRequired
	}
	v := requiredString(name)
	return &v, nil
}

func newKey(key string) (*requiredString, error) {
	return newName(key)
}

func NewName(name string) (Name, error) {
	return newName(name)
}

func NewKey(key string) (Key, error) {
	return newKey(key)
}

type model struct {
	id   id
	name requiredString
}

func (m *model) GetID() *string {
	return m.id.Value()
}

func (m *model) GetName() *string {
	return m.name.Value()
}

func (m *model) equals(value model) bool {
	if m.id.equals(value.id) {
		return true
	}
	return m.name.equals(value.name)
}
