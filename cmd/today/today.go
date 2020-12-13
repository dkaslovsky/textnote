package today

import (
	"fmt"
	"os"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/config"
	"github.com/dkaslovsky/TextNote/pkg/file"
	"github.com/dkaslovsky/TextNote/pkg/template"
	"github.com/spf13/cobra"
)

// CommandOptions are command options
type CommandOptions struct {
	Copy   []string
	Delete bool
}

// CreateTodayCmd creates the package's subcommand
func CreateTodayCmd() *cobra.Command {
	cmdOpts := CommandOptions{}
	cmd := &cobra.Command{
		Use:   "today",
		Short: "open today's note",
		Long:  "open a text based note template for today",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := config.LoadOrCreate()
			if err != nil {
				return err
			}

			shouldCopy := len(cmdOpts.Copy) > 0

			date := time.Now()
			t := template.NewTemplate(opts)
			t.SetDate(date)

			if !shouldCopy {
				err = file.WriteIfNotExists(t)
				if err != nil {
					return err
				}
				return file.OpenInEditor(t)
			}

			if file.Exists(t) {
				err = file.Read(t)
				if err != nil {
					return err
				}
			}

			srcDate := date.Add(-24 * time.Hour)
			src := template.NewTemplate(opts)
			src.SetDate(srcDate)
			err = file.Read(src)
			if err != nil {
				return err
			}

			copyOp := t.CopySectionContents
			if cmdOpts.Delete {
				copyOp = t.MoveSectionContents
			}

			for _, sectionName := range cmdOpts.Copy {
				err = copyOp(src, sectionName)
				if err != nil {
					return err
				}
			}
			if cmdOpts.Delete {
				err = file.Overwrite(src)
				if err != nil {
					return err
				}
			}
			fmt.Println("t:")
			t.Write(os.Stdout)
			err = file.Overwrite(t)
			if err != nil {
				return err
			}
			return nil
		},
	}
	attachOpts(cmd, &cmdOpts)
	return cmd
}

func attachOpts(cmd *cobra.Command, cmdOpts *CommandOptions) {
	flags := cmd.Flags()
	flags.StringSliceVarP(&cmdOpts.Copy, "copy", "c", []string{}, "section names to copy")
	flags.BoolVarP(&cmdOpts.Delete, "delete", "d", false, "delete previous day's section after copy (no-op without copy")
}
