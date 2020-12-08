package template

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	regexSectionName = regexp.MustCompile(fmt.Sprintf("%s[A-Z]+%s", sectionNameWrapper, sectionNameWrapper))
)

// Read parses a note into a Body
func Read(r io.Reader) (b Body, err error) {
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

		section, err := parseSection(text[start:end])
		if err != nil {
			return b, err
		}

		sections = append(sections, section)
	}

	return NewBody(date, sections...), nil
}

func parseSection(s string) (*Section, error) {
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

func whitespace(ch string, repeat int) string {
	return strings.Repeat(string(ch), repeat)
}
