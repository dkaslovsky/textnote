package cmd

import (
	"github.com/dkaslovsky/TextNote/cmd/archive"
	"github.com/dkaslovsky/TextNote/cmd/open"
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
			// run the today command with default options as the default application command
			return open.CreateTodayCmd().Execute()
		},
	}

	cmd.AddCommand(
		open.CreateTodayCmd(),
		open.CreateTomorrowCmd(),
		open.CreateYesterdayCmd(),
		archive.CreateArchiveCmd(),
	)

	return cmd.Execute()
}
