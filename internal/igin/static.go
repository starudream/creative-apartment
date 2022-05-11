package igin

import (
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

type ServeFileSystem interface {
	http.FileSystem

	Exists(prefix string, path string) bool
}

type serveFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

func StaticFile(fs fs.FS, root string, indexes bool) *serveFileSystem {
	return &serveFileSystem{
		FileSystem: http.FS(fs),
		root:       root,
		indexes:    indexes,
	}
}

func (sys *serveFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(sys.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if stats.IsDir() {
			if !sys.indexes {
				index := path.Join(name, "index.html")
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

func Serve(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
	fileServer := http.FileServer(fs)
	if urlPrefix != "" {
		fileServer = http.StripPrefix(urlPrefix, fileServer)
	}
	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}
