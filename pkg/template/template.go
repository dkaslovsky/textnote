package template

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/config"
)

const fileExt = "txt"

// Template contains the structure of a TextNote
type Template struct {
	opts       config.Opts
	date       time.Time
	sections   []*section
	sectionIdx map[string]int
}

// NewTemplate constructs a new Template
func NewTemplate(opts config.Opts, date time.Time) *Template {
	t := &Template{
		opts:       opts,
		date:       date,
		sections:   []*section{},
		sectionIdx: map[string]int{},
	}
	for i, sectionName := range opts.Section.Names {
		t.sections = append(t.sections, newSection(sectionName))
		t.sectionIdx[sectionName] = i
	}
	return t
}

func (t *Template) Write(w io.Writer) error {
	_, err := w.Write([]byte(t.string()))
	return err
}

// GetFileStartLine returns the first line of content of the first Section (used when opening with Vim)
func (t *Template) GetFileStartLine() int {
	return t.opts.Header.TrailingNewlines + 3
}

// GetFilePath generates a name for a file based on the template date
func (t *Template) GetFilePath() string {
	fileName := fmt.Sprintf("%s.%s", t.date.Format(t.opts.File.TimeFormat), fileExt)
	return filepath.Join(t.opts.AppDir, fileName)
}

func (t *Template) getSection(name string) (sec *section, err error) {
	idx, found := t.sectionIdx[name]
	if !found {
		return sec, fmt.Errorf("section [%s] not found", name)
	}
	return t.sections[idx], nil
}

func (t *Template) string() string {
	str := t.makeHeader()
	for _, section := range t.sections {
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
	contents []string
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

func (s *section) getNameString(prefix string, suffix string) string {
	return fmt.Sprintf("%s%s%s\n", prefix, s.name, suffix)
}

func (s *section) getBodyString() string {
	body := ""
	for _, content := range s.contents {
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		body += content
	}
	return body
}
