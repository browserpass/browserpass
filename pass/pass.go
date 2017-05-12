// Package pass provides access to UNIX password stores.
package pass

import (
	"errors"
	"io"
)

// ErrNotFound is returned by Store.Open if the requested item is not found.
var ErrNotFound = errors.New("pass: not found")

// Store is a password store.
type Store interface {
	Search(query string) ([]string, error)
	Lookup(query string) ([]string, error)
	Open(item string) (io.ReadCloser, error)
}
