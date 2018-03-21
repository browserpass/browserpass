package browserpass

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dannyvankooten/browserpass/pass"
	"github.com/dannyvankooten/browserpass/protector"
	"github.com/gokyle/twofactor"
)

// Login represents a single pass login.
type Login struct {
	Username string `json:"u"`
	Password string `json:"p"`
	OTP      string `json:"digits"`
	OTPLabel string `json:"label"`
	URL      string `json:"url"`
}

var endianness = binary.LittleEndian

// Settings info for the browserpass program.
//
// The browser extension will look up settings in its localstorage and find
// which options have been selected by the user, and put them in a JSON object
// which is then passed along with the command over the native messaging api.
type Config struct {
	// Manual searches use FuzzySearch if true, GlobSearch otherwise
	UseFuzzy bool `json:"use_fuzzy_search"`
}

// msg defines a message sent from a browser extension.
type msg struct {
	Settings Config `json:"settings"`
	Action   string `json:"action"`
	Domain   string `json:"domain"`
	Entry    string `json:"entry"`
}

// Run starts browserpass.
func Run(stdin io.Reader, stdout io.Writer, s pass.Store) error {
	protector.Protect("stdio rpath proc exec")
	for {
		// Get message length, 4 bytes
		var n uint32
		if err := binary.Read(stdin, endianness, &n); err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		// Get message body
		var data msg
		lr := &io.LimitedReader{R: stdin, N: int64(n)}
		if err := json.NewDecoder(lr).Decode(&data); err != nil {
			return err
		}

		// Since the pass.Store object is created by the wrapper prior to
		// settings from the browser being made available, we set them here
		s.SetConfig(&data.Settings.UseFuzzy)

		var resp interface{}
		switch data.Action {
		case "search":
			list, err := s.Search(data.Domain)
			if err != nil {
				return err
			}
			resp = list
		case "match_domain":
			list, err := s.GlobSearch(data.Domain)
			if err != nil {
				return err
			}
			resp = list
		case "get":
			rc, err := s.Open(data.Entry)
			if err != nil {
				return err
			}
			defer rc.Close()
			login, err := readLoginGPG(rc)
			if err != nil {
				return err
			}
			if login.Username == "" {
				login.Username = guessUsername(data.Entry)
			}
			resp = login
		default:
			return errors.New("Invalid action")
		}

		var b bytes.Buffer
		if err := json.NewEncoder(&b).Encode(resp); err != nil {
			return err
		}

		if err := binary.Write(stdout, endianness, uint32(b.Len())); err != nil {
			return err
		}
		if _, err := b.WriteTo(stdout); err != nil {
			return err
		}
	}
}

func detectGPGBin() (string, error) {
	binPriorityList := []string{
		"gpg2", "/bin/gpg2", "/usr/bin/gpg2", "/usr/local/bin/gpg2",
		"gpg", "/bin/gpg", "/usr/bin/gpg", "/usr/local/bin/gpg",
	}

	binToUse := ""
	for _, bin := range binPriorityList {
		binCheck := exec.Command(bin, "--version")
		if err := binCheck.Run(); err == nil {
			binToUse = bin
			break
		}
	}

	if binToUse == "" {
		return "", errors.New("Unable to detect the location of gpg binary")
	}

	return binToUse, nil
}

// readLoginGPG reads a encrypted login from r using the system's GPG binary.
func readLoginGPG(r io.Reader) (*Login, error) {
	gpgbin, err := detectGPGBin()
	if err != nil {
		return nil, err
	}

	opts := []string{"--decrypt", "--yes", "--quiet", "--batch", "-"}

	// Run gpg
	cmd := exec.Command(gpgbin, opts...)

	cmd.Stdin = r

	rc, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	var errbuf bytes.Buffer
	cmd.Stderr = &errbuf

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	protector.Protect("stdio")

	// Read decrypted output
	login, err := parseLogin(rc)
	if err != nil {
		return nil, err
	}

	defer rc.Close()

	if err := cmd.Wait(); err != nil {
		return nil, errors.New(err.Error() + "\n" + errbuf.String())
	}
	return login, nil
}

func parseTotp(str string, l *Login) error {
	urlPattern := regexp.MustCompile("^otpauth.*$")
	ourl := urlPattern.FindString(str)

	if ourl == "" {
		tokenPattern := regexp.MustCompile("(?i)^totp(-secret)?:")
		token := tokenPattern.ReplaceAllString(str, "");
		if len(token) != len(str) {
			ourl = "otpauth://totp/?secret=" + strings.TrimSpace(token)
		}
	}
	if ourl != "" {
		o, label, err := twofactor.FromURL(ourl)
		if err != nil {
			return err
		}
		l.OTP = o.OTP()
		l.OTPLabel = label
	}

	return nil
}

// parseLogin parses a login and a password from a decrypted password file.
func parseLogin(r io.Reader) (*Login, error) {
	login := new(Login)

	scanner := bufio.NewScanner(r)

	// The first line is the password
	scanner.Scan()
	login.Password = scanner.Text()

	// Keep reading file for string in "login:", "username:" or "user:" format (case insensitive).
	userPattern := regexp.MustCompile("(?i)^(login|username|user):")
	urlPattern := regexp.MustCompile("(?i)^(url|link|website|web|site):")
	for scanner.Scan() {
		line := scanner.Text()
		parseTotp(line, login)
		replaced := userPattern.ReplaceAllString(line, "")
		if len(replaced) != len(line) {
			login.Username = strings.TrimSpace(replaced)
		}
		if (login.URL == "") {
			replaced = urlPattern.ReplaceAllString(line, "")
			if len(replaced) != len(line) {
				login.URL = strings.TrimSpace(replaced)
			}
		}
	}

	// if an unlabelled OTP is present, label it with the username
	if strings.TrimSpace(login.OTPLabel) == "" && login.OTP != "" {
		login.OTPLabel = login.Username
	}

	return login, nil
}

// guessLogin tries to guess a username from an entry's name.
func guessUsername(name string) string {
	if strings.Count(filepath.ToSlash(name), "/") >= 1 {
		return filepath.Base(name)
	}
	return ""
}
