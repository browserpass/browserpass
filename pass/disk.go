package pass

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"os/user"

	"github.com/mattn/go-zglob"
	sfuzzy "github.com/sahilm/fuzzy"
)

type diskStore struct {
	path           string
	useFuzzy       bool   // Setting to use the fuzzy matcher or the legacy glob matching
	fuzzyAlgorithm string // Setting to choose the fuzzy algorithm
}

func NewDefaultStore() (Store, error) {
	path, err := defaultStorePath()
	if err != nil {
		return nil, err
	}

	return &diskStore{path, false, ""}, nil
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

// Set the configuration options for password matching.
func (s *diskStore) SetConfig(path *string, use_fuzzy *bool) error {
	if path != nil {
		//todo validate path exists
		s.path = *path
	}
	if use_fuzzy != nil {
		s.useFuzzy = *use_fuzzy
	}
	return nil
}

// Do a search. Will call into the correct algoritm (glob, fuzzy (renstrom) or fuzzy (sahilm)
// based on the settings present in the diskStore struct
func (s *diskStore) Search(query string) ([]string, error) {
	// default legacy glob search
	if !s.useFuzzy {
		return s.GlobSearch(query)
	} else {
		return s.sahilmFuzzySerach(query)
	}
}

// Fuzzy searches first get a list of all pass entries by doing a glob search
// for the empty string, then apply appropriate logic to convert results to
// a slice of strings, finally returning all of the unique entries.
func (s *diskStore) sahilmFuzzySerach(query string) ([]string, error) {
	var items []string
	file_list, err := s.GlobSearch("")
	if err != nil {
		return nil, err
	}

	// The resulting match struct does not copy the strings, but rather
	// provides the index to the original array. Copy those strings
	// into the result slice
	matches := sfuzzy.Find(query, file_list)
	for _, match := range matches {
		items = append(items, file_list[match.Index])
	}

	result := unique(items)

	return result, nil
}

func (s *diskStore) GlobSearch(query string) ([]string, error) {
	// Search:
	// 	1. DOMAIN/USERNAME.gpg
	//	2. DOMAIN.gpg
	//	3. DOMAIN/SUBDIRECTORY/USERNAME.gpg

	matches, err := zglob.GlobFollowSymlinks(s.path + "/**/" + query + "*/**/*.gpg")
	if err != nil {
		return nil, err
	}

	matches2, err := zglob.GlobFollowSymlinks(s.path + "/**/" + query + "*.gpg")
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
	if strings.Count(query, ".") >= 2 {
		// try finding additional items by removing subparts of the query
		queryParts := strings.SplitN(query, ".", 2)[1:]
		newItems, err := s.GlobSearch(strings.Join(queryParts, "."))
		if err != nil {
			return nil, err
		}
		items = append(items, newItems...)
	}

	result := unique(items)
	sort.Strings(result)

	return result, nil
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

func unique(elems []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, elem := range elems {
		if !seen[elem] {
			seen[elem] = true
			result = append(result, elem)
		}
	}
	return result
}
