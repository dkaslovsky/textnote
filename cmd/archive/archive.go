package archive

import (
	"fmt"
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
	dryRun  bool
}

// CreateArchiveCmd creates the today subcommand
func CreateArchiveCmd() *cobra.Command {
	cmdOpts := commandOptions{}
	cmd := &cobra.Command{
		Use:          "archive",
		Short:        "consolidate notes into archive files",
		Long:         "consolidate notes into monthly archive files",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := config.Load()
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
	flags.BoolVarP(&cmdOpts.noWrite, "no-write", "n", false, "disable writing archive files (helpful for deleting previously archived files)")
	flags.BoolVar(&cmdOpts.dryRun, "dry-run", false, "print file names to be deleted instead of performing deletes (other flags are ignored)")
}

func run(templateOpts config.Opts, cmdOpts commandOptions) error {
	archiver := archive.NewArchiver(templateOpts, file.NewReadWriter(), time.Now())

	files, err := os.ReadDir(templateOpts.AppDir)
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

	// print file names for dry-run
	if cmdOpts.dryRun {
		files := archiver.GetArchivedFiles()
		fmt.Printf("running \"archive --delete\" will remove [%d] files\n", len(files))
		for _, fileName := range files {
			fmt.Printf("- %s\n", fileName)
		}
		return nil
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
	numDeleted := 0
	for _, fileName := range archiver.GetArchivedFiles() {
		err = os.Remove(fileName)
		if err != nil {
			log.Printf("unable to remove file [%s]: %s", fileName, err)
			continue
		}
		numDeleted++
	}
	log.Printf("removed [%d] files after archiving", numDeleted)

	return nil
}
