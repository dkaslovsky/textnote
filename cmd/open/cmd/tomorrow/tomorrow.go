package tomorrow

import (
	"time"

	"github.com/dkaslovsky/textnote/cmd/open"
	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/spf13/cobra"
)

// CreateTomorrowCmd creates the tomorrow subcommand
func CreateTomorrowCmd() *cobra.Command {
	cmdOpts := open.CommandOptions{}
	cmd := &cobra.Command{
		Use:   open.MakeUse("tomorrow"),
		Short: "open tomorrow's note",
		Long:  open.MakeLong("open a note template for tomorrow"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := config.LoadOrCreate()
			if err != nil {
				return err
			}

			copyDate := time.Now()
			date := copyDate.Add(24 * time.Hour)
			return open.Run(opts, cmdOpts, args, date, copyDate)
		},
	}
	open.AttachOpts(cmd, &cmdOpts)
	return cmd
}
