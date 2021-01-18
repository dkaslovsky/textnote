package today

import (
	"time"

	"github.com/dkaslovsky/textnote/cmd/open"
	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/spf13/cobra"
)

// CreateTodayCmd creates the today subcommand
func CreateTodayCmd() *cobra.Command {
	cmdOpts := open.CommandOptions{}
	cmd := &cobra.Command{
		Use:   open.MakeUse("today"),
		Short: "open today's note",
		Long:  open.MakeLong("open a note template for today"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := config.LoadOrCreate()
			if err != nil {
				return err
			}

			date := time.Now()
			copyDate := date.Add(-24 * time.Hour)
			return open.Run(opts, cmdOpts, args, date, copyDate)
		},
	}
	open.AttachOpts(cmd, &cmdOpts)
	return cmd
}
