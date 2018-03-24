// Package pass provides access to UNIX password stores.
package pass

import (
	"io"
)

// Store is a password store.
type Store interface {
	Search(query string) ([]string, error)
	Open(item string) (io.ReadCloser, error)
	GlobSearch(query string) ([]string, error)
}
