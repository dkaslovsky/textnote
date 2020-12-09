package cmd

import (
	"github.com/dkaslovsky/TextNote/cmd/open"
	"github.com/spf13/cobra"
)

func Run() error {
	cmd := &cobra.Command{
		Use:   "textnote",
		Short: "open today's note",
		Long:  "open a text based note template for today",
		Run: func(cmd *cobra.Command, args []string) {
			// run the open command as the default
			open.CreateCmd().Execute()
		},
	}
	cmd.AddCommand(open.CreateCmd())
	return cmd.Execute()
}
