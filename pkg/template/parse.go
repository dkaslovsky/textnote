package template

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/config"
	"github.com/pkg/errors"
)

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

	if isEmptyContents(contents) {
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
		if isItemHeader(line, prefix, suffix, format) {
			contents = append(contents, contentItem{
				header: header,
				text:   strings.Join(body, "\n"),
			})

			header = line
			body = []string{}
			continue
		}

		body = append(body, line)
	}

	// ensure remaining content is appended
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

func isEmptyContents(c []contentItem) bool {
	if len(c) == 0 {
		return true
	}
	if len(c) == 1 {
		// do not include trailing newlines as content for empty section
		strippedTxt := strings.Replace(c[0].text, "\n", "", -1)
		return len(strippedTxt) == 0
	}
	return false
}
