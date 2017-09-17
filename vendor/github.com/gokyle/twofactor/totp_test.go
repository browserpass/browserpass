package twofactor

import (
	"crypto"
	"fmt"
	"testing"
)

var rfcTotpKey = []byte("12345678901234567890")
var rfcTotpStep uint64 = 30

var rfcTotpTests = []struct {
	Time uint64
	Code string
	T    uint64
	Algo crypto.Hash
}{
	{59, "94287082", 1, crypto.SHA1},
	//{59, "46119246", 1, crypto.SHA256},
	//{59, "90693936", 1, crypto.SHA512},
	{1111111109, "07081804", 37037036, crypto.SHA1},
	// {1111111109, "68084774", 37037036, crypto.SHA256},
	// {1111111109, "25091201", 37037036, crypto.SHA512},
	{1111111111, "14050471", 37037037, crypto.SHA1},
	// {1111111111, "67062674", 37037037, crypto.SHA256},
	// {1111111111, "99943326", 37037037, crypto.SHA512},
	{1234567890, "89005924", 41152263, crypto.SHA1},
	// {1234567890, "91819424", 41152263, crypto.SHA256},
	// {1234567890, "93441116", 41152263, crypto.SHA512},
	{2000000000, "69279037", 66666666, crypto.SHA1},
	// {2000000000, "90698825", 66666666, crypto.SHA256},
	// {2000000000, "38618901", 66666666, crypto.SHA512},
	{20000000000, "65353130", 666666666, crypto.SHA1},
	// {20000000000, "77737706", 666666666, crypto.SHA256},
	// {20000000000, "47863826", 666666666, crypto.SHA512},
}

func TestTotpRFC(t *testing.T) {
	for _, tc := range rfcTotpTests {
		otp := NewTOTP(rfcTotpKey, 0, 30, 8, tc.Algo)
		if otp.otpCounter(tc.Time) != tc.T {
			fmt.Println("twofactor: invalid T")
			fmt.Printf("\texpected: %d\n", tc.T)
			fmt.Printf("\t  actual: %d\n", otp.otpCounter(tc.Time))
			t.FailNow()
		}

		if code := otp.otp(otp.otpCounter(tc.Time)); code != tc.Code {
			fmt.Println("twofactor: invalid TOTP")
			fmt.Printf("\texpected: %s\n", tc.Code)
			fmt.Printf("\t  actual: %s\n", code)
			t.FailNow()
		}
	}
}
