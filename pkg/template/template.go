package template

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/config"
)

// Template contains the structure of a TextNote
type Template struct {
	opts       config.Opts
	date       time.Time
	sections   []*section
	sectionIdx map[string]int
}

// NewTemplate constructs a new template
func NewTemplate(opts config.Opts) *Template {
	t := &Template{
		opts:       opts,
		sections:   []*section{},
		sectionIdx: map[string]int{},
	}
	for i, sectionName := range opts.Section.Names {
		t.sections = append(t.sections, newSection(sectionName))
		t.sectionIdx[sectionName] = i
	}
	return t
}

// SetDate is a setter for a Template's date field
func (t *Template) SetDate(date time.Time) {
	t.date = date
}

func (t *Template) Write(w io.Writer) error {
	if t.date.IsZero() {
		return fmt.Errorf("must set date before writing a template")
	}
	_, err := w.Write([]byte(t.string()))
	return err
}

// CopySectionContents copies the contents of the specified section from a source template by
// appending to the contents of the target's section
func (t *Template) CopySectionContents(src *Template, sectionName string) error {
	tgtIdx, found := t.sectionIdx[sectionName]
	if !found {
		return fmt.Errorf("section [%s] not found in template", sectionName)
	}
	srcIdx, found := src.sectionIdx[sectionName]
	if !found {
		return fmt.Errorf("section [%s] not found in source template", sectionName)
	}

	tgtSec := t.sections[tgtIdx]
	srcSec := src.sections[srcIdx]
	tgtSec.contents = append(tgtSec.contents, srcSec.contents...)
	return nil
}

// MoveSectionContents moves the contents of the specified section from a source template by
// appending to the contents of the target's section and deleting from the source
func (t *Template) MoveSectionContents(src *Template, sectionName string) error {
	tgtIdx, found := t.sectionIdx[sectionName]
	if !found {
		return fmt.Errorf("section [%s] not found in template", sectionName)
	}
	srcIdx, found := src.sectionIdx[sectionName]
	if !found {
		return fmt.Errorf("section [%s] not found in source template", sectionName)
	}

	tgtSec := t.sections[tgtIdx]
	srcSec := src.sections[srcIdx]
	tgtSec.contents = append(tgtSec.contents, srcSec.contents...)
	srcSec.deleteContents()
	return nil
}

// GetFirstSectionFirstLine returns the first line of content of the first Section (used when opening with Vim)
func (t *Template) GetFirstSectionFirstLine() int {
	return t.opts.Header.TrailingNewlines + 3
}

// GetFilePath generates a name for a file based on the template date
func (t *Template) GetFilePath() string {
	fileName := fmt.Sprintf("%s.%s", t.date.Format(t.opts.File.TimeFormat), fileExt)
	return filepath.Join(t.opts.AppDir, fileName)
}

func (t *Template) string() string {
	str := t.makeHeader()
	for _, section := range t.sections {
		str += section.string(t.opts.Section.Prefix, t.opts.Section.Suffix, t.opts.Section.TrailingNewlines)
	}
	return str
}

func (t *Template) makeHeader() string {
	return fmt.Sprintf("%s%s%s\n%s",
		t.opts.Header.Prefix,
		t.date.Format(t.opts.Header.TimeFormat),
		t.opts.Header.Suffix,
		strings.Repeat("\n", t.opts.Header.TrailingNewlines),
	)
}

// section is a named section of a note
type section struct {
	name     string
	contents []string // use a slice in case we want to treat contents as a list of bulleted items
}

// newSection constructs a Section
func newSection(name string, contents ...string) *section {
	return &section{
		name:     name,
		contents: contents,
	}
}

func (s *section) appendContents(contents string) {
	s.contents = append(s.contents, contents)
}

func (s *section) deleteContents() {
	s.contents = []string{}
}

func (s *section) string(prefix string, suffix string, trailing int) string {
	str := fmt.Sprintf("%s%s%s\n", prefix, s.name, suffix)
	for _, content := range s.contents {
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		str += content
	}
	return str + strings.Repeat("\n", trailing)
}

// // CleanNewlines mutates a section to remove all newlines
// func (s *Section) CleanNewlines() {
//	regexNewlines    = regexp.MustCompile(`\n{2,}`)
// 	for i, content := range s.Contents {
// 		s.Contents[i] = strings.Trim(regexNewlines.ReplaceAllString(content, "\n"), "\n")
// 	}
// }
