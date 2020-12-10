package template

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

// TODO: this is config!
const (
	fileExt = "txt"

	headerPrefix           = ""
	headerSuffix           = ""
	headerTrailingNewlines = 1

	sectionNamePrefix       = "___"
	sectionNameSuffix       = "___"
	sectionTrailingNewlines = 3

	// SectionNames are the names for each section
	SectionNames = "TODO,DONE,NOTES"

	timeFormatFileName = "2006-01-02"
	timeFormatHeader   = "[Mon] 02 Jan 2006"
)

var (
	// FirstSectionFirstLine is the first line of content of the first Section, used when opening with Vim
	FirstSectionFirstLine = headerTrailingNewlines + 3

	regexNewlines    = regexp.MustCompile(`\n{2,}`)
	regexSectionName = regexp.MustCompile(fmt.Sprintf("%s[A-Za-z]+%s", sectionNamePrefix, sectionNameSuffix))
)

// Body is the structure for a note
type Body struct {
	Date     time.Time
	Sections []*Section
}

// NewBody constructs a Body
func NewBody(date time.Time, sections ...*Section) Body {
	return Body{
		Date:     date,
		Sections: sections,
	}
}

func (b Body) String() string {
	s := b.makeHeader()
	for _, section := range b.Sections {
		s += section.String()
	}
	return s
}

func (b Body) Write(w io.Writer) error {
	_, err := w.Write([]byte(b.String()))
	return err
}

func (b Body) makeHeader() string {
	return fmt.Sprintf("%s%s%s\n%s",
		headerPrefix,
		b.Date.Format(timeFormatHeader),
		headerSuffix,
		strings.Repeat("\n", headerTrailingNewlines),
	)
}

// Section is a named section of a note
type Section struct {
	Name     string
	Contents []string // use a slice in case we want to treat contents as a list of bulleted items
}

// NewSection constructs a Section
func NewSection(name string, contents ...string) *Section {
	return &Section{
		Name:     name,
		Contents: contents,
	}
}

func (s *Section) String() string {
	str := fmt.Sprintf("%s%s%s\n", sectionNamePrefix, s.Name, sectionNameSuffix)
	for _, content := range s.Contents {
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		str += content
	}
	return str + strings.Repeat("\n", sectionTrailingNewlines)
}

// CleanNewlines mutates a section to remove all newlines
func (s *Section) CleanNewlines() {
	for i, content := range s.Contents {
		s.Contents[i] = strings.Trim(regexNewlines.ReplaceAllString(content, "\n"), "\n")
	}
}
