package archive

import (
	"fmt"
	"log"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/archive"
	"github.com/dkaslovsky/TextNote/pkg/config"
	"github.com/dkaslovsky/TextNote/pkg/file"
	"github.com/spf13/cobra"
)

// CreateArchiveCmd creates the today subcommand
func CreateArchiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive",
		Short: "organize an archive of notes",
		Long:  "organize notes into time-based archive groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := config.LoadOrCreate()
			if err != nil {
				return err
			}

			// fake now and files for now
			now := time.Now().AddDate(0, 2, -6)
			files := generateFileNames(opts, now)

			archiver := archive.NewArchiver(opts, file.NewReadWriter(), now)
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
			}

			err = archiver.Write()
			if err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

// fake it...
type fakeFileInfo struct {
	name  string
	isDir bool
}

func (ffi fakeFileInfo) Name() string {
	return ffi.name
}

func (ffi fakeFileInfo) IsDir() bool {
	return ffi.isDir
}

func generateFileNames(opts config.Opts, now time.Time) []fakeFileInfo {
	ffi := []fakeFileInfo{}
	ffi = append(ffi, fakeFileInfo{".config", false})

	end := now
	start := end.AddDate(0, -2, 2)
	for end.After(start) {
		tm := start.Format(opts.File.TimeFormat)
		name := fmt.Sprintf("%s.txt", tm)
		fi := fakeFileInfo{name, false}
		ffi = append(ffi, fi)
		start = start.AddDate(0, 0, 1)
	}

	archive := fakeFileInfo{opts.Archive.FilePrefix + ffi[2].name, false}
	ffi = append(ffi, archive)
	dir := fakeFileInfo{"somedir", true}
	ffi = append(ffi, dir)
	return ffi
}
