package pass

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestDefaultStorePath(t *testing.T) {
	var home, expected, actual string

	usr, err := user.Current()

	if err != nil {
		t.Error(err)
	}

	home = usr.HomeDir

	// default directory
	os.Setenv("PASSWORD_STORE_DIR", "")
	expected = filepath.Join(home, ".password-store")
	actual, _ = defaultStorePath()
	if expected != actual {
		t.Errorf("%s does not match %s", expected, actual)
	}

	// custom directory from $PASSWORD_STORE_DIR
	expected, err = filepath.Abs("browserpass-test")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(expected)
	os.Mkdir(expected, os.ModePerm)
	os.Setenv("PASSWORD_STORE_DIR", expected)
	actual, _ = defaultStorePath()
	if expected != actual {
		t.Errorf("%s does not match %s", expected, actual)
	}

	// clean-up
	os.Setenv("PASSWORD_STORE_DIR", "")
	os.Remove(expected)
}

func TestDiskStore_Search_nomatch(t *testing.T) {
	s, err := NewDefaultStore()
	if err != nil {
		t.Fatal(err)
	}

	domain := "this-most-definitely-does-not-exist"
	logins, err := s.Search(domain)
	if err != nil {
		t.Fatal(err)
	}
	if len(logins) > 0 {
		t.Errorf("%s yielded results, but it should not", domain)
	}
}

func TestDiskStoreSearch(t *testing.T) {
	store := diskStore{"test_store"}
	targetDomain := "abc.com"
	testDomains := []string{"abc.com", "test.abc.com", "testing.test.abc.com"}
	for _, domain := range testDomains {
		searchResults, err := store.Search(domain)
		if err != nil {
			t.Fatal(err)
		}
		// check if result contains abc.com
		found := false
		for _, searchResult := range searchResults {
			if searchResult == targetDomain {
				found = true
				break
			}
		}
		if found != true {
			t.Fatalf("Couldn't find %v in %v", targetDomain, searchResults)
		}
	}
}
