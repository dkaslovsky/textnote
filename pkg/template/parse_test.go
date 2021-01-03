package template

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/template/templatetest"
	"github.com/stretchr/testify/require"
)

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
				fmt.Sprintf("%s\nhello\n  world", templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts())),
				"\n",
			),
			expected: []contentItem{
				{
					header: templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					text:   "hello\n  world",
				},
			},
		},
		"lines with single header with newline at start and end": {
			lines: strings.Split(
				fmt.Sprintf("\n%s\n\nhello\n  world\n", templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts())),
				"\n",
			),
			expected: []contentItem{
				{},
				{
					header: templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					text:   "\nhello\n  world\n",
				},
			},
		},
		"lines with multiple headers": {
			lines: strings.Split(
				fmt.Sprintf("%s\nhello\n  world\n%s\nhello2\n  world2",
					templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
				),
				"\n",
			),
			expected: []contentItem{
				{
					header: templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					text:   "hello\n  world",
				},
				{
					header: templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
					text:   "hello2\n  world2",
				},
			},
		},
		"lines with multiple headers with newline at start and end": {
			lines: strings.Split(
				fmt.Sprintf("\n%s\nhello\n\n  world\n\n\n\n%s\nhello2\n  world2\n\n\n\n\n",
					templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
				),
				"\n",
			),
			expected: []contentItem{
				{},
				{
					header: templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					text:   "hello\n\n  world\n\n\n",
				},
				{
					header: templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
					text:   "hello2\n  world2\n\n\n\n\n",
				},
			},
		},
		"header with no text": {
			lines: strings.Split(
				templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
				"\n",
			),
			expected: []contentItem{
				{
					header: templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					text:   "",
				},
			},
		},
		"multiple headers with no text": {
			lines: strings.Split(
				fmt.Sprintf("%s\n%s",
					templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
				),
				"\n",
			),
			expected: []contentItem{
				{
					header: templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					text:   "",
				},
				{
					header: templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
					text:   "",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			contents := parseSectionContents(test.lines, templatetest.GetOpts().Archive.HeaderPrefix, templatetest.GetOpts().Archive.HeaderSuffix, templatetest.GetOpts().File.TimeFormat)
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
		"valid header with no prefix or suffix": {
			header:   "2020-07-28",
			prefix:   "",
			suffix:   "",
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

func TestIsEmptyContents(t *testing.T) {
	type testCase struct {
		contents []contentItem
		expected bool
	}

	tests := map[string]testCase{
		"empty contents": {
			contents: []contentItem{},
			expected: true,
		},
		"single content with only newlines and empty header": {
			contents: []contentItem{
				contentItem{
					header: "",
					text:   "\n\n\n",
				},
			},
			expected: true,
		},
		"single content with only newlines and populated header": {
			contents: []contentItem{
				contentItem{
					header: "header",
					text:   "\n\n\n",
				},
			},
			expected: true,
		},
		"multiple contents with only newlines and empty headers": {
			contents: []contentItem{
				contentItem{
					header: "",
					text:   "\n\n\n",
				},
				contentItem{
					header: "",
					text:   "\n",
				},
			},
			expected: true,
		},
		"multiple contents with only newlines and populated headers": {
			contents: []contentItem{
				contentItem{
					header: "header1",
					text:   "\n\n\n",
				},
				contentItem{
					header: "header2",
					text:   "\n",
				},
			},
			expected: true,
		},
		"single content with text and no header": {
			contents: []contentItem{
				contentItem{
					header: "",
					text:   "\n\nfoo\n",
				},
			},
			expected: false,
		},
		"single content with text and populated header": {
			contents: []contentItem{
				contentItem{
					header: "header",
					text:   "\n\nfoo\n",
				},
			},
			expected: false,
		},
		"multiple contents with text and no headers": {
			contents: []contentItem{
				contentItem{
					header: "",
					text:   "\n\nfoo\n",
				},
				contentItem{
					header: "",
					text:   "bar",
				},
			},
			expected: false,
		},
		"multiple contents with text and populated headers": {
			contents: []contentItem{
				contentItem{
					header: "header1",
					text:   "\n\nfoo\n",
				},
				contentItem{
					header: "header2",
					text:   "bar",
				},
			},
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			val := isEmptyContents(test.contents)
			require.Equal(t, test.expected, val)
		})
	}
}
