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

const weekHours = 24 * 7

type Archiver struct {
	opts config.Opts
	date time.Time

	Month map[string]*template.MonthArchiveTemplate
	Year  map[string]*template.YearArchiveTemplate
}

func NewArchiver(opts config.Opts, date time.Time) *Archiver {
	return &Archiver{
		opts: opts,
		date: date,

		Month: map[string]*template.MonthArchiveTemplate{},
		Year:  map[string]*template.YearArchiveTemplate{},
	}
}

type fileInfo interface {
	Name() string
	IsDir() bool
}

func (a *Archiver) Add(f fileInfo) error {
	if shouldSkip(f) {
		return nil
	}

	fileDate, err := parseFileName(f.Name(), a.opts.File.TimeFormat)
	if err != nil {
		return fmt.Errorf("cannot add unparsable file name [%s] to archive", f.Name())
	}

	// recent files are not archived
	if a.date.Sub(fileDate).Hours() <= weekHours {
		return nil
	}

	t := template.NewTemplate(a.opts, fileDate)
	err = file.Read(t)
	if err != nil {
		return errors.Wrapf(err, "cannot add unreadable file [%s] to archive", f.Name())
	}

	if fileDate.Year() < a.date.Year() {
		key := fileDate.Format(a.opts.Archive.YearTimeFormat)
		if _, found := a.Year[key]; !found {
			a.Year[key] = template.NewYearArchiveTemplate(a.opts, fileDate)
		}

		for _, section := range a.opts.Section.Names {
			err := template.ArchiveSectionContents(t, a.Year[key], section)
			if err != nil {
				return errors.Wrapf(err, "cannot add contents from [%s] to archive", f.Name())
			}
		}
		return nil
	}

	if fileDate.Month() <= a.date.Month() {
		key := fileDate.Format(a.opts.Archive.MonthTimeFormat)
		if _, found := a.Month[key]; !found {
			a.Month[key] = template.NewMonthArchiveTemplate(a.opts, fileDate)
		}

		for _, section := range a.opts.Section.Names {
			err := template.ArchiveSectionContents(t, a.Month[key], section)
			if err != nil {
				return errors.Wrapf(err, "cannot add contents from [%s] to archive", f.Name())
			}
		}
		return nil
	}

	return fmt.Errorf("cannot archive file with date [%v]", fileDate)
}

func shouldSkip(f fileInfo) bool {
	// skip archive files
	if strings.HasPrefix(f.Name(), template.ArchiveFilePrefix) {
		return true
	}
	// skip hidden files
	if strings.HasPrefix(f.Name(), ".") {
		return true
	}
	// skip directories
	if f.IsDir() {
		return true
	}
	return false
}

func parseFileName(fileName string, format string) (time.Time, error) {
	name := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	return time.Parse(format, name)
}
