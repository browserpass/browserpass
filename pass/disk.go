package pass

import (
	"errors"
	"io"
	"path/filepath"
	"os"
	"strings"

	"github.com/mattn/go-zglob"
	"os/user"
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
	usr, err := user.Current()

	if err != nil {
		return "", err
	}

	path := os.Getenv("PASSWORD_STORE_DIR")
	if path == "" {
		path = filepath.Join(usr.HomeDir, ".password-store")
	}

	// Follow symlinks
	return filepath.EvalSymlinks(path)
}

func (s *diskStore) Search(query string) ([]string, error) {
	// First, search for DOMAIN/USERNAME.gpg
	// Then, search for DOMAIN.gpg
	matches, err := zglob.Glob(s.path + "/**/" + query + "*/*.gpg")
	if err != nil {
		return nil, err
	}

	matches2, err := zglob.Glob(s.path + "/**/*" + query + "*.gpg")
	if err != nil {
		return nil, err
	}

        hostParts := strings.Split(query, ".")
        matches3 := []string{}
        if( len(hostParts) > 2 ) {
        matches3, err = zglob.Glob(s.path + "/**/*." + hostParts[1] + "." + hostParts[2] + "*.gpg")
		if err != nil {
			return nil, err
		}
        } 
        
	items := append(append(matches, matches2...), matches3...)
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
