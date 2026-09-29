// Package migrations exposes the SQL migration files as an embedded filesystem
// so the server can apply them without depending on the process working
// directory (which differs between `go run ./cmd/server`, a built binary and a
// Docker image).
package migrations

import "embed"

// FS holds every migration. Filenames carry their execution order (001_, 002_,
// ...), so lexical sorting is the migration order.
//
//go:embed *.sql
var FS embed.FS