package editor

import (
	"fmt"
	"os"
	"os/exec"
)

// EnvEditor is the name of the environment variable specifying the editor for opening notes
const EnvEditor = "EDITOR"

const (
	editorNameEmacs  = "emacs"
	editorNameNano   = "nano"
	editorNameNeovim = "nvim"
	editorNameVi     = "vi"
	editorNameVim    = "vim"
)

// openable is the interface that an editor opens
type openable interface {
	GetFilePath() string
	GetFileCursorLine() int
}

// Editor encapsulates the commands and args necessary to open an editor in a shell
type Editor struct {
	Cmd       string
	GetArgs   func(int) []string
	Supported bool
	Default   bool
}

// Open opens an object satisfying the openable interface in the editor
// NOTE: it is recommended to use Go >= v.1.15.7 due to call to exec.Command()
// See: https://blog.golang.org/path-security
func (e *Editor) Open(o openable) error {
	args := append(e.GetArgs(o.GetFileCursorLine()), o.GetFilePath())
	cmd := exec.Command(e.Cmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetEditor gets an Editor based on a provided name
func GetEditor(name string) *Editor {
	switch name {
	case editorNameVi, editorNameVim:
		return &Editor{
			Cmd: name,
			GetArgs: func(line int) []string {
				return []string{
					fmt.Sprintf("+%d", line),
				}
			},
			Supported: true,
			Default:   false,
		}
	case editorNameEmacs:
		return &Editor{
			Cmd: editorNameEmacs,
			GetArgs: func(line int) []string {
				return []string{
					fmt.Sprintf("+%d", line),
				}
			},
			Supported: true,
			Default:   false,
		}
	case editorNameNano:
		return &Editor{
			Cmd: editorNameNano,
			GetArgs: func(line int) []string {
				return []string{
					fmt.Sprintf("+%d", line),
				}
			},
			Supported: true,
			Default:   false,
		}
	case editorNameNeovim:
		return &Editor{
			Cmd: editorNameNeovim,
			GetArgs: func(line int) []string {
				return []string{
					fmt.Sprintf("+%d", line),
				}
			},
			Supported: true,
			Default:   false,
		}
	// use Vim as the default editor
	case "":
		return &Editor{
			Cmd: editorNameVim,
			GetArgs: func(line int) []string {
				return []string{
					fmt.Sprintf("+%d", line),
				}
			},
			Supported: true,
			Default:   true,
		}
	// unrecognized editor will be passed no arguments
	default:
		return &Editor{
			Cmd: name,
			GetArgs: func(line int) []string {
				return []string{}
			},
			Supported: false,
			Default:   false,
		}
	}
}
