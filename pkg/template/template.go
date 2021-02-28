package template

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/pkg/errors"
)

// Template contains the structure of a note
type Template struct {
	opts       config.Opts
	date       time.Time
	sections   []*section
	sectionIdx map[string]int // map of section name to index in sections slice
}

// NewTemplate constructs a new Template
func NewTemplate(opts config.Opts, date time.Time) *Template {
	t := &Template{
		opts:       opts,
		date:       date,
		sections:   []*section{},
		sectionIdx: map[string]int{},
	}
	for idx, sectionName := range opts.Section.Names {
		t.sections = append(t.sections, newSection(sectionName))
		t.sectionIdx[sectionName] = idx
	}
	return t
}

// Write writes the template
func (t *Template) Write(w io.Writer) error {
	_, err := w.Write([]byte(t.string()))
	return err
}

// GetDate returns the template's date
func (t *Template) GetDate() time.Time {
	return t.date
}

// GetFileCursorLine returns the line at which to place the cursor when opening the template
func (t *Template) GetFileCursorLine() int {
	return t.opts.File.CursorLine
}

// GetFilePath generates a full path for a file based on the template date
func (t *Template) GetFilePath() string {
	name := filepath.Join(t.opts.AppDir, t.date.Format(t.opts.File.TimeFormat))
	if t.opts.File.Ext == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", name, t.opts.File.Ext)
}

// sectionGettable is the interface for getting a section
type sectionGettable interface {
	getSection(string) (*section, error)
}

// CopySectionContents copies the contents of the specified section from a source template by
// appending to the contents of the receiver's section
func (t *Template) CopySectionContents(src sectionGettable, sectionName string) error {
	tgtSec, err := t.getSection(sectionName)
	if err != nil {
		return errors.Wrap(err, "failed to find section in target")
	}
	srcSec, err := src.getSection(sectionName)
	if err != nil {
		return errors.Wrap(err, "failed to find section in source")
	}
	tgtSec.contents = append(tgtSec.contents, srcSec.contents...)
	return nil
}

// DeleteSectionContents deletes the contents of a specified section
func (t *Template) DeleteSectionContents(sectionName string) error {
	sec, err := t.getSection(sectionName)
	if err != nil {
		return errors.Wrap(err, "cannot delete section")
	}
	sec.deleteContents()
	return nil
}

// Load populates a Template from the contents of a reader
func (t *Template) Load(r io.Reader) error {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "error loading template")
	}
	sectionText := string(raw)

	sectionNameRegex, err := getSectionNameRegex(t.opts.Section.Prefix, t.opts.Section.Suffix)
	if err != nil {
		return errors.Wrap(err, "cannot parse sections")
	}
	sectionBoundaries := sectionNameRegex.FindAllStringSubmatchIndex(sectionText, -1)
	numSections := len(sectionBoundaries)

	// extract sections from sectionText
	for i, idxs := range sectionBoundaries {
		var curSecEnd int
		// end of current section is marked by the beginning of the next section
		// if current section is not the last section
		if i != numSections-1 {
			curSecEnd = sectionBoundaries[i+1][0]
		} else {
			curSecEnd = len(sectionText)
		}

		section, err := parseSection(sectionText[idxs[0]:curSecEnd], t.opts)
		if err != nil {
			return errors.Wrap(err, "failed to parse section while reading textnote")
		}

		idx, found := t.sectionIdx[section.name]
		if !found {
			return fmt.Errorf("cannot load undefined section [%s]", section.name)
		}
		t.sections[idx] = section
	}

	return nil
}

func (t *Template) string() string {
	str := t.makeHeader()
	for _, section := range t.sections {
		name := section.getNameString(t.opts.Section.Prefix, t.opts.Section.Suffix)
		body := section.getContentString()
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

func (t *Template) getSection(name string) (*section, error) {
	idx, found := t.sectionIdx[name]
	if !found {
		return &section{}, fmt.Errorf("section [%s] not found", name)
	}
	return t.sections[idx], nil
}

// ParseTemplateFileName extracts a time.Time from a file name and returns an additional
// bool indicating if name corresponds to a valid template file name
func ParseTemplateFileName(fileName string, opts config.FileOpts) (t time.Time, ok bool) {
	// ensure extension matches template file name convention
	ext := filepath.Ext(fileName)
	if ext == "." {
		return t, false
	}
	if strings.TrimPrefix(ext, ".") != opts.Ext {
		return t, false
	}

	baseName := strings.TrimSuffix(fileName, ext)
	t, err := time.Parse(opts.TimeFormat, baseName)
	if err != nil {
		return t, false
	}
	return t, true
}
