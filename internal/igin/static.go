package igin

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/starudream/creative-apartment/internal/ierr"
)

type ServeFileSystem interface {
	http.FileSystem
	Exists(path string) bool
}

func Serve(urlPrefix string, sfs ServeFileSystem, notFound gin.HandlerFunc) gin.HandlerFunc {
	fileServer := http.FileServer(sfs)
	if urlPrefix != "" {
		fileServer = http.StripPrefix(urlPrefix, fileServer)
	}
	return func(c *gin.Context) {
		if sfs.Exists(c.Request.URL.Path) {
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		} else {
			notFound(c)
		}
	}
}

type localFileSystem struct {
	http.FileSystem
	paths map[string]struct{}
}

var _ ServeFileSystem = (*localFileSystem)(nil)

func (sfs *localFileSystem) Exists(path string) bool {
	_, exist := sfs.paths[path]
	return exist
}

func LocalFileSystem(localPath string) *localFileSystem {
	paths := map[string]struct{}{}
	ierr.CheckErr(filepath.Walk(localPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			reqPath := strings.TrimPrefix(path, localPath)
			idxPath := strings.TrimSuffix(reqPath, "index.html")
			if reqPath != idxPath {
				paths[idxPath] = struct{}{}
			}
			paths[reqPath] = struct{}{}
		}
		return err
	}))
	var ps []string
	for p := range paths {
		if strings.HasSuffix(p, "/") {
			continue
		}
		ps = append(ps, p)
	}
	sort.Strings(ps)
	log.Debug().Msgf("[http] load static files: %s", strings.Join(ps, ", "))
	return &localFileSystem{http.Dir(localPath), paths}
}
