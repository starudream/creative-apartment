package dist

import (
	"embed"
)

//go:embed *.js *.html
var FS embed.FS
