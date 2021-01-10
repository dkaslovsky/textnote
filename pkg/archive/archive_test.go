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

type testWriter struct {
	written string
}

func newTestWriter() *testWriter {
	return &testWriter{
		written: "",
	}
}

func (w *testWriter) write(rw file.ReadWriteable) error {
	buf := new(bytes.Buffer)
	err := rw.Write(buf)
	if err != nil {
		return err
	}
	w.written = buf.String()
	return nil
}

func TestWrite(t *testing.T) {
	type testCase struct {
		key          string
		text         string
		exists       bool
		existingText string
		expected     string
	}

	tests := map[string]testCase{
		"write to new archive": {
			key: "Dec2020",
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
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			opts := templatetest.GetOpts()

			template := template.NewMonthArchiveTemplate(opts, templatetest.Date)
			err := template.Load(strings.NewReader(test.text))
			require.NoError(t, err)

			a := NewArchiver(opts, templatetest.Date)
			a.Months[test.key] = template

			w := newTestWriter()
			err = a.Write(w.write, func(file.ReadWriteable) bool {
				return test.exists
			})
			require.NoError(t, err)
			require.Equal(t, test.expected, w.written)
		})
	}
}

func TestShouldNotArchive(t *testing.T) {
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
			expected: true,
		},
		"directory": {
			file: testFileInfo{
				name:  "somedir",
				isDir: true,
			},
			expected: true,
		},
		"hidden file": {
			file: testFileInfo{
				name:  ".config",
				isDir: false,
			},
			expected: true,
		},
		"template file": {
			file: testFileInfo{
				name:  "2020-12-29",
				isDir: false,
			},
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			a := NewArchiver(templatetest.GetOpts(), templatetest.Date)
			require.Equal(t, test.expected, a.shouldNotArchive(test.file))
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
