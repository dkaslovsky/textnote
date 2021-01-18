package open

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CommandOptions are the standard options for the package's commands
type CommandOptions struct {
	Copy   bool
	Delete bool
}

// AttachOpts attaches parsed CommandOptions to a command
func AttachOpts(cmd *cobra.Command, cmdOpts *CommandOptions) {
	flags := cmd.Flags()
	flags.BoolVarP(&cmdOpts.Copy, "copy", "c", false, "copy sections (all sections copied if not specified by args)")
	flags.BoolVarP(&cmdOpts.Delete, "delete", "d", false, "delete sections after copy (no-op without copy)")
}

// MakeUse returns the text for the Use field of a cobra.Command
func MakeUse(cmdName string) string {
	return fmt.Sprintf("%s [flags] [sections]...", cmdName)
}

// MakeLong returns the text for the Long field of a cobra.Command
func MakeLong(text string) string {
	trailing := "use arguments to specify sections to copy/delete from today (all by default)"
	return fmt.Sprintf("%s\n%s", text, trailing)
}
