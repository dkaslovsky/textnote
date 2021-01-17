package open

import (
	"time"

	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/spf13/cobra"
)

// CreateTomorrowCmd creates the tomorrow subcommand
func CreateTomorrowCmd() *cobra.Command {
	cmdOpts := commandOptions{}
	cmd := &cobra.Command{
		Use:   "tomorrow",
		Short: "open tomorrow's note",
		Long:  "open a note template for tomorrow",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := config.LoadOrCreate()
			if err != nil {
				return err
			}

			copyDate := time.Now()
			date := copyDate.Add(24 * time.Hour)
			return run(opts, cmdOpts, date, copyDate)
		},
	}
	attachOpts(cmd, &cmdOpts)
	return cmd
}
