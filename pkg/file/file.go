package file

import (
	"io"
	"os"
	"os/exec"
)

// ReadWriteable is the interface on which file operations are executed
type ReadWriteable interface {
	Load(io.Reader) error
	Write(io.Writer) error
	GetFilePath() string
}

// ReadWriter executes file operations
type ReadWriter struct{}

// NewReadWriter constructs a new ReadWriter
func NewReadWriter() *ReadWriter {
	return &ReadWriter{}
}

// Read reads from file
func (rw *ReadWriter) Read(rwable ReadWriteable) error {
	r, err := os.Open(rwable.GetFilePath())
	if err != nil {
		return err
	}
	defer r.Close()
	return rwable.Load(r)
}

// Overwrite writes a template to a file, overwriting existing file contents if any
func (rw *ReadWriter) Overwrite(rwable ReadWriteable) error {
	f, err := os.OpenFile(rwable.GetFilePath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = rwable.Write(f)
	if err != nil {
		return err
	}
	return nil
}

// Exists evaluates if a file exists
func (rw *ReadWriter) Exists(rwable ReadWriteable) bool {
	fileName := rwable.GetFilePath()
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

// Openable is the interface for which a file is opened
type Openable interface {
	GetFilePath() string
	GetFileCursorLine() int
}

// Editor is the interface for which an editor is opened
type Editor interface {
	GetCmd() string
	GetArgsFunc() func(int) []string
}

// Open opens a template in an editor
// NOTE: it is recommended to use Go >= v.1.15.7 due to call to exec.Command()
// See: https://blog.golang.org/path-security
func Open(o Openable, ed Editor) error {
	edArgs := append(ed.GetArgsFunc()(o.GetFileCursorLine()), o.GetFilePath())

	cmd := exec.Command(ed.GetCmd(), edArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
