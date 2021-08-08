package restserver

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/n-creativesystem/rbns/handler/restserver/internal/file"
)

type indexHandler struct {
	*readSeeker
	stats fs.FileInfo
}

var _ io.Closer = (*indexHandler)(nil)
var _ io.Reader = (*indexHandler)(nil)

func (h *indexHandler) Readdir(count int) ([]fs.FileInfo, error) {
	return nil, nil
}

func (h *indexHandler) Stat() (fs.FileInfo, error) {
	return h.stats, nil
}

func (h *indexHandler) Close() error {
	return nil
}

const INDEX = "index.html"

var (
	basePathPattern = regexp.MustCompile(`<base href="/"`)
	basePathReplace = `<base href="%s/"`
)

type localFileSystem struct {
	fs          http.FileSystem
	root        string
	indexes     bool
	indexHandle *indexHandler
}

func (l *localFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(l.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if stats.IsDir() {
			if !l.indexes {
				index := path.Join(name, INDEX)
				_, err := os.Stat(index)
				if err != nil {
					return false
				}
			}
		}
		return true
	}
	return false
}

func newFileSystem(root string, indexes bool, baseURL string) file.ServeFileSystem {
	fs := &localFileSystem{
		fs:      http.Dir(root),
		root:    root,
		indexes: indexes,
	}
	indexFileName := path.Join(root, "index.html")
	stats, err := os.Stat(indexFileName)
	if err != nil {
		return fs
	}
	indexBytes, err := os.ReadFile(indexFileName)
	if err != nil {
		return fs
	}
	if baseURL != "/" {
		if strings.HasPrefix(baseURL, "/") && !strings.HasSuffix(baseURL, "/") {
			indexBytes = basePathPattern.ReplaceAll(indexBytes, []byte(fmt.Sprintf(basePathReplace, baseURL)))
		}
	}
	buffer := bytes.NewBuffer(indexBytes)
	indexHandle := &indexHandler{
		stats: stats,
		readSeeker: &readSeeker{
			Reader: buffer,
			buf:    &bytes.Buffer{},
			max:    DefaultMaxBufferSize,
		},
	}
	fs.indexHandle = indexHandle
	return fs
}

func (l *localFileSystem) Open(name string) (http.File, error) {
	// return l.fs.Open(name)
	if l.indexHandle == nil {
		return l.fs.Open(name)
	}
	if name == "/index.html" {
		return l.indexHandle, nil
	} else {
		return l.fs.Open(name)
	}
}
