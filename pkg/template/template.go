package template

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/config"
	"github.com/pkg/errors"
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

// CopySectionContents copies the contents of the specified section from a source template by
// appending to the contents of the receiver's section
func (t *Template) CopySectionContents(src *Template, sectionName string) error {
	tgtSec, err := t.getSection(sectionName)
	if err != nil {
		return errors.Wrap(err, "failed to find section in target")
	}
	srcSec, err := src.getSection(sectionName)
	if err != nil {
		return errors.Wrap(err, "failed to find section in source")
	}
	tgtSec.contents = insert(tgtSec.contents, srcSec.contents)
	return nil
}

// DeleteSectionContents deletes the contents of a specified section
func (t *Template) DeleteSectionContents(sectionName string) error {
	sec, err := t.getSection(sectionName)
	if err != nil {
		return fmt.Errorf("section [%s] does not exist", sectionName)
	}
	sec.deleteContents()
	return nil
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
		str += fmt.Sprintf("%s%s", name, body)
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

// insert inserts contents into tgt before any trailing empty elements, omitting trailing empty
// elements of contents
func insert(tgt []string, contents []string) []string {
	if len(contents) == 0 {
		return tgt
	}
	if len(tgt) == 0 {
		return contents
	}

	contentsIdx := getLastPopulatedIndex(contents) + 1
	insertIdx := getLastPopulatedIndex(tgt) + 1

	updated := []string{}
	updated = append(updated, tgt[:insertIdx]...)
	updated = append(updated, contents[:contentsIdx]...)
	updated = append(updated, tgt[insertIdx:]...)
	return updated
}

func getLastPopulatedIndex(s []string) int {
	ln := len(s)
	for i := 0; i < ln; i++ {
		idx := ln - i - 1
		if s[idx] != "\n" && s[idx] != "" {
			return idx
		}
	}
	return -1
}
