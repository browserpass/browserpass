package main

import (
	"io"
	"log"
	"os"

	"github.com/dannyvankooten/browserpass"
	"github.com/dannyvankooten/browserpass/pass"
)

func main() {
	log.SetPrefix("[Browserpass] ")

	s, err := pass.NewDefaultStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := browserpass.Run(os.Stdin, os.Stdout, s); err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
