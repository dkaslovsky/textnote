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
	date     time.Time
	sections []*section
	opts     config.Opts
}

// NewTemplate constructs a new template
func NewTemplate(opts config.Opts) *Template {
	t := &Template{
		sections: []*section{},
		opts:     opts,
	}
	for _, sectionName := range opts.Section.Names {
		t.AddSection(newSection(sectionName))
	}
	return t
}

// SetDate is a setter for a Template's date field
func (t *Template) SetDate(date time.Time) {
	t.date = date
}

// AddSection adds a section to a Template
func (t *Template) AddSection(s *section) {
	t.sections = append(t.sections, s)
}

func (t *Template) Write(w io.Writer) error {
	if t.date.IsZero() {
		return fmt.Errorf("must set date before writing a template")
	}
	_, err := w.Write([]byte(t.string()))
	return err
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
	Name     string
	Contents []string // use a slice in case we want to treat contents as a list of bulleted items
}

// newSection constructs a Section
func newSection(name string, contents ...string) *section {
	return &section{
		Name:     name,
		Contents: contents,
	}
}

func (s *section) string(prefix string, suffix string, trailing int) string {
	str := fmt.Sprintf("%s%s%s\n", prefix, s.Name, suffix)
	for _, content := range s.Contents {
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
