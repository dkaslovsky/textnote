package archive

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/file"
	"github.com/dkaslovsky/TextNote/pkg/template"
	"github.com/dkaslovsky/TextNote/pkg/template/templatetest"

	"github.com/stretchr/testify/require"
)

//
// mocks
//

type testFileInfo struct {
	name  string
	isDir bool
}

func (t testFileInfo) Name() string {
	return t.name
}

func (t testFileInfo) IsDir() bool {
	return t.isDir
}

type testReadWriter struct {
	exists  bool
	toRead  string
	written string
}

func newTestReadWriter(exists bool, toRead string) *testReadWriter {
	return &testReadWriter{
		exists:  exists,
		toRead:  toRead,
		written: "",
	}
}

func (trw *testReadWriter) Read(rwable file.ReadWriteable) error {
	r := strings.NewReader(trw.toRead)
	return rwable.Load(r)
}

func (trw *testReadWriter) Overwrite(rwable file.ReadWriteable) error {
	buf := new(bytes.Buffer)
	err := rwable.Write(buf)
	if err != nil {
		return err
	}
	trw.written = buf.String()
	return nil
}

func (trw *testReadWriter) Exists(rwable file.ReadWriteable) bool {
	return trw.exists
}

//
// Tests
//

func TestAdd(t *testing.T) {
	type testCase struct {
		file         testFileInfo
		templateText string
		existing     map[string]string
		expected     map[string]string
	}

	tests := map[string]testCase{
		"add template that should not be archived": {
			file: testFileInfo{
				name:  "2020-12-20.txt",
				isDir: false,
			},
			expected: map[string]string{},
		},
		"add template from last day that should not be archived": {
			file: testFileInfo{
				name:  "2020-12-14.txt",
				isDir: false,
			},
			expected: map[string]string{},
		},
		"add template from first day that should be archived": {
			file: testFileInfo{
				name:  "2020-12-13.txt",
				isDir: false,
			},
			templateText: ``,
			expected: map[string]string{
				"Dec2020": ``,
			},
		},
		"add template from current month": {
			file: testFileInfo{
				name:  "2020-12-01.txt",
				isDir: false,
			},
			templateText: ``,
			expected: map[string]string{
				"Dec2020": ``,
			},
		},
		"add template from different month": {
			file: testFileInfo{
				name:  "2020-11-01.txt",
				isDir: false,
			},
			templateText: ``,
			expected: map[string]string{
				"Nov2020": ``,
			},
		},
		"add template from different year": {
			file: testFileInfo{
				name:  "2019-11-01.txt",
				isDir: false,
			},
			templateText: ``,
			expected: map[string]string{
				"Nov2019": ``,
			},
		},
		// "add template to existing archive": {},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			opts := templatetest.GetOpts()
			trw := newTestReadWriter(true, test.templateText)
			a := NewArchiver(opts, trw, templatetest.Date)
			for key, text := range test.existing {
				date, err := parseFileName(key, opts.Archive.MonthTimeFormat)
				require.NoError(t, err)

				m := template.NewMonthArchiveTemplate(opts, date)
				err = m.Load(strings.NewReader(text))
				require.NoError(t, err)

				a.Months[key] = m
			}

			err := a.Add(test.file)
			require.NoError(t, err)

			require.Equal(t, len(test.expected), len(a.Months))
			for key, expectedText := range test.expected {
				date, err := parseFileName(test.file.Name(), opts.File.TimeFormat)
				require.NoError(t, err)

				expectedMonthArchive := template.NewMonthArchiveTemplate(opts, date)
				err = expectedMonthArchive.Load(strings.NewReader(expectedText))
				require.NoError(t, err)

				monthArchive, found := a.Months[key]
				require.True(t, found)
				require.Equal(t, expectedMonthArchive, monthArchive)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	type testCase struct {
		text         string
		exists       bool
		existingText string
		expected     string
	}

	tests := map[string]testCase{
		"write with empty archive in archiver to new archive": {
			exists: false,
			expected: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_



_p_TestSection2_q_



_p_TestSection3_q_



`,
		},
		"write with empty archive in archiver to existing archive": {
			exists: true,
			existingText: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-15]
existingText1a



_p_TestSection2_q_



_p_TestSection3_q_
[2020-12-15]
existingText3a
[2020-12-22]
existingText3b



`,
			expected: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-15]
existingText1a



_p_TestSection2_q_



_p_TestSection3_q_
[2020-12-15]
existingText3a
[2020-12-22]
existingText3b



`,
		},
		"write to new archive": {
			text: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-17]
text1a
[2020-12-19]
text1b

_p_TestSection2_q_

_p_TestSection3_q_
[2020-12-18]
text3a
[2020-12-19]
text3b

`,
			exists: false,
			expected: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-17]
text1a
[2020-12-19]
text1b



_p_TestSection2_q_



_p_TestSection3_q_
[2020-12-18]
text3a
[2020-12-19]
text3b



`,
		},
		"write to existing archive": {
			text: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-17]
text1a
[2020-12-19]
text1b

_p_TestSection2_q_

_p_TestSection3_q_
[2020-12-18]
text3a
[2020-12-19]
text3b

`,
			exists: true,
			existingText: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-15]
existingText1a



_p_TestSection2_q_



_p_TestSection3_q_
[2020-12-15]
existingText3a
[2020-12-22]
existingText3b



`,
			expected: `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-15]
existingText1a
[2020-12-17]
text1a
[2020-12-19]
text1b



_p_TestSection2_q_



_p_TestSection3_q_
[2020-12-15]
existingText3a
[2020-12-18]
text3a
[2020-12-19]
text3b
[2020-12-22]
existingText3b



`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			opts := templatetest.GetOpts()
			date := templatetest.Date
			key := date.Format(opts.Archive.MonthTimeFormat)

			template := template.NewMonthArchiveTemplate(opts, date)
			err := template.Load(strings.NewReader(test.text))
			require.NoError(t, err)

			trw := newTestReadWriter(test.exists, test.existingText)
			a := NewArchiver(opts, trw, date)
			a.Months[key] = template

			err = a.Write()
			require.NoError(t, err)
			require.Equal(t, test.expected, trw.written)
		})
	}
}

func TestShouldArchive(t *testing.T) {
	type testCase struct {
		file     testFileInfo
		expected bool
	}

	tests := map[string]testCase{
		"archive file": {
			file: testFileInfo{
				name:  "archive-Dec2020",
				isDir: false,
			},
			expected: false,
		},
		"directory": {
			file: testFileInfo{
				name:  "somedir",
				isDir: true,
			},
			expected: false,
		},
		"hidden file": {
			file: testFileInfo{
				name:  ".config",
				isDir: false,
			},
			expected: false,
		},
		"template file": {
			file: testFileInfo{
				name:  "2020-12-29",
				isDir: false,
			},
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			opts := templatetest.GetOpts()
			require.Equal(t, test.expected, ShouldArchive(test.file, opts.Archive.FilePrefix))
		})
	}
}

func TestParseFileName(t *testing.T) {
	type testCase struct {
		fileName  string
		shouldErr bool
		expected  time.Time
	}

	tests := map[string]testCase{
		"parsable file name with extension": {
			fileName:  "2020-12-29.txt",
			shouldErr: false,
			expected:  time.Date(2020, 12, 29, 0, 0, 0, 0, time.UTC),
		},
		"parsable file name with empty extension": {
			fileName:  "2020-12-29.",
			shouldErr: false,
			expected:  time.Date(2020, 12, 29, 0, 0, 0, 0, time.UTC),
		},
		"parsable file name without extension": {
			fileName:  "2020-12-29",
			shouldErr: false,
			expected:  time.Date(2020, 12, 29, 0, 0, 0, 0, time.UTC),
		},
		"non-parsable file name with extension": {
			fileName:  "20201229.txt",
			shouldErr: true,
		},
		"non-parsable file name with empty extension": {
			fileName:  "20201229.",
			shouldErr: true,
		},
		"non-parsable file name without extension": {
			fileName:  "20201229",
			shouldErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			date, err := parseFileName(test.fileName, "2006-01-02")
			require.Equal(t, test.shouldErr, err != nil)
			require.Equal(t, test.expected, date)
		})
	}
}
