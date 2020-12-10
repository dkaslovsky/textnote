package cmd

import (
	"github.com/dkaslovsky/TextNote/cmd/open"
	"github.com/spf13/cobra"
)

// Run executes the CLI interface
func Run() error {
	cmd := &cobra.Command{
		Use:   "textnote",
		Short: "open today's note",
		Long:  "open a text based note template for today",
		RunE: func(cmd *cobra.Command, args []string) error {
			// run the today command as the default
			return open.CreateTodayCmd().Execute()
		},
	}

	cmd.AddCommand(
		open.CreateTodayCmd(),
		open.CreateTomorrowCmd(),
	)

	return cmd.Execute()
}
