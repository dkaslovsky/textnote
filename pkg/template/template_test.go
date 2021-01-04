package template

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dkaslovsky/TextNote/pkg/template/templatetest"
	"github.com/stretchr/testify/require"
)

func TestNewTemplate(t *testing.T) {
	type testCase struct {
		sections           []string
		expectedSections   []*section
		expectedSectionIdx map[string]int
	}

	tests := map[string]testCase{
		"no sections": {
			sections:           []string{},
			expectedSections:   []*section{},
			expectedSectionIdx: map[string]int{},
		},
		"single section": {
			sections: []string{
				"section1",
			},
			expectedSections: []*section{
				newSection("section1"),
			},
			expectedSectionIdx: map[string]int{
				"section1": 0,
			},
		},
		"multiple sections": {
			sections: []string{
				"section1",
				"section3",
				"section2",
			},
			expectedSections: []*section{
				newSection("section1"),
				newSection("section3"),
				newSection("section2"),
			},
			expectedSectionIdx: map[string]int{
				"section1": 0,
				"section2": 2,
				"section3": 1,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			opts := templatetest.GetOpts()
			opts.Section.Names = test.sections
			template := NewTemplate(opts, templatetest.Date)

			require.Equal(t, templatetest.Date, template.date)
			require.Equal(t, test.expectedSections, template.sections)
			require.Equal(t, test.expectedSectionIdx, template.sectionIdx)
		})
	}
}

func TestGetFilePath(t *testing.T) {
	t.Run("get file path with extension", func(t *testing.T) {
		opts := templatetest.GetOpts()
		opts.File.Ext = "txt"
		template := NewTemplate(opts, templatetest.Date)
		filePath := template.GetFilePath()
		require.True(t, strings.HasPrefix(filePath, opts.AppDir))
		require.True(t, strings.HasSuffix(filePath, ".txt"))
		require.Equal(t, templatetest.Date.Format(opts.File.TimeFormat), stripPrefixSuffix(filePath,
			fmt.Sprintf("%s/", opts.AppDir), ".txt"),
		)
	})

	t.Run("get file path without extension", func(t *testing.T) {
		opts := templatetest.GetOpts()
		opts.File.Ext = ""
		template := NewTemplate(opts, templatetest.Date)
		filePath := template.GetFilePath()
		require.True(t, strings.HasPrefix(filePath, opts.AppDir))
		require.False(t, strings.HasSuffix(filePath, "."))
		require.Equal(t, templatetest.Date.Format(opts.File.TimeFormat), stripPrefixSuffix(filePath,
			fmt.Sprintf("%s/", opts.AppDir), ""),
		)
	})
}

