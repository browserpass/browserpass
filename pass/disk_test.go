package pass

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestDefaultStorePath(t *testing.T) {
	var home, expectedCustom, expected, actual string

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
		t.Errorf("1: '%s' does not match '%s'", expected, actual)
	}

	// custom directory from $PASSWORD_STORE_DIR
	expected, err = filepath.Abs("browserpass-test")
	if err != nil {
		t.Error(err)
	}

	os.Mkdir(expectedCustom, os.ModePerm)
	os.Setenv("PASSWORD_STORE_DIR", expected)
	actual, err = defaultStorePath()
	if err != nil {
		t.Error(err)
	}
	if expected != actual {
		t.Errorf("2: '%s' does not match '%s'", expected, actual)
	}

	// clean-up
	os.Setenv("PASSWORD_STORE_DIR", "")
	os.Remove(expected)
}

func TestDiskStore_Search_nomatch(t *testing.T) {
	store := diskStore{stores: []StoreDefinition{StoreDefinition{Name: "default", Path: "test_store"}}, useFuzzy: false}

	domain := "this-most-definitely-does-not-exist"
	logins, err := store.Search(domain)
	if err != nil {
		t.Fatal(err)
	}
	if len(logins) > 0 {
		t.Errorf("%s yielded results, but it should not", domain)
	}
}

func TestDiskStoreSearch(t *testing.T) {
	store := diskStore{stores: []StoreDefinition{StoreDefinition{Name: "default", Path: "test_store"}}, useFuzzy: false}
	expectedResult := "default:abc.com"
	testDomains := []string{"abc.com", "test.abc.com", "testing.test.abc.com"}
	for _, domain := range testDomains {
		searchResults, err := store.Search(domain)
		if err != nil {
			t.Fatal(err)
		}
		// check if result contains abc.com
		found := false
		for _, searchResult := range searchResults {
			if searchResult == expectedResult {
				found = true
				break
			}
		}
		if found != true {
			t.Fatalf("Couldn't find %v in %v", expectedResult, searchResults)
		}
	}
}

func TestDiskStoreSearchNoDuplicatesWhenPatternMatchesDomainAndUsername(t *testing.T) {
	store := diskStore{stores: []StoreDefinition{StoreDefinition{Name: "default", Path: "test_store"}}, useFuzzy: false}
	searchResult, err := store.Search("xyz")
	if err != nil {
		t.Fatal(err)
	}
	if len(searchResult) != 1 {
		t.Fatalf("Found %v results instead of 1", len(searchResult))
	}
	expectedResult := "default:xyz.com/xyz_user"
	if searchResult[0] != expectedResult {
		t.Fatalf("Couldn't find %v, found %v instead", expectedResult, searchResult[0])
	}
}

func TestDiskStoreSearchFollowsSymlinkFiles(t *testing.T) {
	store := diskStore{stores: []StoreDefinition{StoreDefinition{Name: "default", Path: "test_store"}}, useFuzzy: false}
	searchResult, err := store.Search("def.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(searchResult) != 1 {
		t.Fatalf("Found %v results instead of 1", len(searchResult))
	}
	expectedResult := "default:def.com"
	if searchResult[0] != expectedResult {
		t.Fatalf("Couldn't find %v, found %v instead", expectedResult, searchResult[0])
	}
}

func TestDiskStoreSearchFollowsSymlinkDirectories(t *testing.T) {
	store := diskStore{stores: []StoreDefinition{StoreDefinition{Name: "default", Path: "test_store"}}, useFuzzy: false}
	searchResult, err := store.Search("amazon.co.uk")
	if err != nil {
		t.Fatal(err)
	}
	if len(searchResult) != 2 {
		t.Fatalf("Found %v results instead of 2", len(searchResult))
	}
	expectedResult := []string{"default:amazon.co.uk/user1", "default:amazon.co.uk/user2"}
	if searchResult[0] != expectedResult[0] || searchResult[1] != expectedResult[1] {
		t.Fatalf("Couldn't find %v, found %v instead", expectedResult, searchResult)
	}
}

