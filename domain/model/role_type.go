package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/n-creativesystem/rbns/utilsconv"
)

// type RoleType string
type RoleLevel uint

const (
	ROLE_VIEWER RoleLevel = iota + 1
	ROLE_EDITOR
	ROLE_ADMIN
)

func (r RoleLevel) Valid() bool {
	return r == ROLE_VIEWER || r == ROLE_ADMIN || r == ROLE_EDITOR
}

func (r RoleLevel) String() string {
	switch r {
	case ROLE_VIEWER:
		return "Viewer"
	case ROLE_EDITOR:
		return "Editor"
	case ROLE_ADMIN:
		return "Admin"
	default:
		panic("no support role level")
	}
}

func String2RoleLevel(str string) (RoleLevel, error) {
	switch strings.ToLower(str) {
	case "viewer":
		return ROLE_VIEWER, nil
	case "editor":
		return ROLE_EDITOR, nil
	case "admin":
		return ROLE_ADMIN, nil
	}
	return RoleLevel(0), fmt.Errorf("invalid role value: %s", str)
}

func (r *RoleLevel) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	lvl, err := String2RoleLevel(str)
	if err != nil {
		return fmt.Errorf("JSON validation error: %s", err.Error())
	}
	*r = lvl
	return nil
}

func (r RoleLevel) MarshalJSON() ([]byte, error) {
	return utilsconv.StringToBytes(r.String()), nil
}

// IsLevelEnabled 引数のロールレベルの方が大きい場合trueを返す
func (r RoleLevel) IsLevelEnabled(lvl RoleLevel) bool {
	return lvl >= r
}
