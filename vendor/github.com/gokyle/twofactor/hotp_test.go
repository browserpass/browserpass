package twofactor

import (
	"fmt"
	"testing"
)

var testKey = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

func newZeroHOTP() *HOTP {
	return NewHOTP(testKey, 0, 6)
}

var rfcHotpKey = []byte("12345678901234567890")
var rfcHotpExpected = []string{
	"755224",
	"287082",
	"359152",
	"969429",
	"338314",
	"254676",
	"287922",
	"162583",
	"399871",
	"520489",
}

// This test runs through the test cases presented in the RFC, and
// ensures that this implementation is in compliance.
func TestHotpRFC(t *testing.T) {
	otp := NewHOTP(rfcHotpKey, 0, 6)
	for i := 0; i < len(rfcHotpExpected); i++ {
		if otp.Counter() != uint64(i) {
			fmt.Printf("twofactor: invalid counter (should be %d, is %d",
				i, otp.Counter())
			t.FailNow()
		}
		code := otp.OTP()
		if code == "" {
			fmt.Printf("twofactor: failed to produce an OTP\n")
			t.FailNow()
		} else if code != rfcHotpExpected[i] {
			fmt.Printf("twofactor: invalid OTP\n")
			fmt.Printf("\tExpected: %s\n", rfcHotpExpected[i])
			fmt.Printf("\t  Actual: %s\n", code)
			fmt.Printf("\t Counter: %d\n", otp.counter)
			t.FailNow()
		}
	}
}

// This test uses a different key than the test cases in the RFC,
// but runs through the same test cases to ensure that they fail as
// expected.
func TestHotpBadRFC(t *testing.T) {
	otp := NewHOTP(testKey, 0, 6)
	for i := 0; i < len(rfcHotpExpected); i++ {
		code := otp.OTP()
		if code == "" {
			fmt.Printf("twofactor: failed to produce an OTP\n")
			t.FailNow()
		} else if code == rfcHotpExpected[i] {
			fmt.Printf("twofactor: should not have received a valid OTP\n")
			t.FailNow()
		}
	}
}
