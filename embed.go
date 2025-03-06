package embed

import (
	"embed"
)

var (
	//go:embed doc/api
	EmbedAPIAssets embed.FS
)
