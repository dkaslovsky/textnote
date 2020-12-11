package template

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/config"
	"github.com/pkg/errors"
)

const fileExt = "txt"

// Load populates a Template from the contents of a TextNote
func (t *Template) Load(r io.Reader) error {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "error loading template")
	}

	// extract header and section text
	split := strings.SplitN(string(raw), "\n", 2)
	if len(split) != 2 {
		return errors.New("failed to parse header while loading textnote")
	}
	dateStr, sectionText := split[0], split[1]

	// parse date from header
	date, err := time.Parse(t.opts.Header.TimeFormat, dateStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse header while loading textnote")
	}
	t.SetDate(date)

	// extract sections from sectionText
	sectionPattern := fmt.Sprintf("%s[A-Za-z]+%s", t.opts.Section.Prefix, t.opts.Section.Suffix)
	sectionNameRegex, err := regexp.Compile(sectionPattern)
	if err != nil {
		return errors.Wrap(err, "failed to parse sections while loading textnote")
	}

	matchIdx := sectionNameRegex.FindAllStringSubmatchIndex(sectionText, -1)
	for i, idx := range matchIdx {
		// get start and end indices for each section
		var start, end int
		start = idx[0]
		if i+1 < len(matchIdx) {
			end = matchIdx[i+1][0]
		} else {
			end = len(sectionText) - 1
		}

		section, err := parseSectionText(sectionText[start:end], t.opts.Section)
		if err != nil {
			return errors.Wrap(err, "failed to parse Section while reading textnote")
		}
		t.AddSection(section)
	}

	return nil
}

func parseSectionText(text string, opts config.SectionOpts) (*section, error) {
	if len(text) == 0 {
		return nil, errors.New("cannot parse Section from empty input")
	}

	// trim off trailing newlines
	text = strings.TrimSuffix(text, strings.Repeat("\n", opts.TrailingNewlines))
	// split into lines
	lines := strings.Split(text, "\n")
	// extract name from first line
	name := parseSectionName(lines[0], opts.Prefix, opts.Suffix)
	if len(lines) == 1 {
		return newSection(name), nil
	}

	// reform without the first line
	contents := strings.Join(lines[1:], "\n")
	if contents == "" {
		return newSection(name), nil
	}
	return newSection(name, contents), nil
}

// remove prefix and suffix from a line to get the section name
func parseSectionName(line string, prefix string, suffix string) string {
	return strings.TrimPrefix(strings.TrimSuffix(line, suffix), prefix)
}
