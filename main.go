package main

import (
	"log"

	"github.com/dkaslovsky/TextNote/cmd"
)

func main() {
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
