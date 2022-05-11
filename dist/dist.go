package dist

import (
	"embed"
	"io/fs"
	"path/filepath"
	"sort"

	"github.com/starudream/creative-apartment/internal/ierr"
)

var (
	//go:embed *.js *.html *.ico
	FS embed.FS

	Files []string
)

func init() {
	es, err := fs.ReadDir(FS, ".")
	ierr.CheckErr(err)
	for i := 0; i < len(es); i++ {
		Files = append(Files, es[i].Name())
	}
	sort.Slice(Files, func(i, j int) bool {
		ie, je := filepath.Ext(Files[i]), filepath.Ext(Files[j])
		if ie == je {
			return Files[i] < Files[j]
		} else {
			return ie < je
		}
	})
}
