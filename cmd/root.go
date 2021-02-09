package cmd

import (
	"fmt"
	"strings"

	"github.com/dkaslovsky/textnote/cmd/archive"
	"github.com/dkaslovsky/textnote/cmd/config"
	"github.com/dkaslovsky/textnote/cmd/open"
	pkgconf "github.com/dkaslovsky/textnote/pkg/config"
	"github.com/spf13/cobra"
)

// Run executes the CLI
func Run(name string, version string) error {
	cmd := &cobra.Command{
		Use:           name,
		Long:          fmt.Sprintf("Name:\n  %s - a simple tool for creating and organizing daily notes on the command line", name),
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
		config.CreateConfigCmd(),
	)

	setVersion(cmd, version)
	setHelp(cmd, name)

	return cmd.Execute()
}

func setVersion(cmd *cobra.Command, version string) {
	if version != "" {
		cmd.Version = version
		return
	}

	cmd.Version = "unavailable"
	cmd.SetVersionTemplate(
		fmt.Sprintf("%s: built from source", strings.TrimSuffix(cmd.VersionTemplate(), "\n")),
	)
}

func setHelp(cmd *cobra.Command, name string) {
	// set custom help message for the root command
	defaultHelpFunc := cmd.HelpFunc()
	cmd.SetHelpFunc(func(cmd *cobra.Command, s []string) {
		defaultHelpFunc(cmd, s)
		if cmd.Name() != name {
			return
		}
		if description := pkgconf.DescribeEnvVars(); description != "" {
			fmt.Printf("\nOverride configuration using environment variables:%s", description)
		}
	})
}
