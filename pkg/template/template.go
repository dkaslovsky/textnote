package template

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

const (
	fileExt = "txt"

	headerTrailingNewlines = 1

	sectionNameWrapper      = "___"
	sectionTrailingNewlines = 3

	timeFormatFileName = "2006-01-02"
	timeFormatHeader   = "[Mon] 02 Jan 2006"
)

var (
	regexNewlines = regexp.MustCompile(`\n{2,}`)
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
	return fmt.Sprintf("%s\n%s",
		b.Date.Format(timeFormatHeader),
		whitespace("\n", headerTrailingNewlines),
	)
}

// GetFileName formats a time.Time object into a format used as a filename
func GetFileName(t time.Time) string {
	return fmt.Sprintf("%s.%s", t.Format(timeFormatFileName), fileExt)
}

// GetFirstSectionLine calculates the line number of the first Section for Vim to open on
func GetFirstSectionLine() int {
	// there are 2 lines (the header itself and the section title) plus the header whitespace
	// to get to the section title, so add one to start cursor inside the section
	return headerTrailingNewlines + 3
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
	str := fmt.Sprintf("%s%s%s\n", sectionNameWrapper, s.Name, sectionNameWrapper)
	for _, content := range s.Contents {
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		str += content
	}
	return str + whitespace("\n", sectionTrailingNewlines)
}

// CleanNewlines mutates a section to remove all newlines
func (s *Section) CleanNewlines() {
	for i, content := range s.Contents {
		s.Contents[i] = strings.Trim(regexNewlines.ReplaceAllString(content, "\n"), "\n")
	}
}
