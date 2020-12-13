package cmd

import (
	"github.com/dkaslovsky/TextNote/cmd/today"
	"github.com/spf13/cobra"
)

// Run executes the CLI
func Run() error {
	cmd := &cobra.Command{
		Use:           "textnote",
		Short:         "open today's note",
		Long:          "open a text based note template for today",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// run the today command as the default
			return today.CreateTodayCmd().Execute()
		},
	}

	cmd.AddCommand(
		today.CreateTodayCmd(),
	)

	return cmd.Execute()
}
