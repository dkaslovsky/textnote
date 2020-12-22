package template

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/config"
)

// ArchiveFilePrefix is the prefix attached to the file name of all archive files
const ArchiveFilePrefix = "archive-"

// MonthArchiveTemplate contains the structure of a TextNote month archive
type MonthArchiveTemplate struct {
	*Template
}

// NewMonthArchiveTemplate constructs a new MonthArchiveTemplate
func NewMonthArchiveTemplate(opts config.Opts, date time.Time) *MonthArchiveTemplate {
	return &MonthArchiveTemplate{
		NewTemplate(opts, date),
	}
}

func (t *MonthArchiveTemplate) Write(w io.Writer) error {
	_, err := w.Write([]byte(t.string()))
	return err
}

// GetFilePath generates a name for a file based on the template date
func (t *MonthArchiveTemplate) GetFilePath() string {
	fileName := fmt.Sprintf("%s%s.%s", ArchiveFilePrefix, t.date.Format(t.opts.Archive.MonthTimeFormat), fileExt)
	return filepath.Join(t.opts.AppDir, fileName)
}

func (t *MonthArchiveTemplate) string() string {
	str := t.makeHeader()
	for _, section := range t.sections {
		name := section.getNameString(t.opts.Section.Prefix, t.opts.Section.Suffix)
		body := section.getBodyString()
		str += fmt.Sprintf("%s%s%s",
			name,
			body,
			strings.Repeat("\n", t.opts.Section.TrailingNewlines),
		)
	}
	return str
}

func (t *MonthArchiveTemplate) makeHeader() string {
	return fmt.Sprintf("%s%s%s\n%s",
		t.opts.Archive.HeaderPrefix,
		t.date.Format(t.opts.Archive.MonthTimeFormat),
		t.opts.Archive.HeaderSuffix,
		strings.Repeat("\n", t.opts.Header.TrailingNewlines),
	)
}

func (t *MonthArchiveTemplate) makeSectionContentPrefix(date time.Time) string {
	return fmt.Sprintf("\n%s%s%s",
		t.opts.Archive.SectionContentPrefix,
		date.Format(t.opts.Archive.SectionContentTimeFormat),
		t.opts.Archive.SectionContentSuffix,
	)
}
