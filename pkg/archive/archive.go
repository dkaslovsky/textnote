package archive

import (
	"fmt"
	"log"
	"time"

	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/dkaslovsky/textnote/pkg/file"
	"github.com/dkaslovsky/textnote/pkg/template"
)

// Archiver consolidates templates into archives
type Archiver struct {
	opts config.Opts
	rw   readWriter
	date time.Time // timestamp for calculating if a file is old enough to be archived

	// monthArchives maintains a map of formatted month timestamp to the corresponding archive
	monthArchives map[string]*template.MonthArchiveTemplate
	// archivedFiles maintains the file names that have been archived
	archivedFiles []string
}

// NewArchiver constructs a new Archiver
func NewArchiver(opts config.Opts, rw readWriter, date time.Time) *Archiver {
	return &Archiver{
		opts: opts,
		rw:   rw,
		date: date,

		monthArchives: map[string]*template.MonthArchiveTemplate{},
		archivedFiles: []string{},
	}
}

// Add adds a template corresponding to a date to the archive
func (a *Archiver) Add(date time.Time) error {
	// recent files are not archived
	if a.date.Sub(date).Hours() <= float64(a.opts.Archive.AfterDays*24) {
		return nil
	}

	t := template.NewTemplate(a.opts, date)
	err := a.rw.Read(t)
	if err != nil {
		return fmt.Errorf("cannot add unreadable file [%s] to archive: %w", t.GetFilePath(), err)
	}

	monthKey := date.Format(a.opts.Archive.MonthTimeFormat)
	if _, found := a.monthArchives[monthKey]; !found {
		a.monthArchives[monthKey] = template.NewMonthArchiveTemplate(a.opts, date)
	}

	archive := a.monthArchives[monthKey]
	for _, section := range a.opts.Section.Names {
		err := archive.ArchiveSectionContents(t, section)
		if err != nil {
			return fmt.Errorf("cannot add contents from [%s] to archive: %w", t.GetFilePath(), err)
		}
	}

	a.archivedFiles = append(a.archivedFiles, t.GetFilePath())
	return nil
}

// Write writes all of the archive templates stored in the Archiver
func (a *Archiver) Write() error {
	for _, t := range a.monthArchives {
		if a.rw.Exists(t) {
			existing := template.NewMonthArchiveTemplate(a.opts, t.GetDate())
			err := a.rw.Read(existing)
			if err != nil {
				return fmt.Errorf("unable to open existing archive file [%s]: %w", existing.GetFilePath(), err)
			}
			err = t.Merge(existing)
			if err != nil {
				return fmt.Errorf("unable to from merge existing archive file [%s] %w", existing.GetFilePath(), err)
			}
		}

		err := a.rw.Overwrite(t)
		if err != nil {
			return fmt.Errorf("failed to write archive file [%s]: %w", t.GetFilePath(), err)
		}
		log.Printf("wrote archive file [%s]", t.GetFilePath())
	}
	return nil
}

// GetArchivedFiles returns the files that have been archived
func (a *Archiver) GetArchivedFiles() []string {
	return a.archivedFiles
}

// readWriter is the interface for executing file operations
type readWriter interface {
	Read(file.ReadWriteable) error
	Overwrite(file.ReadWriteable) error
	Exists(file.ReadWriteable) bool
}
