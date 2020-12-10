package template

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Read parses a textnote into a Body
func Read(r io.Reader) (b Body, err error) {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return b, err
	}

	// separate header from text
	split := strings.SplitN(string(raw), "\n", 2)
	if len(split) != 2 {
		return b, errors.New("failed to parse header while reading textnote")
	}
	dateStr, text := split[0], split[1]
	date, err := time.Parse(timeFormatHeader, dateStr)
	if err != nil {
		return b, errors.Wrap(err, "failed to parse header while reading textnote")
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

		section, err := parseSection(text[start:end])
		if err != nil {
			return b, errors.Wrap(err, "failed to parse Section while reading textnote")
		}

		sections = append(sections, section)
	}

	return NewBody(date, sections...), nil
}

// GetFileName formats a time.Time object into a format used as a filename
func GetFileName(t time.Time) string {
	return fmt.Sprintf("%s.%s", t.Format(timeFormatFileName), fileExt)
}

func parseSection(s string) (*Section, error) {
	if len(s) == 0 {
		return nil, errors.New("cannot parse Section from empty input")
	}

	// trim off trailing newlines
	s = strings.TrimSuffix(s, strings.Repeat("\n", sectionTrailingNewlines))
	// split into lines
	lines := strings.Split(s, "\n")
	// extract name from first line
	name := parseSectionName(lines[0])
	if len(lines) == 1 {
		return NewSection(name), nil
	}

	// reform without the first line
	contents := strings.Join(lines[1:], "\n")
	//contents = strings.TrimRight(contents, "\n")
	if contents == "" {
		return NewSection(name), nil
	}
	return NewSection(name, contents), nil
}

func parseSectionName(line string) string {
	// trim the prefix and suffix to get the section name
	return strings.TrimPrefix(strings.TrimSuffix(line, sectionNameSuffix), sectionNamePrefix)
}
