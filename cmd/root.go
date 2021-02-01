package cmd

import (
	"fmt"

	"github.com/dkaslovsky/textnote/cmd/open"

	"github.com/dkaslovsky/textnote/cmd/archive"
	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/spf13/cobra"
)

// Run executes the CLI
func Run(name string, version string) error {
	cmd := &cobra.Command{
		Use:           name,
		Long:          fmt.Sprintf("Name:\n  %s - a simple tool for creating and organizing daily notes on the command line", name),
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// run the open command with default options as the default application command
			return open.CreateOpenCmd().Execute()
		},
	}

	cmd.AddCommand(
		open.CreateOpenCmd(),
		archive.CreateArchiveCmd(),
	)

	// set custom help message for the root command
	defaultHelpFunc := cmd.HelpFunc()
	cmd.SetHelpFunc(func(cmd *cobra.Command, s []string) {
		defaultHelpFunc(cmd, s)
		if cmd.Name() != name {
			return
		}
		if description := config.DescribeEnvVars(); description != "" {
			fmt.Printf("\nOverride configuration using environment variables:%s", description)
		}
	})

	return cmd.Execute()
}