func TestDiskStoreSearchSubDirectories(t *testing.T) {
	store := diskStore{stores: []StoreDefinition{StoreDefinition{Name: "default", Path: "test_store"}}, useFuzzy: false}
	searchTermsMatches := map[string][]string{
		"abc.org": []string{"default:abc.org/user3", "default:abc.org/wiki/user4", "default:abc.org/wiki/work/user5"},
		"wiki":    []string{"default:abc.org/wiki/user4", "default:abc.org/wiki/work/user5"},
		"work":    []string{"default:abc.org/wiki/work/user5"},
	}

	for term, expectedResult := range searchTermsMatches {
		searchResult, err := store.Search(term)
		if err != nil {
			t.Fatal(err)
		}
		if len(searchResult) != len(expectedResult) {
			t.Fatalf("For term %v found %v results (%v) instead of %v (%v)", term, len(searchResult), searchResult, len(expectedResult), expectedResult)
		}
		for i := 0; i < len(expectedResult); i++ {
			if searchResult[i] != expectedResult[i] {
				t.Fatalf("Couldn't find %v, found %v instead", expectedResult, searchResult)
			}
		}
	}
}

func TestDiskStorePartSearch(t *testing.T) {
	store := diskStore{stores: []StoreDefinition{StoreDefinition{Name: "default", Path: "test_store"}}, useFuzzy: false}
	searchResult, err := store.Search("ab")
	if err != nil {
		t.Fatal(err)
	}
	if len(searchResult) != 4 {
		t.Fatalf("Found %v results instead of 4", len(searchResult))
	}
	expectedResult := []string{"default:abc.com", "default:abc.org/user3", "default:abc.org/wiki/user4", "default:abc.org/wiki/work/user5"}
	for i := 0; i < len(expectedResult); i++ {
		if searchResult[i] != expectedResult[i] {
			t.Fatalf("Couldn't find %v, found %v instead", expectedResult, searchResult)
		}
	}
}

func TestFuzzySearch(t *testing.T) {
	store := diskStore{stores: []StoreDefinition{StoreDefinition{Name: "default", Path: "test_store"}}, useFuzzy: true}
	searchResult, err := store.Search("amaz2")

	if err != nil {
		t.Fatal(err)
	}
	if len(searchResult) != 2 {
		t.Fatalf("Result size was: %d expected 2", len(searchResult))
	}

	expectedResult := map[string]bool{
		"default:amazon.co.uk/user2": true,
		"default:amazon.com/user2":   true,
	}

	for _, res := range searchResult {
		if !expectedResult[res] {
			t.Fatalf("Result %s not expected!", res)
		}
	}
}

func TestFuzzySearchNoResult(t *testing.T) {
	store := diskStore{stores: []StoreDefinition{StoreDefinition{Name: "default", Path: "test_store"}}, useFuzzy: true}
	searchResult, err := store.Search("vvv")

	if err != nil {
		t.Fatal(err)
	}
	if len(searchResult) != 0 {
		t.Fatalf("Result size was: %d expected 0", len(searchResult))
	}
}

func TestFuzzySearchTopLevelEntries(t *testing.T) {
	store := diskStore{stores: []StoreDefinition{StoreDefinition{Name: "default", Path: "test_store"}}, useFuzzy: true}
	searchResult, err := store.Search("def")

	if err != nil {
		t.Fatal(err)
	}
	if len(searchResult) != 1 {
		t.Fatalf("Result size was: %d expected 1", len(searchResult))
	}

	expectedResult := map[string]bool{
		"default:def.com": true,
	}

	for _, res := range searchResult {
		if !expectedResult[res] {
			t.Fatalf("Result %s not expected!", res)
		}
	}
}

func TestGlobSearchMultipleStores(t *testing.T) {
	store := diskStore{stores: []StoreDefinition{StoreDefinition{Name: "default", Path: "test_store"}, StoreDefinition{Name: "custom", Path: "test_store_2"}}, useFuzzy: false}
	searchResults, err := store.Search("abc.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(searchResults) != 2 {
		t.Fatalf("Found %v results instead of 2", len(searchResults))
	}
	expectedResults := []string{"custom:abc.com", "default:abc.com"}
	if searchResults[0] != expectedResults[0] || searchResults[1] != expectedResults[1] {
		t.Fatalf("Couldn't find %v, found %v instead", expectedResults, searchResults)
	}
}

func TestFuzzySearchMultipleStores(t *testing.T) {
	store := diskStore{stores: []StoreDefinition{StoreDefinition{Name: "default", Path: "test_store"}, StoreDefinition{Name: "custom", Path: "test_store_2"}}, useFuzzy: true}
	searchResults, err := store.Search("abc.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(searchResults) != 2 {
		t.Fatalf("Found %v results instead of 2", len(searchResults))
	}
	expectedResults := []string{"default:abc.com", "custom:abc.com"}
	if searchResults[0] != expectedResults[0] || searchResults[1] != expectedResults[1] {
		t.Fatalf("Couldn't find %v, found %v instead", expectedResults, searchResults)
	}
}
