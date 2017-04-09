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
	// Number of prefix matches to be checked is the number of periods
	// with lower bound of 1
	count := strings.Count(query, ".")
	if count < 1 {
		count = 1
	}

	partials := make([]string, count)
	for i:=0; i < count; i++ {
		partials[i] = query
		j := strings.Index(query, ".")
		if j >= 0 {
			query = query[j+1:]
		}
	}

	// First, search for DOMAIN/USERNAME.gpg
	// Then, search for DOMAIN.gpg
	matches, err := zglob.Glob(s.path + "/**/*" + partials[count-1] + "*/*.gpg")
	if err != nil {
		return nil, err
	}

	matches2, err := zglob.Glob(s.path + "/**/*" + partials[count-1] + "*.gpg")
	if err != nil {
		return nil, err
	}

	items := append(matches, matches2...)
	for i, path := range items {
		item, err := filepath.Rel(s.path, path)
		if err != nil {
			return nil, err
		}
		items[i] = strings.TrimSuffix(item, ".gpg")
	}

	for _, partial := range partials {
		ret := []string{}
		for _, item := range items {
			if strings.HasPrefix(item, partial) {
				ret = append(ret, item)
			}
		}
		if len(ret) > 0 {
			return ret, nil
		}
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
