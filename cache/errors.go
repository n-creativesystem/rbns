package cache

import "errors"

var (
	// ErrNotFound データが見つかりません。
	ErrNotFound = errors.New("データが見つかりません。")
)
