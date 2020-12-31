package template

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/config"
	"github.com/stretchr/testify/require"
)

var testOpts = config.Opts{
	Header: config.HeaderOpts{
		Prefix:           "-^-",
		Suffix:           "-v-",
		TrailingNewlines: 1,
		TimeFormat:       "[Mon] 02 Jan 2006",
	},
	Section: config.SectionOpts{
		Prefix:           "_p_",
		Suffix:           "_q_",
		TrailingNewlines: 3,
		Names: []string{
			"TestSection1",
			"TestSection2",
			"TestSection3",
		},
	},
	File: config.FileOpts{
		TimeFormat: "2006-01-02",
	},
	Archive: config.ArchiveOpts{
		HeaderPrefix:             "ARCHIVEPREFIX ",
		HeaderSuffix:             " ARCHIVESUFFIX",
		SectionContentPrefix:     "[",
		SectionContentSuffix:     "]",
		SectionContentTimeFormat: "2006-01-02",
		MonthTimeFormat:          "Jan2006",
	},
}

var testDate = time.Date(2020, 12, 20, 1, 1, 1, 1, time.UTC)

func makeTestItemHeader(date time.Time, opts config.Opts) string {
	return fmt.Sprintf("%s%s%s",
		opts.Archive.HeaderPrefix,
		date.Format(opts.Archive.SectionContentTimeFormat),
		opts.Archive.HeaderSuffix,
	)
}

func TestParseSectionContents(t *testing.T) {
	type testCase struct {
		lines    []string
		expected []contentItem
	}

	tests := map[string]testCase{
		"empty lines": {
			lines:    []string{},
			expected: []contentItem{},
		},
		"single empty string line": {
			lines: []string{""},
			expected: []contentItem{
				{},
			},
		},
		"lines with no header": {
			lines: strings.Split("hello\n  world", "\n"),
			expected: []contentItem{
				{
					header: "",
					text:   "hello\n  world",
				},
			},
		},
		"lines with no header with newline at start and end": {
			lines: strings.Split("\n\nhello\n  world\n\n", "\n"),
			expected: []contentItem{
				{
					header: "",
					text:   "\n\nhello\n  world\n\n",
				},
			},
		},
		"lines with single header": {
			lines: strings.Split(
				fmt.Sprintf("%s\nhello\n  world", makeTestItemHeader(testDate, testOpts)),
				"\n",
			),
			expected: []contentItem{
				{
					header: makeTestItemHeader(testDate, testOpts),
					text:   "hello\n  world",
				},
			},
		},
		"lines with single header with newline at start and end": {
			lines: strings.Split(
				fmt.Sprintf("\n%s\n\nhello\n  world\n", makeTestItemHeader(testDate, testOpts)),
				"\n",
			),
			expected: []contentItem{
				{},
				{
					header: makeTestItemHeader(testDate, testOpts),
					text:   "\nhello\n  world\n",
				},
			},
		},
		"lines with multiple headers": {
			lines: strings.Split(
				fmt.Sprintf("%s\nhello\n  world\n%s\nhello2\n  world2",
					makeTestItemHeader(testDate, testOpts),
					makeTestItemHeader(testDate.Add(24*time.Hour), testOpts),
				),
				"\n",
			),
			expected: []contentItem{
				{
					header: makeTestItemHeader(testDate, testOpts),
					text:   "hello\n  world",
				},
				{
					header: makeTestItemHeader(testDate.Add(24*time.Hour), testOpts),
					text:   "hello2\n  world2",
				},
			},
		},
		"lines with multiple headers with newline at start and end": {
			lines: strings.Split(
				fmt.Sprintf("\n%s\nhello\n  world\n\n%s\nhello2\n  world2\n",
					makeTestItemHeader(testDate, testOpts),
					makeTestItemHeader(testDate.Add(24*time.Hour), testOpts),
				),
				"\n",
			),
			expected: []contentItem{
				{},
				{
					header: makeTestItemHeader(testDate, testOpts),
					text:   "hello\n  world\n",
				},
				{
					header: makeTestItemHeader(testDate.Add(24*time.Hour), testOpts),
					text:   "hello2\n  world2\n",
				},
			},
		},
		"header with no text": {
			lines: strings.Split(
				makeTestItemHeader(testDate, testOpts),
				"\n",
			),
			expected: []contentItem{
				{
					header: makeTestItemHeader(testDate, testOpts),
					text:   "",
				},
			},
		},
		"multiple headers with no text": {
			lines: strings.Split(
				fmt.Sprintf("%s\n%s",
					makeTestItemHeader(testDate, testOpts),
					makeTestItemHeader(testDate.Add(24*time.Hour), testOpts),
				),
				"\n",
			),
			expected: []contentItem{
				{
					header: makeTestItemHeader(testDate, testOpts),
					text:   "",
				},
				{
					header: makeTestItemHeader(testDate.Add(24*time.Hour), testOpts),
					text:   "",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			contents := parseSectionContents(test.lines, testOpts.Archive.HeaderPrefix, testOpts.Archive.HeaderSuffix, testOpts.File.TimeFormat)
			require.Equal(t, test.expected, contents)
		})
	}
}

func TestIsItemHeader(t *testing.T) {
	type testCase struct {
		header   string
		prefix   string
		suffix   string
		format   string
		expected bool
	}

	tests := map[string]testCase{
		"valid header": {
			header:   "[2020-07-28]",
			prefix:   "[",
			suffix:   "]",
			format:   "2006-01-02",
			expected: true,
		},
		"invalid header with wrong prefix": {
			header:   "<2020-07-28]",
			prefix:   "[",
			suffix:   "]",
			format:   "2006-01-02",
			expected: false,
		},
		"invalid header with wrong suffix": {
			header:   "[2020-07-28>",
			prefix:   "[",
			suffix:   "]",
			format:   "2006-01-02",
			expected: false,
		},
		"invalid header with wrong format": {
			header:   "[2020-July-28]",
			prefix:   "[",
			suffix:   "]",
			format:   "2006-01-02",
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			val := isItemHeader(test.header, test.prefix, test.suffix, test.format)
			require.Equal(t, test.expected, val)
		})
	}
}
