package pass

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mattn/go-zglob"
	"github.com/mattn/go-zglob/fastwalk"
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

// lookup will find sites matching or partly matching the passed domain.
func lookup(domain string, sites []*site) []*site {
	results := make([]*site, 0)
	domainParts := reverse(strings.Split(domain, "."))
	for _, s := range sites {
		parts := reverse(strings.Split(s.domain, "."))
		if subMatch(domainParts, parts, 2) {
			results = append(results, s)
		}
	}

	// sort by length of domain, longest first.
	sort.Slice(results, func(i, j int) bool {
		return len(results[i].domain) > len(results[j].domain)
	})

	return results
}

// reverse reverses the element order of a string slice.
func reverse(strings []string) []string {
	for i, j := 0, len(strings)-1; i < j; i, j = i+1, j-1 {
		strings[i], strings[j] = strings[j], strings[i]
	}

	return strings
}

// subMatch validates whether the candidate string slice matches at least min
// parts of the query slice.
// Example:
//    query = ["org", "example", "my"]
//    candidate = ["org", "example"]
//    min = 2
// this will match because both elements of candidate is found in the same
// order in query.
func subMatch(query, candidate []string, min int) bool {
	if len(candidate) < min {
		return false
	}

	if len(query) < len(candidate) {
		return false
	}

	matches := 0
	for i := len(candidate) - 1; i > -1; i-- {
		if candidate[i] == query[i] {
			matches++
		} else {
			return false
		}

		if matches >= min {
			return true
		}
	}

	return false
}

// site defines a domain and the related users stored for this domain.
type site struct {
	domain string
	users  []string
}

// Lookup looks up domains in the password store based on the domainQuery string.
// the lookup will return a list of domains/subdomains matching the query with
// the most precise domain first.
// Assuming the domains "sub1.domain.tld", "sub2.domain.tld" and "domain.tld"
// are defined in the store, and the query domain is "sub1.domain.tld" the
// lookup will return the list ["sub1.domain.tld", "domain.tld"]. If
// "domain.tld" is the query domain then only ["domain.tld"] will be returned
// as the matching is done from front to end thus not matching the subdoamins.
func (s *diskStore) Lookup(domainQuery string) ([]string, error) {
	sites := make([]*site, 0)
	siteCh := make(chan *site)
	errCh := make(chan error)

	go func() {
		// use FastWalk to collect all domains/users defined in the password store.
		err := fastwalk.FastWalk(s.path, func(dir string, typ os.FileMode) error {
			if dir == s.path {
				return nil
			}
			if typ&os.ModeDir != 0 {
				files, err := ioutil.ReadDir(dir)
				if err != nil {
					return err
				}

				users := make([]string, 0, len(files))
				for _, file := range files {
					users = append(users, strings.TrimSuffix(file.Name(), ".gpg"))
				}
				siteCh <- &site{domain: path.Base(dir), users: users}
				return filepath.SkipDir
			}
			return nil
		})
		close(siteCh)
		errCh <- err
	}()

	for site := range siteCh {
		sites = append(sites, site)
	}

	err := <-errCh
	if err != nil {
		return nil, err
	}

	sites = lookup(domainQuery, sites)

	results := make([]string, 0, len(sites))
	for _, site := range sites {
		for _, user := range site.users {
			results = append(results, site.domain+"/"+user)
		}
	}

	return results, nil
}

func (s *diskStore) Search(query string) ([]string, error) {
	// First, search for DOMAIN/USERNAME.gpg
	// Then, search for DOMAIN.gpg
	matches, err := zglob.Glob(s.path + "/**/" + query + "*/*.gpg")
	if err != nil {
		return nil, err
	}

	matches2, err := zglob.Glob(s.path + "/**/" + query + "*.gpg")
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