func TestCopySectionContents(t *testing.T) {
	type testCase struct {
		sectionName      string
		existingContents []contentItem
		incomingContents []contentItem
	}

	tests := map[string]testCase{
		"copy empty contents into empty section": {
			sectionName:      "TestSection1",
			existingContents: []contentItem{},
			incomingContents: []contentItem{},
		},
		"copy empty contents into populated section": {
			sectionName: "TestSection1",
			existingContents: []contentItem{
				contentItem{
					header: "existingHeader",
					text:   "existingText1",
				},
			},
			incomingContents: []contentItem{},
		},
		"copy contents with single element into empty section": {
			sectionName:      "TestSection1",
			existingContents: []contentItem{},
			incomingContents: []contentItem{
				contentItem{
					header: "header",
					text:   "text1",
				},
			},
		},
		"copy contents with single element into populated section": {
			sectionName: "TestSection1",
			existingContents: []contentItem{
				contentItem{
					header: "existingHeader",
					text:   "existingText1",
				},
			},
			incomingContents: []contentItem{
				contentItem{
					header: "header",
					text:   "text1",
				},
			},
		},
		"copy contents with multiple element into empty section": {
			sectionName:      "TestSection1",
			existingContents: []contentItem{},
			incomingContents: []contentItem{
				contentItem{
					header: "header1",
					text:   "text1",
				},
				contentItem{
					header: "header2",
					text:   "text2",
				},
			},
		},
		"copy contents with multiple elements into populated section": {
			sectionName: "TestSection1",
			existingContents: []contentItem{
				contentItem{
					header: "existingHeader",
					text:   "existingText1",
				},
			},
			incomingContents: []contentItem{
				contentItem{
					header: "header1",
					text:   "text1",
				},
				contentItem{
					header: "header2",
					text:   "text2",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			opts := templatetest.GetOpts()
			src := NewTemplate(opts, templatetest.Date)
			src.sections[src.sectionIdx[test.sectionName]].contents = test.incomingContents
			template := NewTemplate(opts, templatetest.Date)
			template.sections[template.sectionIdx[test.sectionName]].contents = test.existingContents

			err := template.CopySectionContents(src, test.sectionName)
			require.NoError(t, err)
			for _, content := range test.incomingContents {
				require.Contains(t, template.sections[template.sectionIdx[test.sectionName]].contents, content)
			}
		})
	}
}

func TestCopySectionContentsFail(t *testing.T) {
	t.Run("section does not exist in template", func(t *testing.T) {
		toCopy := "toBeCopied"
		opts := templatetest.GetOpts()
		template := NewTemplate(opts, templatetest.Date)
		src := NewTemplate(opts, templatetest.Date)
		src.sections = append(src.sections, newSection(toCopy))
		src.sectionIdx[toCopy] = len(src.sections) - 1

		err := template.CopySectionContents(src, toCopy)
		require.Error(t, err)
	})

	t.Run("section does not exist in source", func(t *testing.T) {
		toCopy := "toBeCopied"
		opts := templatetest.GetOpts()
		template := NewTemplate(opts, templatetest.Date)
		template.sections = append(template.sections, newSection(toCopy))
		template.sectionIdx[toCopy] = len(template.sections) - 1
		src := NewTemplate(opts, templatetest.Date)

		err := template.CopySectionContents(src, toCopy)
		require.Error(t, err)
	})
}

func TestDeleteSectionContents(t *testing.T) {
	t.Run("delete section with no contents", func(t *testing.T) {
		toDelete := "sectionToBeDeleted"
		template := NewTemplate(templatetest.GetOpts(), templatetest.Date)
		template.sections = append(template.sections, newSection(toDelete))
		template.sectionIdx[toDelete] = len(template.sections) - 1

		err := template.DeleteSectionContents(toDelete)
		require.NoError(t, err)
		require.Empty(t, template.sections[len(template.sections)-1].contents)
	})

	t.Run("delete section with contents", func(t *testing.T) {
		toDelete := "sectionToBeDeleted"
		template := NewTemplate(templatetest.GetOpts(), templatetest.Date)
		template.sections = append(template.sections, newSection(toDelete, contentItem{
			header: "header",
			text:   "text goes here",
		}))
		template.sectionIdx[toDelete] = len(template.sections) - 1

		err := template.DeleteSectionContents(toDelete)
		require.NoError(t, err)
		require.Empty(t, template.sections[len(template.sections)-1].contents)
	})

	t.Run("delete non-existent section", func(t *testing.T) {
		toDelete := "sectionToBeDeleted"
		opts := templatetest.GetOpts()
		template := NewTemplate(opts, templatetest.Date)

		err := template.DeleteSectionContents(toDelete)
		require.Error(t, err)
	})
}

func TestLoad(t *testing.T) {
	type testCase struct {
		text             string
		expectedSections []*section
	}

	tests := map[string]testCase{
		"empty text": {
			text: ``,
			expectedSections: []*section{
				newSection("TestSection1"),
				newSection("TestSection2"),
				newSection("TestSection3"),
			},
		},
		"no sections in text": {
			text: `-^-[Sun] 20 Dec 2020-v-

`,
			expectedSections: []*section{
				newSection("TestSection1"),
				newSection("TestSection2"),
				newSection("TestSection3"),
			},
		},
		"single empty section in text": {
			text: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_`,
			expectedSections: []*section{
				newSection("TestSection1"),
				newSection("TestSection2"),
				newSection("TestSection3"),
			},
		},
		"single empty section with trainling newlines in text": {
			text: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_



`,
			expectedSections: []*section{
				newSection("TestSection1"),
				newSection("TestSection2"),
				newSection("TestSection3"),
			},
		},
		"single empty section with too many trainling newlines in text": {
			text: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_





`,
			expectedSections: []*section{
				newSection("TestSection1"),
				newSection("TestSection2"),
				newSection("TestSection3"),
			},
		},
		"single empty second section in text": {
			text: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection2_q_`,
			expectedSections: []*section{
				newSection("TestSection1"),
				newSection("TestSection2"),
				newSection("TestSection3"),
			},
		},
		"multiple empty sections in text": {
			text: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_
_p_TestSection2_q_
_p_TestSection3_q_`,
			expectedSections: []*section{
				newSection("TestSection1"),
				newSection("TestSection2"),
				newSection("TestSection3"),
			},
		},
		"multiple empty sections with trailing newlines in text": {
			text: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_



_p_TestSection2_q_



_p_TestSection3_q_



`,
			expectedSections: []*section{
				newSection("TestSection1"),
				newSection("TestSection2"),
				newSection("TestSection3"),
			},
		},
		"single section with contents in text": {
			text: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_
text1
  text2


_p_TestSection2_q_
_p_TestSection3_q_
`,
			expectedSections: []*section{
				newSection("TestSection1",
					contentItem{
						header: "",
						text:   "text1\n  text2\n\n\n",
					},
				),
				newSection("TestSection2"),
				newSection("TestSection3"),
			},
		},
		"multiple sections with contents in text": {
			text: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_
text1
  text2


_p_TestSection2_q_
  text3
_p_TestSection3_q_

text4

`,
			expectedSections: []*section{
				newSection("TestSection1",
					contentItem{
						header: "",
						text:   "text1\n  text2\n\n\n",
					},
				),
				newSection("TestSection2",
					contentItem{
						header: "",
						text:   "  text3\n",
					},
				),
				newSection("TestSection3",
					contentItem{
						header: "",
						text:   "\ntext4\n\n",
					}),
			},
		},
		"section with single item header": {
			text: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_
[2020-12-18]
text1a
  text1b


_p_TestSection2_q_
_p_TestSection3_q_
`,
			expectedSections: []*section{
				newSection("TestSection1",
					contentItem{
						header: "[2020-12-18]",
						text:   "text1a\n  text1b\n\n\n",
					},
				),
				newSection("TestSection2"),
				newSection("TestSection3"),
			},
		},
		"section with multiple item headers": {
			text: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_
[2020-12-16]
text1a
  text1b


[2020-12-17]
text1c
[2020-12-18]
_p_TestSection2_q_
_p_TestSection3_q_
`,
			expectedSections: []*section{
				newSection("TestSection1",
					contentItem{
						header: "[2020-12-16]",
						text:   "text1a\n  text1b\n\n",
					},
					contentItem{
						header: "[2020-12-17]",
						text:   "text1c",
					},
					contentItem{
						header: "[2020-12-18]",
						text:   "",
					},
				),
				newSection("TestSection2"),
				newSection("TestSection3"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			template := NewTemplate(templatetest.GetOpts(), templatetest.Date)
			err := template.Load(strings.NewReader(test.text))
			require.NoError(t, err)
			for _, expectedSection := range test.expectedSections {
				sec, err := template.getSection(expectedSection.name)
				require.NoError(t, err)
				require.Equal(t, expectedSection, sec)
			}
		})
	}
}

func TestString(t *testing.T) {
	type testCase struct {
		sections []*section
		expected string
	}

	tests := map[string]testCase{
		"empty template": {
			sections: []*section{},
			expected: `-^-[Sun] 20 Dec 2020-v-

`,
		},
		"single empty section": {
			sections: []*section{
				newSection("TestSection1"),
			},
			expected: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_



`,
		},
		"single section with text": {
			sections: []*section{
				newSection("TestSection1",
					contentItem{
						header: "",
						text:   "text",
					},
				),
			},
			expected: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_
text
`,
		},
		"single section with multiline text": {
			sections: []*section{
				newSection("TestSection1",
					contentItem{
						header: "",
						text:   "text1\ntext2\n\n text3text4\n- text5\n\n  -text6\n\n",
					},
				),
			},
			expected: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_
text1
text2

 text3text4
- text5

  -text6

`,
		},
		"single section with text and header": {
			sections: []*section{
				newSection("TestSection1",
					contentItem{
						// in practice a Template will not have sections with headers
						// and as such we expect no formatting to be applied
						header: "header",
						text:   "text",
					},
				),
			},
			expected: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_
header
text
`,
		},
		"single section with multiple contents": {
			sections: []*section{
				newSection("TestSection1",
					contentItem{
						header: "",
						text:   "text1",
					},
					// in practice a Template will not have sections with multiple contents
					contentItem{
						header: "",
						text:   "text2",
					},
				),
			},
			expected: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_
text1
text2
`,
		},
		"multiple empty sections": {
			sections: []*section{
				newSection("TestSection1"),
				newSection("TestSection2"),
				newSection("TestSection3"),
			},
			expected: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_



_p_TestSection2_q_



_p_TestSection3_q_



`,
		},
		"multiple sections with only first populated": {
			sections: []*section{
				newSection("TestSection1",
					contentItem{
						header: "",
						text:   "text",
					},
				),
				newSection("TestSection2"),
				newSection("TestSection3"),
			},
			expected: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_
text
_p_TestSection2_q_



_p_TestSection3_q_



`,
		},
		"multiple sections with only middle populated": {
			sections: []*section{
				newSection("TestSection1"),
				newSection("TestSection2",
					contentItem{
						header: "",
						text:   "text",
					},
				),
				newSection("TestSection3"),
			},
			expected: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_



_p_TestSection2_q_
text
_p_TestSection3_q_



`,
		},
		"multiple sections with only last populated": {
			sections: []*section{
				newSection("TestSection1"),
				newSection("TestSection2"),
				newSection("TestSection3",
					contentItem{
						header: "",
						text:   "text",
					},
				),
			},
			expected: `-^-[Sun] 20 Dec 2020-v-

_p_TestSection1_q_



_p_TestSection2_q_



_p_TestSection3_q_
text
`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			opts := templatetest.GetOpts()
			names := []string{}
			for _, section := range test.sections {
				names = append(names, section.name)
			}
			opts.Section.Names = names

			template := NewTemplate(opts, templatetest.Date)
			for i, section := range test.sections {
				template.sections[i] = section
			}

			require.Equal(t, test.expected, template.string())
		})
	}
}
