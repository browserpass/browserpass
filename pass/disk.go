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

// StoreDefinition defines a password store object
type StoreDefinition struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type diskStore struct {
	stores   []StoreDefinition
	useFuzzy bool // Setting for FuzzySearch or GlobSearch in manual searches
}

func NewDefaultStore(stores []StoreDefinition, useFuzzy bool) (Store, error) {
	if stores == nil || len(stores) == 0 {
		defaultPath, err := defaultStorePath()
		if err != nil {
			return nil, err
		}
		stores = []StoreDefinition{{Name: "default", Path: defaultPath}}
	}

	// Expand paths, follow symlinks
	for i, store := range stores {
		path := store.Path
		if strings.HasPrefix(path, "~/") {
			path = filepath.Join("$HOME", path[2:])
		}
		path = os.ExpandEnv(path)
		path, err := filepath.EvalSymlinks(path)
		if err != nil {
			return nil, err
		}
		stores[i].Path = path
	}

	return &diskStore{stores, useFuzzy}, nil
}

func defaultStorePath() (string, error) {
	path := os.Getenv("PASSWORD_STORE_DIR")
	if path != "" {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	path = filepath.Join(usr.HomeDir, ".password-store")
	return path, nil
}

// Do a search. Will call into the correct algoritm (glob or fuzzy)
// based on the settings present in the diskStore struct
func (s *diskStore) Search(query string) ([]string, error) {
	if s.useFuzzy {
		return s.FuzzySearch(query)
	}
	return s.GlobSearch(query)
}

// Fuzzy searches first get a list of all pass entries by doing a glob search
// for the empty string, then apply appropriate logic to convert results to
// a slice of strings, finally returning all of the unique entries.
func (s *diskStore) FuzzySearch(query string) ([]string, error) {
	entries, err := s.GlobSearch("")
	if err != nil {
		return nil, err
	}

	// GlobSearch now results `storename:filename`, for fuzzy search we need to provide only file names
	var fileNames []string
	for _, entry := range entries {
		fileNames = append(fileNames, strings.SplitN(entry, ":", 2)[1])
	}

	// The resulting match struct does not copy the strings, but rather
	// provides the index to the original array. Copy those strings
	// into the result slice
	var results []string
	matches := sfuzzy.Find(query, fileNames)
	for _, match := range matches {
		results = append(results, entries[match.Index])
	}

	return results, nil
}

func (s *diskStore) GlobSearch(query string) ([]string, error) {
	// Search:
	//	1. DOMAIN/USERNAME.gpg
	//	2. DOMAIN.gpg
	//	3. DOMAIN/SUBDIRECTORY/USERNAME.gpg

	items := []string{}

	for _, store := range s.stores {
		matches, err := zglob.GlobFollowSymlinks(store.Path + "/**/" + query + "*/**/*.gpg")
		if err != nil {
			return nil, err
		}

		matches2, err := zglob.GlobFollowSymlinks(store.Path + "/**/" + query + "*.gpg")
		if err != nil {
			return nil, err
		}

		allMatches := append(matches, matches2...)

		for i, match := range allMatches {
			item, err := filepath.Rel(store.Path, match)
			if err != nil {
				return nil, err
			}
			allMatches[i] = store.Name + ":" + strings.TrimSuffix(item, ".gpg")
		}

		items = append(items, allMatches...)
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
	parts := strings.SplitN(item, ":", 2)

	for _, store := range s.stores {
		if store.Name != parts[0] {
			continue
		}
		path := filepath.Join(store.Path, parts[1]+".gpg")
		f, err := os.Open(path)
		if os.IsNotExist(err) {
			continue
		}
		return f, err
	}
	return nil, errors.New("Unable to find the item on disk")
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
