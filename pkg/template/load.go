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
			end = matchIdx[i+1][0] - 1
		} else {
			end = len(sectionText) - 1
		}

		section, err := parseSection(sectionText[start:end], t.opts.Section)
		if err != nil {
			return errors.Wrap(err, "failed to parse Section while reading textnote")
		}

		idx, found := t.sectionIdx[section.name]
		if !found {
			return fmt.Errorf("cannot load undefined section [%s]", section.name)
		}
		t.sections[idx] = section
	}

	return nil
}

func parseSection(text string, opts config.SectionOpts) (*section, error) {
	if len(text) == 0 {
		return nil, errors.New("cannot parse Section from empty input")
	}

	lines := strings.Split(text, "\n")
	name := parseSectionName(lines[0], opts.Prefix, opts.Suffix)
	// if len(lines) == 1 {
	// 	return newSection(name), nil
	// }
	contents := parseSectionContents(lines[1:])
	return newSection(name, contents...), nil
}

func parseSectionContents(lines []string) []string {
	contents := []string{}
	curItems := []string{lines[0]}
	for _, line := range lines[1:] {
		// if line is not a continuation then reform and add as an element of contents
		if !strings.HasPrefix(line, " ") {
			contents = append(contents, strings.Join(curItems, "\n"))
			curItems = []string{}
		}
		curItems = append(curItems, line)
	}
	// ensure last set of items are appended
	if len(curItems) > 0 {
		contents = append(contents, strings.Join(curItems, "\n"))
	}
	return contents
}

func parseSectionName(line string, prefix string, suffix string) string {
	return strings.TrimPrefix(strings.TrimSuffix(line, suffix), prefix)
}
