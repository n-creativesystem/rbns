package model

import (
	"encoding/json"
	"fmt"
)

type RoleType string

const (
	ROLE_VIEWER RoleType = "Viewer"
	ROLE_EDITOR RoleType = "Editor"
	ROLE_ADMIN  RoleType = "Admin"
)

func (r RoleType) Valid() bool {
	return r == ROLE_VIEWER || r == ROLE_ADMIN || r == ROLE_EDITOR
}

func (r RoleType) Includes(other RoleType) bool {
	if r == ROLE_ADMIN {
		return true
	}

	if r == ROLE_EDITOR {
		return other != ROLE_ADMIN
	}

	return r == other
}

func (r RoleType) Children() []RoleType {
	switch r {
	case ROLE_ADMIN:
		return []RoleType{ROLE_EDITOR, ROLE_VIEWER}
	case ROLE_EDITOR:
		return []RoleType{ROLE_VIEWER}
	default:
		return nil
	}
}

func (r *RoleType) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	*r = RoleType(str)

	if !r.Valid() {
		if (*r) != "" {
			return fmt.Errorf("JSON validation error: invalid role value: %s", *r)
		}

		*r = ROLE_VIEWER
	}

	return nil
}
