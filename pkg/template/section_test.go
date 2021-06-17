package template

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dkaslovsky/textnote/pkg/template/templatetest"
	"github.com/stretchr/testify/require"
)

func TestGetNameString(t *testing.T) {
	type testCase struct {
		name     string
		prefix   string
		suffix   string
		expected string
	}

	tests := map[string]testCase{
		"empty name, empty prefix and suffix": {
			name:     "",
			prefix:   "",
			suffix:   "",
			expected: "\n",
		},
		"empty name, non-empty prefix and suffix": {
			name:     "",
			prefix:   "p ",
			suffix:   " s",
			expected: "p  s\n",
		},
		"non-empty name, empty prefix and suffix": {
			name:     "name",
			prefix:   "",
			suffix:   "",
			expected: "name\n",
		},
		"non-empty name, non-empty prefix and suffix": {
			name:     "name",
			prefix:   "p ",
			suffix:   " s",
			expected: "p name s\n",
		},
		"non-empty name with spaces, non-empty prefix and suffix": {
			name:     " na me ",
			prefix:   "p ",
			suffix:   " s",
			expected: "p  na me  s\n",
		},
		"non-empty name with newlines, non-empty prefix and suffix": {
			name:     " na \n me ",
			prefix:   "p ",
			suffix:   " s",
			expected: "p  na \n me  s\n",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := newSection(test.name)
			str := s.getNameString(test.prefix, test.suffix)
			require.Equal(t, test.expected, str)
		})
	}
}

func TestGetContentString(t *testing.T) {
	type testCase struct {
		contents []contentItem
		expected string
	}

	tests := map[string]testCase{
		"empty contents": {
			contents: []contentItem{},
			expected: "",
		},
		"single empty string contents with no header": {
			contents: []contentItem{
				{
					header: "",
					text:   "",
				},
			},
			expected: "\n",
		},
		"single empty string contents with header": {
			contents: []contentItem{
				{
					header: "header",
					text:   "",
				},
			},
			expected: "header\n",
		},
		"multiple empty string contents with no header": {
			contents: []contentItem{
				{
					header: "",
					text:   "",
				},
				{
					header: "",
					text:   "",
				},
			},
			expected: "\n\n",
		},
		"multiple empty string contents with header": {
			contents: []contentItem{
				{
					header: "header1",
					text:   "",
				},
				{
					header: "header2",
					text:   "",
				},
			},
			expected: "header1\nheader2\n",
		},
		"single nonempty contents with no header missing trailing newline": {
			contents: []contentItem{
				{
					header: "",
					text:   "text\n goes\n  here",
				},
			},
			expected: "text\n goes\n  here\n",
		},
		"single nonempty contents with no header": {
			contents: []contentItem{
				{
					header: "",
					text:   "text\n goes\n  here\n",
				},
			},
			expected: "text\n goes\n  here\n",
		},
		"single nonempty contents with header": {
			contents: []contentItem{
				{
					header: "header",
					text:   "text\n goes\n  here\n",
				},
			},
			expected: "header\ntext\n goes\n  here\n",
		},
		"multiple nonempty contents with no headers": {
			contents: []contentItem{
				{
					header: "",
					text:   "text\n goes\n  here\n",
				},
				{
					header: "",
					text:   "text2\n goes2\n  here2 \n",
				},
			},
			expected: "text\n goes\n  here\ntext2\n goes2\n  here2 \n",
		},
		"multiple nonempty contents with headers": {
			contents: []contentItem{
				{
					header: "header1 ",
					text:   "text\n goes\n  here\n",
				},
				{
					header: " header2",
					text:   "text2\n goes2\n  here2 \n",
				},
			},
			expected: "header1 \ntext\n goes\n  here\n header2\ntext2\n goes2\n  here2 \n",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := newSection("name", test.contents...)
			str := s.getContentString()
			require.Equal(t, test.expected, str)
		})
	}
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
			opts := templatetest.GetOpts()
			contents := parseSectionContents(test.lines, opts.Archive.SectionContentPrefix, opts.Archive.SectionContentSuffix, opts.File.TimeFormat)
			require.Equal(t, test.expected, contents)
		})
	}
}

// func TestIsEmptyContents(t *testing.T) {
// 	type testCase struct {
// 		contents []contentItem
// 		expected bool
// 	}

