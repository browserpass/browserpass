package twofactor

import "fmt"
import "io"
import "testing"

func TestHOTPString(t *testing.T) {
	hotp := NewHOTP(nil, 0, 6)
	hotpString := otpString(hotp)
	if hotpString != "OATH-HOTP, 6" {
		fmt.Println("twofactor: invalid OTP string")
		t.FailNow()
	}
}

// This test generates a new OTP, outputs the URL for that OTP,
// and attempts to parse that URL. It verifies that the two OTPs
// are the same, and that they produce the same output.
func TestURL(t *testing.T) {
	var ident = "testuser@foo"
	otp := NewHOTP(testKey, 0, 6)
	url := otp.URL("testuser@foo")
	otp2, id, err := FromURL(url)
	if err != nil {
		fmt.Printf("hotp: failed to parse HOTP URL\n")
		t.FailNow()
	} else if id != ident {
		fmt.Printf("hotp: bad label\n")
		fmt.Printf("\texpected: %s\n", ident)
		fmt.Printf("\t  actual: %s\n", id)
		t.FailNow()
	} else if otp2.Counter() != otp.Counter() {
		fmt.Printf("hotp: OTP counters aren't synced\n")
		fmt.Printf("\toriginal: %d\n", otp.Counter())
		fmt.Printf("\t  second: %d\n", otp2.Counter())
		t.FailNow()
	}

	code1 := otp.OTP()
	code2 := otp2.OTP()
	if code1 != code2 {
		fmt.Printf("hotp: mismatched OTPs\n")
		fmt.Printf("\texpected: %s\n", code1)
		fmt.Printf("\t  actual: %s\n", code2)
	}

	// There's not much we can do test the QR code, except to
	// ensure it doesn't fail.
	_, err = otp.QR(ident)
	if err != nil {
		fmt.Printf("hotp: failed to generate QR code PNG (%v)\n", err)
		t.FailNow()
	}

	// This should fail because the maximum size of an alphanumeric
	// QR code with the lowest-level of error correction should
	// max out at 4296 bytes. 8k may be a bit overkill... but it
	// gets the job done. The value is read from the PRNG to
	// increase the likelihood that the returned data is
	// uncompressible.
	var tooBigIdent = make([]byte, 8192)
	_, err = io.ReadFull(PRNG, tooBigIdent)
	if err != nil {
		fmt.Printf("hotp: failed to read identity (%v)\n", err)
		t.FailNow()
	} else if _, err = otp.QR(string(tooBigIdent)); err == nil {
		fmt.Println("hotp: QR code should fail to encode oversized URL")
		t.FailNow()
	}
}

// This test attempts a variety of invalid urls against the parser
// to ensure they fail.
func TestBadURL(t *testing.T) {
	var urlList = []string{
		"http://google.com",
		"",
		"-",
		"foo",
		"otpauth:/foo/bar/baz",
		"://",
		"otpauth://hotp/secret=bar",
		"otpauth://hotp/?secret=QUJDRA&algorithm=SHA256",
		"otpauth://hotp/?digits=",
		"otpauth://hotp/?secret=123",
		"otpauth://hotp/?secret=MFRGGZDF&digits=ABCD",
		"otpauth://hotp/?secret=MFRGGZDF&counter=ABCD",
	}

	for i := range urlList {
		if _, _, err := FromURL(urlList[i]); err == nil {
			fmt.Println("hotp: URL should not have parsed successfully")
			fmt.Printf("\turl was: %s\n", urlList[i])
			t.FailNow()
		}
	}
}
