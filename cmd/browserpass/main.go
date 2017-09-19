package main

import (
	"io"
	"log"
	"os"

	"github.com/dannyvankooten/browserpass"
	"github.com/dannyvankooten/browserpass/pass"
	"github.com/dannyvankooten/browserpass/protector"
)

func main() {
	protector.Protect("stdio rpath proc exec getpw")
	log.SetPrefix("[Browserpass] ")

	s, err := pass.NewDefaultStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := browserpass.Run(os.Stdin, os.Stdout, s); err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
