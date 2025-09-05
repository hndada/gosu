package resource

import "embed"

//go:embed *
var DefaultFS embed.FS
