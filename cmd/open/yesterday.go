package open

import (
	"time"

	"github.com/dkaslovsky/TextNote/pkg/config"
	"github.com/spf13/cobra"
)

// CreateYesterdayCmd creates the yesterday subcommand
func CreateYesterdayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yesterday",
		Short: "open yesterday's note",
		Long:  "open a text based note template from yesterday",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := config.LoadOrCreate()
			if err != nil {
				return err
			}
			date := time.Now().Add(-24 * time.Hour)
			return open(opts, date)
		},
	}
	return cmd
}
