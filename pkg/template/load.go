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
	sectionText := string(raw)

	// extract sections from sectionText
	sectionNameRegex, err := getSectionNameRegex(t.opts.Section.Prefix, t.opts.Section.Suffix)
	if err != nil {
		return errors.Wrap(err, "cannot parse sections")
	}
	matchIdx := sectionNameRegex.FindAllStringSubmatchIndex(sectionText, -1)
	for i, idx := range matchIdx {
		// get start and end indices for each section
		var start, end int
		start = idx[0]
		if i+1 < len(matchIdx) {
			end = matchIdx[i+1][0]
		} else {
			end = len(sectionText)
		}

		section, err := parseSection(sectionText[start:end], t.opts)
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

func parseSection(text string, opts config.Opts) (*section, error) {
	if len(text) == 0 {
		return nil, errors.New("cannot parse Section from empty input")
	}

	lines := strings.Split(text, "\n")
	name := stripPrefixSuffix(lines[0], opts.Section.Prefix, opts.Section.Suffix)
	contents := parseSectionContents(
		lines[1:],
		opts.Archive.SectionContentPrefix,
		opts.Archive.SectionContentSuffix,
		opts.File.TimeFormat,
	)

	// TODO - clean up this hack
	// do not include trailing newlines for empty sections as content
	if len(contents) == 1 && contents[0].text == strings.Repeat("\n", opts.Section.TrailingNewlines) {
		return newSection(name), nil
	}

	return newSection(name, contents...), nil
}

func parseSectionContents(lines []string, prefix string, suffix string, format string) []contentItem {
	contents := []contentItem{}
	if len(lines) == 0 {
		return contents
	}

	// parse first line
	line := lines[0]
	header := ""
	body := []string{}
	if isItemHeader(line, prefix, suffix, format) {
		header = line
	} else {
		body = append(body, line)
	}

	for _, line := range lines[1:] {
		if !isItemHeader(line, prefix, suffix, format) {
			body = append(body, line)
			continue
		}

		contents = append(contents, contentItem{
			header: header,
			text:   strings.Join(body, "\n"),
		})

		header = line
		body = []string{}
	}

	if len(body) != 0 || header != "" {
		contents = append(contents, contentItem{
			header: header,
			text:   strings.Join(body, "\n"),
		})
	}
	return contents
}

func stripPrefixSuffix(line string, prefix string, suffix string) string {
	return strings.TrimPrefix(strings.TrimSuffix(line, suffix), prefix)
}

func getSectionNameRegex(prefix string, suffix string) (*regexp.Regexp, error) {
	sectionPattern := fmt.Sprintf("%s[A-Za-z]+%s", prefix, suffix)
	sectionNameRegex, err := regexp.Compile(sectionPattern)
	if err != nil {
		return sectionNameRegex, fmt.Errorf("invalid section prefix [%s] or suffix [%s]", prefix, suffix)
	}
	return sectionNameRegex, nil
}

func isItemHeader(line string, prefix string, suffix string, format string) bool {
	if !strings.HasPrefix(line, prefix) {
		return false
	}
	if !strings.HasSuffix(line, suffix) {
		return false
	}
	_, err := time.Parse(format, stripPrefixSuffix(line, prefix, suffix))
	if err != nil {
		return false
	}
	return true
}
