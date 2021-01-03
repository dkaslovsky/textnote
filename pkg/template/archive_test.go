package template

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/template/templatetest"
	"github.com/stretchr/testify/require"
)

func TestArchiveGetFilePath(t *testing.T) {
	t.Run("get file path", func(t *testing.T) {
		opts := templatetest.GetOpts()
		template := NewMonthArchiveTemplate(opts, templatetest.Date)
		filePath := template.GetFilePath()
		require.True(t, strings.HasPrefix(filePath, opts.AppDir))
		require.True(t, strings.HasSuffix(filePath, fileExt))
		require.Equal(t,
			ArchiveFilePrefix+templatetest.Date.Format(opts.Archive.MonthTimeFormat),
			stripPrefixSuffix(filePath,
				fmt.Sprintf("%s/", opts.AppDir),
				fmt.Sprintf(".%s", fileExt),
			))
	})
}

func TestArchiveCopySectionContents(t *testing.T) {
	type testCase struct {
		sectionName      string
		existingContents []contentItem
		sourceDate       time.Time
		sourceContents   []contentItem
		expectedContents []contentItem
	}

	tests := map[string]testCase{
		"copy empty contents into empty section": {
			sectionName:      "TestSection1",
			existingContents: []contentItem{},
			sourceDate:       templatetest.Date,
			sourceContents:   []contentItem{},
			expectedContents: []contentItem{},
		},
		"copy empty contents into populated section": {
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
		"copy contents with single element into empty section": {
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
		"copy contents with single element into populated section": {
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
		"copy contents with multiple element into empty section": {
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
		"copy contents with multiple elements into populated section": {
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
		"copy contents from source with same date": {
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

			err := template.CopySectionContents(src, test.sectionName)
			require.NoError(t, err)
			require.Equal(t, template.sections[template.sectionIdx[test.sectionName]].contents, test.expectedContents)
		})
	}
}

func TestArchiveCopySectionContentsFail(t *testing.T) {
	t.Run("section does not exist in template", func(t *testing.T) {
		toCopy := "toBeCopied"
		opts := templatetest.GetOpts()
		template := NewMonthArchiveTemplate(opts, templatetest.Date)
		src := NewTemplate(opts, templatetest.Date)
		src.sections = append(src.sections, newSection(toCopy))
		src.sectionIdx[toCopy] = len(src.sections) - 1

		err := template.CopySectionContents(src, toCopy)
		require.Error(t, err)
	})

	t.Run("section does not exist in source", func(t *testing.T) {
		toCopy := "toBeCopied"
		opts := templatetest.GetOpts()
		template := NewMonthArchiveTemplate(opts, templatetest.Date)
		template.sections = append(template.sections, newSection(toCopy))
		template.sectionIdx[toCopy] = len(template.sections) - 1
		src := NewTemplate(opts, templatetest.Date)

		err := template.CopySectionContents(src, toCopy)
		require.Error(t, err)
	})
}
