package main

import (
	"log"

	"github.com/dkaslovsky/textnote/cmd"
	"github.com/dkaslovsky/textnote/pkg/config"
)

const name = "textnote"

var version string // set by build ldflags

func main() {
	err := config.EnsureAppDir()
	if err != nil {
		log.Fatal(err)
	}

	err = config.CreateIfNotExists()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Run(name, version)
	if err != nil {
		log.Fatal(err)
	}
}
