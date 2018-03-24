package main

import (
	"io"
	"log"
	"os"

	"github.com/dannyvankooten/browserpass"
	"github.com/dannyvankooten/browserpass/protector"
)

func main() {
	protector.Protect("stdio rpath proc exec getpw")
	log.SetPrefix("[Browserpass] ")

	if err := browserpass.Run(os.Stdin, os.Stdout); err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
