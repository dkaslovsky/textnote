package file

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// ReadWriteable is the interface for which file operations are executed
type ReadWriteable interface {
	Load(io.Reader) error
	Write(io.Writer) error
	GetFilePath() string
	GetFileStartLine() int
}

// Exists evaluates if a file exists
func Exists(rw ReadWriteable) bool {
	fileName := rw.GetFilePath()
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

// Read reads from file
func Read(rw ReadWriteable) error {
	r, err := os.Open(rw.GetFilePath())
	if err != nil {
		return err
	}
	defer r.Close()
	return rw.Load(r)
}

// Overwrite writes a template to a file, overwriting existing file contents if any
func Overwrite(rw ReadWriteable) error {
	f, err := os.OpenFile(rw.GetFilePath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = rw.Write(f)
	if err != nil {
		return err
	}
	return nil
}

// WriteIfNotExists writes a template to a file if the file does not already exist
func WriteIfNotExists(rw ReadWriteable) error {
	if Exists(rw) {
		return nil
	}
	return Overwrite(rw)
}

// OpenInEditor opens a template in Vim
func OpenInEditor(rw ReadWriteable) error {
	lineArg := fmt.Sprintf("+%d", rw.GetFileStartLine())
	cmd := exec.Command("vim", lineArg, rw.GetFilePath())
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
