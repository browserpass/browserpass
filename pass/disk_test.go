package pass

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"
)

type domainUser struct {
	domain, user string
}

func setupTestStore(t *testing.T) string {
	storeDir := path.Join(os.TempDir(), fmt.Sprintf("browserpass-%d", time.Now().UTC().UnixNano()))

	domains := []*domainUser{
		{"foo.bar", "u1.gpg"},
		{"foo.bar", "u2.gpg"},
		{"baz.foo.bar", "u1.gpg"},
		{"domain.tld", "u1.gpg"},
		{"sub.domain.tld", "u1.gpg"},
		{"sub.sub.domain.tld", "u1.gpg"},
	}

	for _, du := range domains {
		err := os.MkdirAll(path.Join(storeDir, du.domain), os.ModePerm)
		if err != nil {
			t.Errorf("should not fail: %s", err)
		}

		f, err := os.Create(path.Join(storeDir, du.domain, du.user))
		if err != nil {
			t.Errorf("should not fail: %s", err)
		}
		err = f.Close()
		if err != nil {
			t.Errorf("should not fail: %s", err)
		}
	}

	return storeDir
}

func cleanTestStore(t *testing.T, storeDir string) {
	err := os.RemoveAll(storeDir)
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}
}

func TestDefaultStorePath(t *testing.T) {
	var home, expected, actual string
	home = os.Getenv("HOME")

	// default directory
	os.Setenv("PASSWORD_STORE_DIR", "")
	expected = home + "/.password-store"
	actual, _ = defaultStorePath()
	if expected != actual {
		t.Errorf("%s does not match %s", expected, actual)
	}

	// custom directory from $PASSWORD_STORE_DIR
	expected = "/tmp/browserpass-test"
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
	storeDir := setupTestStore(t)
	defer cleanTestStore(t, storeDir)

	os.Setenv("PASSWORD_STORE_DIR", storeDir)
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

func TestDiskStoreLookup(t *testing.T) {
	storeDir := setupTestStore(t)
	defer cleanTestStore(t, storeDir)

	// TODO fix
	os.Setenv("PASSWORD_STORE_DIR", storeDir)

	s, err := NewDefaultStore()
	if err != nil {
		t.Error(err)
	}

	for _, tc := range []struct {
		msg      string
		domain   string
		expected []string
	}{
		{
			msg:      "should find sub.domain.tld followed by domain.tld",
			domain:   "sub.domain.tld",
			expected: []string{"sub.domain.tld/u1", "domain.tld/u1"},
		},
		{
			msg:      "should find sub.sub.domain.tld, sub.domain.tld followed by domain.tld",
			domain:   "sub.sub.domain.tld",
			expected: []string{"sub.sub.domain.tld/u1", "sub.domain.tld/u1", "domain.tld/u1"},
		},
		{
			msg:      "should find two users for foo.bar",
			domain:   "foo.bar",
			expected: []string{"foo.bar/u1", "foo.bar/u2"},
		},
	} {
		t.Run(tc.msg, func(t *testing.T) {
			logins, err := s.Lookup(tc.domain)
			if err != nil {
				t.Error(err)
			}

			if len(logins) != len(tc.expected) {
				t.Errorf("expected %d results, got %d", len(tc.expected), len(logins))
			}

			for i, l := range logins {
				if l != tc.expected[i] {
					t.Errorf("expected %s, got %s", tc.expected[i], l)
				}
			}
		})
	}
}

func TestLookup(t *testing.T) {
	sites := []*site{
		{domain: "example.org", users: nil},
		{domain: "my.example.org", users: nil},
		{domain: "other.example.org", users: nil},
		{domain: "google.com", users: nil},
	}

	domain := "my.example.org"

	sitesFound := lookup(domain, sites)
	if len(sitesFound) != 2 {
		t.Fatalf("expected 2 sites matching, found %d", len(sitesFound))
	}

	if sitesFound[0].domain != domain {
		t.Fatalf("expected first domain to be '%s', got '%s'", domain, sitesFound[0].domain)
	}
}
