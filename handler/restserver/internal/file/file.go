package file

import "net/http"

type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}
