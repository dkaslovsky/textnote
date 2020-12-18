package main

import (
	"log"

	"github.com/dkaslovsky/TextNote/cmd"
	"github.com/dkaslovsky/TextNote/pkg/config"
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
