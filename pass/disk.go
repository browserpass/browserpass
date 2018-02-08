package pass

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"os/exec"
	"os/user"

	"github.com/mattn/go-zglob"
	//"github.com/kballard/go-shellquote"
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
	external_cmd := os.Getenv("BROWSERPASS_SEARCH_COMMAND")
	filter_cmd := os.Getenv("BROWSERPASS_FILTER_COMMAND")

	var res []string
	var err error

	if filter_cmd != "" {
		res, err = s.filterSearch(query, filter_cmd)
	} else if external_cmd == "" {
		res, err = s.builtinSearch(query)
	} else {
		res, err = s.externalSearch(query, external_cmd)
	}
	if err != nil {
		return nil, err
	}

	result := unique(res)
	sort.Strings(result)

	return result, nil
}

func (s *diskStore) filterSearch(query, cmd_str string) ([]string, error) {
	var result []string

	items, err := s.builtinSearch("")
	if err != nil {
		return nil, err
	}
	// Commands like fzf and grep will have an exit code of 1, causing an error
	// include single quotes to avoid shell execution issues
	cmd_str += " '" + query + "' || true"
	cmd := exec.Command("/bin/sh", "-c", cmd_str)

	ri, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	rc, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(rc)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	for _, item := range items {
		str := item + "\n"
		ri.Write([]byte(str))
	}

	ri.Close()

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			result = append(result, line)
		}
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return result, nil

}
func (s *diskStore) externalSearch(query, cmd_str string) ([]string, error) {
	var result []string

	// Commands like fzf and grep will have an exit code of 1, causing an error
	// include single quotes to avoid shell execution issues
	cmd_str += " '" + query + "' || true"
	cmd := exec.Command("/bin/sh", "-c", cmd_str)

	rc, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(rc)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			result = append(result, line)
		}
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *diskStore) builtinSearch(query string) ([]string, error) {
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
		newItems, err := s.Search(strings.Join(queryParts, "."))
		if err != nil {
			return nil, err
		}
		items = append(items, newItems...)
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
