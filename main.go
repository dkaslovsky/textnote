package main

import (
	"log"

	"github.com/dkaslovsky/textnote/cmd"
	"github.com/dkaslovsky/textnote/pkg/config"
)

const (
	appName    = "textnote"
	appVersion = "undefined"
)

func main() {
	err := config.EnsureAppDir()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Run(appName, appVersion)
	if err != nil {
		log.Fatal(err)
	}
}
