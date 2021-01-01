package template

import (
	"testing"

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
				contentItem{
					header: "",
					text:   "",
				},
			},
			expected: "\n",
		},
		"single empty string contents with header": {
			contents: []contentItem{
				contentItem{
					header: "header",
					text:   "",
				},
			},
			expected: "header\n",
		},
		"multiple empty string contents with no header": {
			contents: []contentItem{
				contentItem{
					header: "",
					text:   "",
				},
				contentItem{
					header: "",
					text:   "",
				},
			},
			expected: "\n\n",
		},
		"multiple empty string contents with header": {
			contents: []contentItem{
				contentItem{
					header: "header1",
					text:   "",
				},
				contentItem{
					header: "header2",
					text:   "",
				},
			},
			expected: "header1\nheader2\n",
		},
		"single nonempty contents with no header missing trailing newline": {
			contents: []contentItem{
				contentItem{
					header: "",
					text:   "text\n goes\n  here",
				},
			},
			expected: "text\n goes\n  here\n",
		},
		"single nonempty contents with no header": {
			contents: []contentItem{
				contentItem{
					header: "",
					text:   "text\n goes\n  here\n",
				},
			},
			expected: "text\n goes\n  here\n",
		},
		"single nonempty contents with header": {
			contents: []contentItem{
				contentItem{
					header: "header",
					text:   "text\n goes\n  here\n",
				},
			},
			expected: "header\ntext\n goes\n  here\n",
		},
		"multiple nonempty contents with no headers": {
			contents: []contentItem{
				contentItem{
					header: "",
					text:   "text\n goes\n  here\n",
				},
				contentItem{
					header: "",
					text:   "text2\n goes2\n  here2 \n",
				},
			},
			expected: "text\n goes\n  here\ntext2\n goes2\n  here2 \n",
		},
		"multiple nonempty contents with headers": {
			contents: []contentItem{
				contentItem{
					header: "header1 ",
					text:   "text\n goes\n  here\n",
				},
				contentItem{
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
