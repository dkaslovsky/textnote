package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dkaslovsky/TextNote/cmd"
)

const envAppDir = "TEXTNOTE_DIR"

func main() {
	err := ensureAppDir()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func ensureAppDir() error {
	appDir := os.Getenv(envAppDir)
	if appDir == "" {
		return fmt.Errorf("required environment variable [%s] is not set", envAppDir)
	}
	finfo, err := os.Stat(appDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(appDir, 0755)
		if err != nil {
			return err
		}
		log.Printf("created directory [%s]", appDir)
		return nil
	}
	if !finfo.IsDir() {
		return fmt.Errorf("[%s=%s] must be a directory", envAppDir, appDir)
	}
	return nil
}
