package www

import (
	"embed"
)

//go:embed css/* javascript/* intersects/* point-in-polygon/* point-in-polygon-with-tile/* *.html
var FS embed.FS
