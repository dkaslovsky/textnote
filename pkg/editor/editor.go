package editor

import "fmt"

// EnvEditor is the name of the environment variable specifying the editor for opening notes
const EnvEditor = "EDITOR"

const (
	editorNameEmacs  = "emacs"
	editorNameNano   = "nano"
	editorNameNeovim = "nvim"
	editorNameVi     = "vi"
	editorNameVim    = "vim"
)

// Editor encapsulates the commands and args necessary to open an editor in a shell
type Editor struct {
	Cmd       string
	GetArgs   func(int) []string
	Supported bool
	Default   bool
}

// GetCmd returns the editor's command
func (e *Editor) GetCmd() string {
	return e.Cmd
}

// GetArgsFunc returns the editor's GetArgs function
func (e *Editor) GetArgsFunc() func(int) []string {
	return e.GetArgs
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
