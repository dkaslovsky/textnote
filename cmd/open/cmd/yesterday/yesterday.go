package yesterday

import (
	"fmt"
	"time"

	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/dkaslovsky/textnote/pkg/file"
	"github.com/dkaslovsky/textnote/pkg/template"
	"github.com/spf13/cobra"
)

// CreateYesterdayCmd creates the yesterday subcommand
func CreateYesterdayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "yesterday",
		Short:                 "open yesterday's note",
		Long:                  "open a note template for yesterday",
		Args:                  cobra.MaximumNArgs(0),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := config.LoadOrCreate()
			if err != nil {
				return err
			}

			rw := file.NewReadWriter()

			date := time.Now().Add(-24 * time.Hour)
			t := template.NewTemplate(opts, date)
			if !rw.Exists(t) {
				return fmt.Errorf("file [%s] for template does not exist", t.GetFilePath())
			}
			return file.OpenInVim(t)
		},
	}
	return cmd
}
