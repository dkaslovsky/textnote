package template

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/template/templatetest"
	"github.com/stretchr/testify/require"
)

func TestNewMonthArchiveTemplate(t *testing.T) {
	type testCase struct {
		date     time.Time
		expected time.Time
	}

	tests := map[string]testCase{
		"first of the month": {
			date:     time.Date(2020, 12, 1, 2, 3, 4, 5, time.UTC),
			expected: time.Date(2020, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		"not first of the month": {
			date:     time.Date(2020, 12, 15, 2, 3, 4, 5, time.UTC),
			expected: time.Date(2020, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		"non UTC location": {
			date:     time.Date(2020, 12, 15, 2, 3, 4, 5, time.FixedZone("UTC-8", -8*60*60)),
			expected: time.Date(2020, 12, 1, 0, 0, 0, 0, time.FixedZone("UTC-8", -8*60*60)),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := NewMonthArchiveTemplate(templatetest.GetOpts(), test.date)
			require.Equal(t, test.expected, m.date)
		})
	}
}

func TestArchiveGetFilePath(t *testing.T) {
	t.Run("get file path with extension", func(t *testing.T) {
		opts := templatetest.GetOpts()
		opts.File.Ext = "txt"
		template := NewMonthArchiveTemplate(opts, templatetest.Date)
		filePath := template.GetFilePath()
		require.True(t, strings.HasPrefix(filePath, opts.AppDir))
		require.True(t, strings.HasSuffix(filePath, ".txt"))
		require.Equal(t,
			opts.Archive.FilePrefix+templatetest.Date.Format(opts.Archive.MonthTimeFormat),
			stripPrefixSuffix(filePath, fmt.Sprintf("%s/", opts.AppDir), ".txt"),
		)
	})

	t.Run("get file path without extension", func(t *testing.T) {
		opts := templatetest.GetOpts()
		opts.File.Ext = ""
		template := NewMonthArchiveTemplate(opts, templatetest.Date)
		filePath := template.GetFilePath()
		require.True(t, strings.HasPrefix(filePath, opts.AppDir))
		require.False(t, strings.HasSuffix(filePath, "."))
		require.Equal(t,
			opts.Archive.FilePrefix+templatetest.Date.Format(opts.Archive.MonthTimeFormat),
			stripPrefixSuffix(filePath, fmt.Sprintf("%s/", opts.AppDir), ""),
		)
	})
}

func TestArchiveSectionContents(t *testing.T) {
	type testCase struct {
		sectionName      string
		existingContents []contentItem
		sourceDate       time.Time
		sourceContents   []contentItem
		expectedContents []contentItem
	}

	tests := map[string]testCase{
		"archive empty contents into empty section": {
			sectionName:      "TestSection1",
			existingContents: []contentItem{},
			sourceDate:       templatetest.Date,
			sourceContents:   []contentItem{},
			expectedContents: []contentItem{},
		},
		"archive empty contents into populated section": {
			sectionName: "TestSection1",
			existingContents: []contentItem{
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					text:   "existingText1",
				},
			},
			sourceDate:     templatetest.Date.Add(24 * time.Hour),
			sourceContents: []contentItem{},
			expectedContents: []contentItem{
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					text:   "existingText1",
				},
			},
		},
		"archive contents with single element into empty section": {
			sectionName:      "TestSection1",
			existingContents: []contentItem{},
			sourceDate:       templatetest.Date.Add(24 * time.Hour),
			sourceContents: []contentItem{
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
					text:   "text1",
				},
			},
			expectedContents: []contentItem{
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
					text:   "text1",
				},
			},
		},
		"archive contents with single element into populated section": {
			sectionName: "TestSection1",
			existingContents: []contentItem{
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					text:   "existingText1",
				},
			},
			sourceDate: templatetest.Date.Add(24 * time.Hour),
			sourceContents: []contentItem{
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
					text:   "sourceText1",
				},
			},
			expectedContents: []contentItem{
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					text:   "existingText1",
				},
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
					text:   "sourceText1",
				},
			},
		},
		"archive contents with multiple element into empty section": {
			sectionName:      "TestSection1",
			existingContents: []contentItem{},
			sourceDate:       templatetest.Date.Add(24 * time.Hour),
			sourceContents: []contentItem{
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
					text:   "text1\n",
				},
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
					text:   "text2\n\n",
				},
			},
			expectedContents: []contentItem{
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
					text:   "text1\ntext2\n\n",
				},
			},
		},
		"archive contents with multiple elements into populated section": {
			sectionName: "TestSection1",
			existingContents: []contentItem{
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date.Add(-24*time.Hour), templatetest.GetOpts()),
					text:   "existingText",
				},
			},
			sourceDate: templatetest.Date.Add(24 * time.Hour),
			sourceContents: []contentItem{
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
					text:   "text1\n",
				},
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
					text:   "text2\n\n",
				},
			},
			expectedContents: []contentItem{
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date.Add(-24*time.Hour), templatetest.GetOpts()),
					text:   "existingText",
				},
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
					text:   "text1\ntext2\n\n",
				},
			},
		},
		"archive contents from source with same date": {
			sectionName: "TestSection1",
			existingContents: []contentItem{
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					text:   "existingText",
				},
			},
			sourceDate: templatetest.Date,
			sourceContents: []contentItem{
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					text:   "text1\n",
				},
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					text:   "text2\n\n",
				},
			},
			expectedContents: []contentItem{
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					text:   "existingText",
				},
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date, templatetest.GetOpts()),
					text:   "text1\ntext2\n\n",
				},
			},
		},
		"source header does not matter": {
			sectionName:      "TestSection1",
			existingContents: []contentItem{},
			sourceDate:       templatetest.Date.Add(24 * time.Hour),
			sourceContents: []contentItem{
				contentItem{
					header: "doesn't matter 1",
					text:   "text1\n",
				},
				contentItem{
					header: "doesn't matter 2",
					text:   "text2\n\n",
				},
			},
			expectedContents: []contentItem{
				contentItem{
					header: templatetest.MakeItemHeader(templatetest.Date.Add(24*time.Hour), templatetest.GetOpts()),
					text:   "text1\ntext2\n\n",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			opts := templatetest.GetOpts()
			src := NewTemplate(opts, test.sourceDate)
			src.sections[src.sectionIdx[test.sectionName]].contents = test.sourceContents
			template := NewMonthArchiveTemplate(opts, templatetest.Date)
			template.sections[template.sectionIdx[test.sectionName]].contents = test.existingContents

			err := template.ArchiveSectionContents(src, test.sectionName)
			require.NoError(t, err)
			require.Equal(t, template.sections[template.sectionIdx[test.sectionName]].contents, test.expectedContents)
		})
	}
}

