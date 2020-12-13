package file

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dkaslovsky/TextNote/pkg/template"
)

// Exists evaluates if a file for a template exists
func Exists(t *template.Template) bool {
	fileName := t.GetFilePath()
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

// Read reads a template from file
func Read(t *template.Template) error {
	r, err := os.Open(t.GetFilePath())
	if err != nil {
		return err
	}
	defer r.Close()
	return t.Load(r)
}

// Overwrite writes a template to a file, overwriting existing file contents if any
func Overwrite(t *template.Template) error {
	f, err := os.OpenFile(t.GetFilePath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Write(f)
	if err != nil {
		return err
	}
	return nil
}

// WriteIfNotExists writes a template to a file if the file does not already exist
func WriteIfNotExists(t *template.Template) error {
	if Exists(t) {
		return nil
	}
	return Overwrite(t)
}

// OpenInEditor opens a template in Vim
func OpenInEditor(t *template.Template) error {
	lineArg := fmt.Sprintf("+%d", t.GetFileStartLine())
	cmd := exec.Command("vim", lineArg, t.GetFilePath())
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
