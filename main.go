package main

import (
	"log"

	"github.com/dkaslovsky/textnote/cmd"
	"github.com/dkaslovsky/textnote/pkg/config"
)

func main() {
	err := config.EnsureAppDir()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
