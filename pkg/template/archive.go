package template

import (
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/config"
	"github.com/pkg/errors"
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

// Write writes the template
func (t *MonthArchiveTemplate) Write(w io.Writer) error {
	_, err := w.Write([]byte(t.string()))
	return err
}

// GetFilePath generates a full path for a file based on the template date
func (t *MonthArchiveTemplate) GetFilePath() string {
	fileName := fmt.Sprintf("%s%s.%s", ArchiveFilePrefix, t.date.Format(t.opts.Archive.MonthTimeFormat), fileExt)
	return filepath.Join(t.opts.AppDir, fileName)
}

type dateSectionGettable interface {
	GetDate() time.Time
	getSection(string) (*section, error)
}

// CopySectionContents archives the contents of the specified section from a source template by
// concatenating its contents' text and appending to the contents of the receiver's section with a
// header corresponding to the source template's date
func (t *MonthArchiveTemplate) CopySectionContents(src dateSectionGettable, sectionName string) error {
	tgtSec, err := t.getSection(sectionName)
	if err != nil {
		return errors.Wrap(err, "failed to find section in target")
	}
	srcSec, err := src.getSection(sectionName)
	if err != nil {
		return errors.Wrap(err, "failed to find section in source")
	}

	// flatten text from contents into a single string
	txt := ""
	for _, content := range srcSec.contents {
		txt += content.text
	}
	if len(txt) == 0 {
		return nil
	}

	tgtSec.contents = append(tgtSec.contents, contentItem{
		header: t.makeContentHeader(src.GetDate()),
		text:   txt,
	})
	return nil
}

// Merge merges a source MonthArchiveTemplate into the receiver
func (t *MonthArchiveTemplate) Merge(src *MonthArchiveTemplate) error {
	for sectionName := range t.sectionIdx {
		tgtSec, err := t.getSection(sectionName)
		if err != nil {
			return errors.Wrap(err, "failed to find section in target")
		}
		srcSec, err := src.getSection(sectionName)
		if err != nil {
			return errors.Wrap(err, "failed to find section in source")
		}

		tgtSec.contents = append(tgtSec.contents, srcSec.contents...)
	}
	return nil
}

func (t *MonthArchiveTemplate) string() string {
	str := t.makeHeader()
	for _, section := range t.sections {
		name := section.getNameString(t.opts.Section.Prefix, t.opts.Section.Suffix)

		section.sortContents()
		body := section.getContentString()
		body = regexp.MustCompile(`\n{2,}`).ReplaceAllString(body, "\n") // remove blank lines

		str += fmt.Sprintf("%s%s%s", name, body, strings.Repeat("\n", t.opts.Section.TrailingNewlines))
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

func (t *MonthArchiveTemplate) makeContentHeader(date time.Time) string {
	return fmt.Sprintf("%s%s%s",
		t.opts.Archive.SectionContentPrefix,
		date.Format(t.opts.Archive.SectionContentTimeFormat),
		t.opts.Archive.SectionContentSuffix,
	)
}
