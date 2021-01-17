package archive

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/dkaslovsky/textnote/pkg/file"
	"github.com/dkaslovsky/textnote/pkg/template"
	"github.com/pkg/errors"
)

// Archiver consolidates TextNotes into archive files
type Archiver struct {
	opts config.Opts
	rw   readWriter
	// timestamp used to calculate whether a file is old enough to be archived, usually time.Now()
	date time.Time

	// archive templates by month keyed by formatted month timestamp
	Months map[string]*template.MonthArchiveTemplate
}

// NewArchiver constructs a new Archiver
func NewArchiver(opts config.Opts, rw readWriter, date time.Time) *Archiver {
	return &Archiver{
		opts: opts,
		rw:   rw,
		date: date,

		Months: map[string]*template.MonthArchiveTemplate{},
	}
}

// readWriter is the interface for executing file operations
type readWriter interface {
	Read(file.ReadWriteable) error
	Overwrite(file.ReadWriteable) error
	Exists(file.ReadWriteable) bool
}

// fileInfo is the interface for adding a file to the Archvier
type fileInfo interface {
	Name() string
	IsDir() bool
}

// Add adds a file to the archive
func (a *Archiver) Add(f fileInfo) error {
	fileDate, err := parseFileName(f.Name(), a.opts.File.TimeFormat)
	if err != nil {
		return fmt.Errorf("cannot add unparsable file name [%s] to archive", f.Name())
	}

	// recent files are not archived
	if a.date.Sub(fileDate).Hours() <= float64(a.opts.Archive.AfterDays*24) {
		return nil
	}

	t := template.NewTemplate(a.opts, fileDate)
	err = a.rw.Read(t)
	if err != nil {
		return errors.Wrapf(err, "cannot add unreadable file [%s] to archive", f.Name())
	}

	monthKey := fileDate.Format(a.opts.Archive.MonthTimeFormat)
	if _, found := a.Months[monthKey]; !found {
		a.Months[monthKey] = template.NewMonthArchiveTemplate(a.opts, fileDate)
	}

	archive := a.Months[monthKey]
	for _, section := range a.opts.Section.Names {
		err := archive.ArchiveSectionContents(t, section)
		if err != nil {
			return errors.Wrapf(err, "cannot add contents from [%s] to archive", f.Name())
		}
	}
	return nil
}

func (a *Archiver) Write() error {
	for _, t := range a.Months {
		if a.rw.Exists(t) {
			existing := template.NewMonthArchiveTemplate(a.opts, t.GetDate())
			err := a.rw.Read(existing)
			if err != nil {
				return errors.Wrapf(err, "unable to open existing archive file [%s]", existing.GetFilePath())
			}
			err = t.Merge(existing)
			if err != nil {
				return errors.Wrapf(err, "unable to from merge existing archive file [%s]", existing.GetFilePath())
			}
		}

		err := a.rw.Overwrite(t)
		if err != nil {
			return errors.Wrapf(err, "failed to write archive file [%s]", t.GetFilePath())
		}
	}
	return nil
}

// ShouldArchive determines whether a file should be included in an archive
func ShouldArchive(f fileInfo, archivePrefix string) bool {
	switch {
	// skip archive files
	case strings.HasPrefix(f.Name(), archivePrefix):
		return false
	// skip hidden files
	case strings.HasPrefix(f.Name(), "."):
		return false
	// skip directories
	case f.IsDir():
		return false
	default:
		return true
	}
}

func parseFileName(fileName string, format string) (time.Time, error) {
	name := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	return time.Parse(format, name)
}
