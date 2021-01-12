package cmd

import (
	"fmt"

	"github.com/dkaslovsky/textnote/cmd/archive"
	"github.com/dkaslovsky/textnote/cmd/open"
	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/spf13/cobra"
)

var appName = "textnote"

// Run executes the CLI
func Run() error {
	cmd := &cobra.Command{
		Use:           appName,
		Short:         "open today's note",
		Long:          "open a note template for today",
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
		open.CreateNextCmd(),
		archive.CreateArchiveCmd(),
	)

	// custom help message
	defaultHelpFunc := cmd.HelpFunc()
	cmd.SetHelpFunc(func(cmd *cobra.Command, s []string) {
		defaultHelpFunc(cmd, s)
		if cmd.Name() == appName {
			description := config.DescribeEnvVars()
			if description != "" {
				fmt.Printf("\nOverride configuration using environment variables:%s", description)
			}
		}
	})

	return cmd.Execute()
}
