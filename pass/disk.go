package pass

import (
	"errors"
	"io"
	"path/filepath"
	"os"
	"strings"

	"github.com/mattn/go-zglob"
)

type diskStore struct {
	path string
}

func NewDefaultStore() (Store, error) {
	path, err := defaultStorePath()
	if err != nil {
		return nil, err
	}

	return &diskStore{path}, nil
}

func defaultStorePath() (string, error) {
	path := os.Getenv("PASSWORD_STORE_DIR")
	if path == "" {
		path = filepath.Join(os.Getenv("HOME"), ".password-store")
	}

	// Follow symlinks
	return filepath.EvalSymlinks(path)
}

func (s *diskStore) Search(query string) ([]string, error) {
	var matches, matches2 []string
	var err error

	for {
		// First, search for DOMAIN/USERNAME.gpg
		// Then, search for DOMAIN.gpg
		matches, err = zglob.Glob(s.path + "/**/" + query + "*/*.gpg")
		if err != nil {
			return nil, err
		}

		matches2, err = zglob.Glob(s.path + "/**/" + query + "*.gpg")
		if err != nil {
			return nil, err
		}

		// Exit loop if some results were found
		if len(matches) > 0 || len(matches2) > 0 {
			break
		}

		// Remove component from query for next iteration
		i := strings.Index(query, ".")
		if i >= 0 {
			query = query[i+1:len(query)]
		}

		// Exit loop if we are at tld
		if strings.Index(query, ".") == -1 {
			break
		}
	}

	items := append(matches, matches2...)
	for i, path := range items {
		item, err := filepath.Rel(s.path, path)
		if err != nil {
			return nil, err
		}
		items[i] = strings.TrimSuffix(item, ".gpg")
	}

	return items, nil
}

func (s *diskStore) Open(item string) (io.ReadCloser, error) {
	p := filepath.Join(s.path, item+".gpg")
	if !filepath.HasPrefix(p, s.path) {
		// Make sure the requested item is *in* the password store
		return nil, errors.New("invalid item path")
	}

	f, err := os.Open(p)
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	}
	return f, err
}
