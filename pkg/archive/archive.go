package archive

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/config"
	"github.com/dkaslovsky/TextNote/pkg/file"
	"github.com/dkaslovsky/TextNote/pkg/template"
	"github.com/pkg/errors"
)

// Archiver consolidates TextNotes into archive files
type Archiver struct {
	opts config.Opts
	date time.Time

	Months map[string]*template.MonthArchiveTemplate
}

// NewArchiver constructs a new Archiver
func NewArchiver(opts config.Opts, date time.Time) *Archiver {
	return &Archiver{
		opts: opts,
		date: date,

		Months: map[string]*template.MonthArchiveTemplate{},
	}
}

type fileInfo interface {
	Name() string
	IsDir() bool
}

// Add adds a file to the archive
func (a *Archiver) Add(f fileInfo) error {
	if a.shouldNotArchive(f) {
		return nil
	}

	fileDate, err := parseFileName(f.Name(), a.opts.File.TimeFormat)
	if err != nil {
		return fmt.Errorf("cannot add unparsable file name [%s] to archive", f.Name())
	}

	// recent files are not archived
	if a.date.Sub(fileDate).Hours() <= float64(a.opts.Archive.AfterDays*24) {
		return nil
	}

	t := template.NewTemplate(a.opts, fileDate)
	err = file.Read(t)
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

func (a *Archiver) Write(write func(file.ReadWriteable) error, exists func(file.ReadWriteable) bool) error {
	for _, t := range a.Months {
		if exists(t) {
			existing := template.NewMonthArchiveTemplate(a.opts, t.GetDate())
			err := file.Read(existing)
			if err != nil {
				return errors.Wrapf(err, "unable to open existing archive file [%s]", existing.GetFilePath())
			}
			err = t.Merge(existing)
			if err != nil {
				return errors.Wrapf(err, "unable to from merge existing archive file [%s]", existing.GetFilePath())
			}
		}

		err := write(t)
		if err != nil {
			return errors.Wrapf(err, "failed to write archive file [%s]", t.GetFilePath())
		}
	}
	return nil
}

func (a *Archiver) shouldNotArchive(f fileInfo) bool {
	switch {
	// skip archive files
	case strings.HasPrefix(f.Name(), a.opts.Archive.FilePrefix):
		return true
	// skip hidden files
	case strings.HasPrefix(f.Name(), "."):
		return true
	// skip directories
	case f.IsDir():
		return true
	default:
		return false
	}
}

func parseFileName(fileName string, format string) (time.Time, error) {
	name := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	return time.Parse(format, name)
}
