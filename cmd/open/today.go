package open

import (
	"time"

	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/spf13/cobra"
)

// CreateTodayCmd creates the today subcommand
func CreateTodayCmd() *cobra.Command {
	cmdOpts := commandOptions{}
	cmd := &cobra.Command{
		Use:   "today",
		Short: "open today's note",
		Long:  "open a text based note template for today",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := config.LoadOrCreate()
			if err != nil {
				return err
			}

			date := time.Now()
			copyDate := date.Add(-24 * time.Hour)
			return run(opts, cmdOpts, date, copyDate)
		},
	}
	attachOpts(cmd, &cmdOpts)
	return cmd
}
