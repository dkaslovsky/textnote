package template

import (
	"fmt"
	"path/filepath"
	"sort"
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
		Template: NewTemplate(opts, date),
	}
}

// GetFilePath generates a name for a file based on the template date
func (t *MonthArchiveTemplate) GetFilePath() string {
	fileName := fmt.Sprintf("%s%s.%s", ArchiveFilePrefix, t.date.Format(t.opts.Archive.MonthTimeFormat), fileExt)
	return filepath.Join(t.opts.AppDir, fileName)
}

func (t *MonthArchiveTemplate) string() string {
	str := t.makeHeader()
	for _, section := range t.sections {
		sort.Strings(section.contents)
		name := section.getNameString(t.opts.Section.Prefix, t.opts.Section.Suffix)
		body := section.getBodyString()
		// default to trailing whitespace for empty body
		if len(body) == 0 {
			body = strings.Repeat("\n", t.opts.Section.TrailingNewlines)
		}
		str += name + body
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
	return makeArchiveSectionContentPrefix(t.opts.Archive, date)
}

// YearArchiveTemplate contains the structure of a TextNote year archive
type YearArchiveTemplate struct {
	*Template
}

// NewYearArchiveTemplate constructs a new YearArchiveTemplate
func NewYearArchiveTemplate(opts config.Opts, date time.Time) *YearArchiveTemplate {
	return &YearArchiveTemplate{
		Template: NewTemplate(opts, date),
	}
}

// GetFilePath generates a name for a file based on the template date
func (t *YearArchiveTemplate) GetFilePath() string {
	fileName := fmt.Sprintf("%s%s.%s", ArchiveFilePrefix, t.date.Format(t.opts.Archive.YearTimeFormat), fileExt)
	return filepath.Join(t.opts.AppDir, fileName)
}

// adds a sort to each section contents
func (t *YearArchiveTemplate) string() string {
	str := t.makeHeader()
	for _, section := range t.sections {
		sort.Strings(section.contents)
		name := section.getNameString(t.opts.Section.Prefix, t.opts.Section.Suffix)
		body := section.getBodyString()
		// default to trailing whitespace for empty body
		if len(body) == 0 {
			body = strings.Repeat("\n", t.opts.Section.TrailingNewlines)
		}
		str += name + body
	}
	return str
}

func (t *YearArchiveTemplate) makeHeader() string {
	return fmt.Sprintf("%s%s%s\n%s",
		t.opts.Archive.HeaderPrefix,
		t.date.Format(t.opts.Archive.YearTimeFormat),
		t.opts.Archive.HeaderSuffix,
		strings.Repeat("\n", t.opts.Header.TrailingNewlines),
	)
}

func (t *YearArchiveTemplate) makeSectionContentPrefix(date time.Time) string {
	return makeArchiveSectionContentPrefix(t.opts.Archive, date)
}

func makeArchiveSectionContentPrefix(opts config.ArchiveOpts, date time.Time) string {
	return fmt.Sprintf("%s%s%s",
		opts.SectionContentPrefix,
		date.Format(opts.SectionContentTimeFormat),
		opts.SectionContentSuffix,
	)
}
