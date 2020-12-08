package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	headerTrailingNewlines = 1

	sectionNameWrapper      = "___"
	sectionTrailingNewlines = 3

	timeFormatFileName = "2006-01-02"
	timeFormatHeader   = "[Mon] 02 Jan 2006"
)

var (
	regexSectionName = regexp.MustCompile(`___[A-Z]+___`)
	regexNewlines    = regexp.MustCompile(`\n{2,}`)
)

type Body struct {
	Date     time.Time
	Sections []*Section
}

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

func (b Body) getFileName() string {
	return b.Date.Format(timeFormatFileName)
}

func (b Body) makeHeader() string {
	return fmt.Sprintf("%s\n%s",
		b.Date.Format(timeFormatHeader),
		whitespace("\n", headerTrailingNewlines),
	)
}

type Section struct {
	Name     string
	Contents []string // use a slice in case we want to treat contents as a list of bulleted items
}

func NewSection(name string, contents ...string) *Section {
	return &Section{
		Name:     name,
		Contents: contents,
	}
}

func ParseSection(s string) (*Section, error) {
	if len(s) == 0 {
		return nil, errors.New("cannot parse Section from empty input")
	}

	lines := strings.Split(strings.TrimSuffix(s, whitespace("\n", sectionTrailingNewlines)), "\n")
	name := strings.Trim(lines[0], sectionNameWrapper)
	if len(lines) == 1 {
		return NewSection(name), nil
	}

	contents := strings.TrimRight(strings.Join(lines[1:], "\n"), "\n")
	if contents == "" {
		return NewSection(name), nil
	}
	return NewSection(name, contents), nil
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

func (s *Section) CleanNewlines() {
	for i, content := range s.Contents {
		s.Contents[i] = strings.Trim(regexNewlines.ReplaceAllString(content, "\n"), "\n")
	}
}

func ParseBody(r io.Reader) (b Body, err error) {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return b, err
	}

	headerAndBody := strings.SplitN(string(raw), "\n", 2)
	if len(headerAndBody) != 2 {
		return b, errors.New("failed to parse header")
	}

	dateStr, text := headerAndBody[0], headerAndBody[1]
	date, err := time.Parse(timeFormatHeader, dateStr)
	if err != nil {
		return b, errors.Wrap(err, "failed to parse header")
	}

	sections := []*Section{}

	matchIdx := regexSectionName.FindAllStringSubmatchIndex(text, -1)
	for i, idx := range matchIdx {

		// get start and end indices for each section
		var start, end int
		start = idx[0]
		if i+1 < len(matchIdx) {
			end = matchIdx[i+1][0]
		} else {
			end = len(text) - 1
		}

		section, err := ParseSection(text[start:end])
		if err != nil {
			return b, err
		}

		sections = append(sections, section)
	}

	return NewBody(date, sections...), nil
}

func whitespace(ch string, repeat int) string {
	return strings.Repeat(string(ch), repeat)
}

func main() {
	date := time.Now()
	b := NewBody(date,
		NewSection("TODO",
			"- nothing",
			"- things\n  - all of them\n\n  - some of them",
			"- others\n\n\n\n",
			"- still more",
		),
		NewSection("DONE"),
		NewSection("NOTES", "- foo\n", "- bar"),
	)

	fmt.Println(b.getFileName())
	fmt.Println("-------------------")

	err := b.Write(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("-------------------")

	file := strings.NewReader(b.String())

	b2, err := ParseBody(file)
	if err != nil {
		log.Fatal(err)
	}

	err = b2.Write(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

// func ParseSection(s string) (*Section, error) {
// 	if len(s) == 0 {
// 		return nil, errors.New("some error goes here")
// 	}

// 	lines := strings.Split(s, "\n")
// 	name := strings.Trim(lines[0], sectionNameWrapper)
// 	if len(lines) == 1 {
// 		return NewSection(name), nil
// 	}

// 	contents := []string{}
// 	content := []string{}
// 	for _, line := range lines[1:] {
// 		// if line == "" {
// 		// 	continue
// 		// }
// 		if strings.HasPrefix(line, "-") {
// 			if len(content) > 0 {
// 				contents = append(contents, strings.Join(content, "\n"))
// 			}
// 			content = []string{}
// 		}
// 		content = append(content, line)
// 	}
// 	if len(content) > 0 {
// 		contents = append(contents, strings.Join(content, "\n"))
// 	}
// 	sec := NewSection(name, contents...)
// 	return sec, nil
// }
