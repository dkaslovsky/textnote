package initialize

import (
	"github.com/spf13/cobra"

	"github.com/dkaslovsky/textnote/pkg/config"
)

// CreateInitCmd creates the init subcommand
func CreateInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize the application",
		Long:  "initialize the application's required directories and files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return config.InitApp()
		},
	}
	return cmd
}
