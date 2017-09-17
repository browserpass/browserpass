package twofactor

import (
	"fmt"
	"testing"
)

var sha1Hmac = []byte{
	0x1f, 0x86, 0x98, 0x69, 0x0e,
	0x02, 0xca, 0x16, 0x61, 0x85,
	0x50, 0xef, 0x7f, 0x19, 0xda,
	0x8e, 0x94, 0x5b, 0x55, 0x5a,
}

var truncExpect int64 = 0x50ef7f19

// This test runs through the truncation example given in the RFC.
func TestTruncate(t *testing.T) {
	if result := truncate(sha1Hmac); result != truncExpect {
		fmt.Printf("hotp: expected truncate -> %d, saw %d\n",
			truncExpect, result)
		t.FailNow()
	}

	sha1Hmac[19]++
	if result := truncate(sha1Hmac); result == truncExpect {
		fmt.Println("hotp: expected truncation to fail")
		t.FailNow()
	}
}
