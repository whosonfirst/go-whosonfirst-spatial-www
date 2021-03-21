package static

import (
	"embed"
)

//go:embed css/* javascript/* images/*
var FS embed.FS
