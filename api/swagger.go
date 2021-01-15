package openapi

import (
	"embed"
)

//go:embed openapi.yml
var OpenAPI embed.FS