// 	tests := map[string]testCase{
// 		"empty contents": {
// 			contents: []contentItem{},
// 			expected: true,
// 		},
// 		"single content with only newlines and empty header": {
// 			contents: []contentItem{
// 				{
// 					header: "",
// 					text:   "\n\n\n",
// 				},
// 			},
// 			expected: true,
// 		},
// 		"single content with only newlines and populated header": {
// 			contents: []contentItem{
// 				{
// 					header: "header",
// 					text:   "\n\n\n",
// 				},
// 			},
// 			expected: true,
// 		},
// 		"multiple contents with only newlines and empty headers": {
// 			contents: []contentItem{
// 				{
// 					header: "",
// 					text:   "\n\n\n",
// 				},
// 				{
// 					header: "",
// 					text:   "\n",
// 				},
// 			},
// 			expected: true,
// 		},
// 		"multiple contents with only newlines and populated headers": {
// 			contents: []contentItem{
// 				{
// 					header: "header1",
// 					text:   "\n\n\n",
// 				},
// 				{
// 					header: "header2",
// 					text:   "\n",
// 				},
// 			},
// 			expected: true,
// 		},
// 		"single content with text and no header": {
// 			contents: []contentItem{
// 				{
// 					header: "",
// 					text:   "\n\nfoo\n",
// 				},
// 			},
// 			expected: false,
// 		},
// 		"single content with text and populated header": {
// 			contents: []contentItem{
// 				{
// 					header: "header",
// 					text:   "\n\nfoo\n",
// 				},
// 			},
// 			expected: false,
// 		},
// 		"multiple contents with text and no headers": {
// 			contents: []contentItem{
// 				{
// 					header: "",
// 					text:   "\n\nfoo\n",
// 				},
// 				{
// 					header: "",
// 					text:   "bar",
// 				},
// 			},
// 			expected: false,
// 		},
// 		"multiple contents with text and populated headers": {
// 			contents: []contentItem{
// 				{
// 					header: "header1",
// 					text:   "\n\nfoo\n",
// 				},
// 				{
// 					header: "header2",
// 					text:   "bar",
// 				},
// 			},
// 			expected: false,
// 		},
// 	}

// 	for name, test := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			val := isEmptyContents(test.contents)
// 			require.Equal(t, test.expected, val)
// 		})
// 	}
// }

func TestIsEmpty(t *testing.T) {
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
				{
					header: "",
					text:   "\n\n\n",
				},
			},
			expected: true,
		},
		"single content with only newlines and populated header": {
			contents: []contentItem{
				{
					header: "header",
					text:   "\n\n\n",
				},
			},
			expected: true,
		},
		"multiple contents with only newlines and empty headers": {
			contents: []contentItem{
				{
					header: "",
					text:   "\n\n\n",
				},
				{
					header: "",
					text:   "\n",
				},
			},
			expected: true,
		},
		"multiple contents with only newlines and populated headers": {
			contents: []contentItem{
				{
					header: "header1",
					text:   "\n\n\n",
				},
				{
					header: "header2",
					text:   "\n",
				},
			},
			expected: true,
		},
		"single content with text and no header": {
			contents: []contentItem{
				{
					header: "",
					text:   "\n\nfoo\n",
				},
			},
			expected: false,
		},
		"single content with text and populated header": {
			contents: []contentItem{
				{
					header: "header",
					text:   "\n\nfoo\n",
				},
			},
			expected: false,
		},
		"multiple contents with text and no headers": {
			contents: []contentItem{
				{
					header: "",
					text:   "\n\nfoo\n",
				},
				{
					header: "",
					text:   "bar",
				},
			},
			expected: false,
		},
		"multiple contents with text and populated headers": {
			contents: []contentItem{
				{
					header: "header1",
					text:   "\n\nfoo\n",
				},
				{
					header: "header2",
					text:   "bar",
				},
			},
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := newSection("name", test.contents...)
			val := s.isEmpty()
			require.Equal(t, test.expected, val)
		})
	}
}

func TestContentItemIsEmpty(t *testing.T) {
	type testCase struct {
		item     contentItem
		expected bool
	}

	tests := map[string]testCase{
		"empty": {
			item:     contentItem{},
			expected: true,
		},
		"only newlines and empty header": {
			item: contentItem{
				header: "",
				text:   "\n\n\n",
			},
			expected: true,
		},
		"only newlines and populated header": {
			item: contentItem{
				header: "header",
				text:   "\n\n\n",
			},
			expected: true,
		},
		"text and no header": {
			item: contentItem{
				header: "",
				text:   "\n\nfoo\n",
			},
			expected: false,
		},
		"text and populated header": {
			item: contentItem{
				header: "header",
				text:   "\n\nfoo\n",
			},
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			val := test.item.isEmpty()
			require.Equal(t, test.expected, val)
		})
	}

}
