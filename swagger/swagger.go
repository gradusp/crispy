package swagger

import (
	"embed"
)

//go:embed openapi.yml
var OpenAPI embed.FS

//go:embed ui/*
var UI embed.FS
