package browserpass

import (
	"strings"
	"testing"

	"github.com/gokyle/twofactor"
)

func TestParseLogin(t *testing.T) {
	r := strings.NewReader("password\n\nfoo\nlogin: bar")

	login, err := parseLogin(r)
	if err != nil {
		t.Fatal(err)
	}

	if login.Password != "password" {
		t.Errorf("Password is %s, expected %s", login.Password, "password")
	}
	if login.Username != "bar" {
		t.Errorf("Username is %s, expected %s", login.Username, "bar")
	}
}

func TestOtp(t *testing.T) {
	r := strings.NewReader("password\n\nfoo\nlogin: bar\notpauth://totp/totp-secret?secret=GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ&issuer=test")

	login, err := parseLogin(r)
	if err != nil {
		t.Fatal(err)
	}

	if login.OTPLabel != "totp-secret" {
		t.Errorf("OTPLabel is '%s', expected '%s'", login.OTPLabel, "totp-secret")
	}

	o, err := twofactor.NewGoogleTOTP("GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ")
	if err != nil {
		t.Error(err)
	}

	otp := o.OTP()

	if login.OTP != otp {
		t.Errorf("OTP is %s, expected %s", login.OTP, otp)
	}
}

func TestGuessUsername(t *testing.T) {
	tests := map[string]string{
		"foo":     "",
		"foo/bar": "bar",
	}

	for input, expected := range tests {
		if username := guessUsername(input); username != expected {
			t.Errorf("guessUsername(%s): expected %s, got %s", input, expected, username)
		}
	}
}
