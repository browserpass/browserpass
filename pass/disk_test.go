package pass

import (
	"os"
	"testing"
)

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

func TestDiskStore_Search(t *testing.T) {
	storedir := "/tmp/browserpass-test/"

	foo_u1_path := storedir + "foo.com/u1.gpg"
	foo_u2_path := storedir + "foo.com/u2.gpg"
	os.MkdirAll(storedir + "foo.com", os.ModePerm)
	os.OpenFile(foo_u1_path, os.O_RDONLY|os.O_CREATE, 0666)
	os.OpenFile(foo_u2_path, os.O_RDONLY|os.O_CREATE, 0666)

	a_foo_path := storedir + "a.foo.com.gpg"
	os.OpenFile(a_foo_path, os.O_RDONLY|os.O_CREATE, 0666)

	os.Setenv("PASSWORD_STORE_DIR", storedir)
	s, err := NewDefaultStore()
	if err != nil {
		t.Fatal(err)
	}

	domain := "b.foo.com"
	logins, err := s.Search(domain)
	if err != nil {
		t.Fatal(err)
	}
	if len(logins) != 2 {
		t.Fatal("%s yielded %d results, expected 2 results", domain, len(logins))
	}
	if logins[0] != "foo.com/u1" && logins[0] != "foo.com/u2" {
		t.Fatal("%s yielded %s, which isn't in expected results: foo.com/{u1,u2}",
			domain, logins[0])
	}
	if logins[1] != "foo.com/u1" && logins[1] != "foo.com/u2" {
		t.Fatal("%s yielded %s, which isn't in expected results: foo.com/{u1,u2}",
			domain, logins[1])
	}

	domain = "x.y.a.foo.com"
	logins, err = s.Search(domain)
	if err != nil {
		t.Fatal(err)
	}
	if len(logins) != 1 {
		t.Fatal("%s yielded %d results, expected 1 result", domain, len(logins))
	}
	if logins[0] != "a.foo.com" {
		t.Fatal("%s yielded %s, expected a.foo.com", domain, logins[0])
	}

	domain = "bar.com"
	logins, err = s.Search(domain)
	if err != nil {
		t.Fatal(err)
	}
	if len(logins) != 0 {
		t.Fatal("%s yielded %d results, expected no results", domain, len(logins))
	}

	// clean-up
	os.Setenv("PASSWORD_STORE_DIR", "")
	os.RemoveAll(storedir)
}