func TestArchiveSectionContentsFail(t *testing.T) {
	t.Run("section does not exist in template", func(t *testing.T) {
		toCopy := "toBeArchived"
		opts := templatetest.GetOpts()
		template := NewMonthArchiveTemplate(opts, templatetest.Date)
		src := NewTemplate(opts, templatetest.Date)
		src.sections = append(src.sections, newSection(toCopy))
		src.sectionIdx[toCopy] = len(src.sections) - 1

		err := template.ArchiveSectionContents(src, toCopy)
		require.Error(t, err)
	})

	t.Run("section does not exist in source", func(t *testing.T) {
		toCopy := "toBeArchived"
		opts := templatetest.GetOpts()
		template := NewMonthArchiveTemplate(opts, templatetest.Date)
		template.sections = append(template.sections, newSection(toCopy))
		template.sectionIdx[toCopy] = len(template.sections) - 1
		src := NewTemplate(opts, templatetest.Date)

		err := template.ArchiveSectionContents(src, toCopy)
		require.Error(t, err)
	})
}

func TestArchiveString(t *testing.T) {
	type testCase struct {
		sections []*section
		expected string
	}

	tests := map[string]testCase{
		"empty template": {
			sections: []*section{},
			expected: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

`,
		},
		"single empty section": {
			sections: []*section{
				newSection("TestSection1"),
			},
			expected: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_



`,
		},
		"single section": {
			sections: []*section{
				newSection("TestSection1",
					contentItem{
						header: "[2020-12-19]",
						text:   "text",
					},
				),
			},
			expected: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-19]
text



`,
		},
		"single section with multiline text": {
			sections: []*section{
				newSection("TestSection1",
					contentItem{
						header: "[2020-12-19]",
						text:   "text1\ntext2\n\n text3text4\n- text5\n\n  -text6\n\n",
					},
				),
			},
			expected: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-19]
text1
text2
 text3text4
- text5
  -text6



`,
		},
		"single section with multiple contents": {
			sections: []*section{
				newSection("TestSection1",
					contentItem{
						header: "[2020-12-18]",
						text:   "text1\n",
					},
					contentItem{
						header: "[2020-12-19]",
						text:   "text2\n",
					},
				),
			},
			expected: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-18]
text1
[2020-12-19]
text2



`,
		},
		"multiple empty sections": {
			sections: []*section{
				newSection("TestSection1"),
				newSection("TestSection2"),
				newSection("TestSection3"),
			},
			expected: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_



_p_TestSection2_q_



_p_TestSection3_q_



`,
		},
		"multiple sections with only first populated": {
			sections: []*section{
				newSection("TestSection1",
					contentItem{
						header: "[2020-12-18]",
						text:   "text",
					},
				),
				newSection("TestSection2"),
				newSection("TestSection3"),
			},
			expected: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-18]
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
						header: "[2020-12-18]",
						text:   "text",
					},
				),
				newSection("TestSection3"),
			},
			expected: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_



_p_TestSection2_q_
[2020-12-18]
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
						header: "[2020-12-18]",
						text:   "text",
					},
				),
			},
			expected: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_



_p_TestSection2_q_



_p_TestSection3_q_
[2020-12-18]
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

			template := NewMonthArchiveTemplate(opts, templatetest.Date)
			for i, section := range test.sections {
				template.sections[i] = section
			}

			require.Equal(t, test.expected, template.string())
		})
	}
}
