package archive

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dkaslovsky/textnote/pkg/archive"
	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/dkaslovsky/textnote/pkg/file"
	"github.com/dkaslovsky/textnote/pkg/template"
	"github.com/spf13/cobra"
)

type commandOptions struct {
	delete  bool
	noWrite bool
}

// CreateArchiveCmd creates the today subcommand
func CreateArchiveCmd() *cobra.Command {
	cmdOpts := commandOptions{}
	cmd := &cobra.Command{
		Use:   "archive",
		Short: "consolidate notes into archive files",
		Long:  "consolidate notes into monthly archive files",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := config.LoadOrCreate()
			if err != nil {
				return err
			}
			return run(opts, cmdOpts)
		},
	}
	attachOpts(cmd, &cmdOpts)
	return cmd
}

func attachOpts(cmd *cobra.Command, cmdOpts *commandOptions) {
	flags := cmd.Flags()
	flags.BoolVarP(&cmdOpts.delete, "delete", "x", false, "delete individual files after archiving")
	flags.BoolVarP(&cmdOpts.noWrite, "nowrite", "n", false, "disable writing archive file (helpful for deleting previously archived files)")
}

func run(templateOpts config.Opts, cmdOpts commandOptions) error {
	archiver := archive.NewArchiver(templateOpts, file.NewReadWriter(), time.Now())

	files, err := ioutil.ReadDir(config.AppDir)
	if err != nil {
		return err
	}

	// add template files to archiver
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		// parse date from template file name, skipping non-template files
		templateDate, ok := template.ParseTemplateFileName(f.Name(), templateOpts.File)
		if !ok {
			continue
		}

		err := archiver.Add(templateDate)
		if err != nil {
			log.Printf("skipping unarchivable file [%s]: %s", f.Name(), err)
			continue
		}
	}

	// write archive files
	if !cmdOpts.noWrite {
		err = archiver.Write()
		if err != nil {
			return err
		}
	}

	// return if not deleting archived files
	if !cmdOpts.delete {
		return nil
	}

	// delete individual archived files
	for _, fileName := range archiver.GetArchivedFiles() {
		err = os.Remove(fileName)
		if err != nil {
			log.Printf("unable to remove file [%s]: %s", fileName, err)
			continue
		}
	}

	return nil
}
