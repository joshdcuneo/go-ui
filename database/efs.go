package database

import "embed"

//go:embed "migrations" "seeds"
var EFS embed.FS
