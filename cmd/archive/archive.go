package archive

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dkaslovsky/textnote/pkg/archive"
	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/dkaslovsky/textnote/pkg/file"
	"github.com/spf13/cobra"
)

type commandOptions struct {
	Delete  bool
	NoWrite bool
}

// CreateArchiveCmd creates the today subcommand
func CreateArchiveCmd() *cobra.Command {
	cmdOpts := commandOptions{}
	cmd := &cobra.Command{
		Use:   "archive",
		Short: "organize an archive of notes",
		Long:  "organize notes into time-based archive groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := config.LoadOrCreate()
			if err != nil {
				return err
			}

			archiver := archive.NewArchiver(opts, file.NewReadWriter(), time.Now())

			files, err := ioutil.ReadDir(config.AppDir)
			if err != nil {
				return err
			}

			// add files to archiver
			archived := []string{}
			for _, f := range files {
				if !archive.ShouldArchive(f, opts.Archive.FilePrefix) {
					log.Printf("file [%s] will not be archived", f.Name())
					continue
				}
				err := archiver.Add(f)
				if err != nil {
					log.Printf("file [%s] will not be archived: %s", f.Name(), err)
					continue
				}
				archived = append(archived, f.Name())
			}

			// write archive files
			if !cmdOpts.NoWrite {
				err = archiver.Write()
				if err != nil {
					return err
				}
			}

			// return if not deleting archived files
			if !cmdOpts.Delete {
				return nil
			}

			// delete individual archived files
			for _, f := range archived {
				err = os.Remove(filepath.Join(config.AppDir, f))
				if err != nil {
					log.Printf("unable to remove file [%s]: %s", f, err)
					continue
				}
			}

			return nil
		},
	}
	attachOpts(cmd, &cmdOpts)
	return cmd
}

func attachOpts(cmd *cobra.Command, cmdOpts *commandOptions) {
	flags := cmd.Flags()
	flags.BoolVarP(&cmdOpts.Delete, "delete", "d", false, "delete individual files after archiving")
	flags.BoolVarP(&cmdOpts.NoWrite, "nowrite", "n", false, "disable writing archive file (helpful for deleting previously archived files)")
}
