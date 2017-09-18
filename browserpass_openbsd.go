// +build openbsd

package browserpass

import (
	"golang.org/x/sys/unix"
)

func Protector(s string) {
	unix.Pledge(s, nil)
}
