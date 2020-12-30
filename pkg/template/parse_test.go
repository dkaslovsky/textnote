package template

import (
	"strings"
	"testing"

	"github.com/dkaslovsky/TextNote/pkg/config"
	"github.com/stretchr/testify/require"
)

var testOpts = config.Opts{
	Header: config.HeaderOpts{
		Prefix:           "--",
		Suffix:           "--",
		TrailingNewlines: 1,
		TimeFormat:       "[Mon] 02 Jan 2006",
	},
	Section: config.SectionOpts{
		Prefix:           "___",
		Suffix:           "___",
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
		HeaderPrefix:             "ARCHIVE ",
		HeaderSuffix:             " ARCHIVE",
		SectionContentPrefix:     "[",
		SectionContentSuffix:     "]",
		SectionContentTimeFormat: "2006-01-02",
		MonthTimeFormat:          "Jan2006",
	},
}

func TestParseSectionContents(t *testing.T) {
	type testCase struct {
		body     string
		expected []contentItem
	}

	tests := map[string]testCase{
		"non-archive": {
			body: "\n\nhello\n  world\n\n",
			expected: []contentItem{
				{
					header: "",
					text:   "\n\nhello\n  world\n\n",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			lines := strings.Split(test.body, "\n")
			contents := parseSectionContents(lines, testOpts.Section.Prefix, testOpts.Section.Suffix, testOpts.File.TimeFormat)
			require.Equal(t, test.expected, contents)
		})
	}
}
