package open

import (
	"fmt"
	"time"

	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/spf13/cobra"
)

type nextCommandOptions struct {
	CmdOpts commandOptions
	Weekday int
}

// CreateNextCmd creates the next subcommand
func CreateNextCmd() *cobra.Command {
	nextCmdOpts := nextCommandOptions{}
	cmd := &cobra.Command{
		Use:   "next",
		Short: "open a note for the next specified day",
		Long:  "open a note for the next specified day of the week",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := config.LoadOrCreate()
			if err != nil {
				return err
			}

			if nextCmdOpts.Weekday < 0 || nextCmdOpts.Weekday > 6 {
				return fmt.Errorf("invalid day of the week [%d], must be in range [0, 6]", nextCmdOpts.Weekday)
			}
			targetDay := time.Weekday(nextCmdOpts.Weekday)

			copyDate := time.Now()
			date := time.Now()

			if date.Weekday() == targetDay {
				date = date.Add(24 * time.Hour)
			}
			for date.Weekday() != targetDay {
				date = date.Add(24 * time.Hour)
			}

			return run(opts, nextCmdOpts.CmdOpts, date, copyDate)
		},
	}
	attachNextOpts(cmd, &nextCmdOpts)
	return cmd
}

func attachNextOpts(cmd *cobra.Command, nextCmdOpts *nextCommandOptions) {
	flags := cmd.Flags()
	flags.IntVarP(&nextCmdOpts.Weekday, "weekday", "w", 1, "day of the week to open (0=Sunday, 1=Monday, ...")
	attachOpts(cmd, &nextCmdOpts.CmdOpts)
}
