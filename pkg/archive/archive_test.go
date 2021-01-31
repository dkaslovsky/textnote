package archive

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/dkaslovsky/textnote/pkg/file"
	"github.com/dkaslovsky/textnote/pkg/template"
	"github.com/dkaslovsky/textnote/pkg/template/templatetest"

	"github.com/stretchr/testify/require"
)

//
// mocks
//

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
		fileName     string
		templateText string
		existing     map[string]string
		expected     map[string]string
	}

	tests := map[string]testCase{
		"add template that should not be archived": {
			fileName: "2020-12-20.txt",
			expected: map[string]string{},
		},
		"add template from last day that should not be archived": {
			fileName: "2020-12-14.txt",
			expected: map[string]string{},
		},
		"add template from first day that should be archived": {
			fileName: "2020-12-13.txt",
			templateText: `-^-[Sun] 13 Dec 2020-v-

_p_TestSection1_q_
text1
  text2



_p_TestSection2_q_



_p_TestSection3_q_



`,
			expected: map[string]string{
				"Dec2020": `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-13]
text1
  text2



_p_TestSection2_q_



_p_TestSection3_q_



`,
			},
		},
		"add template from current month": {
			fileName: "2020-12-01.txt",
			templateText: `-^-[Tue] 01 Dec 2020-v-

_p_TestSection1_q_
text1
  text2



_p_TestSection2_q_



_p_TestSection3_q_



`,
			expected: map[string]string{
				"Dec2020": `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-01]
text1
  text2



_p_TestSection2_q_



_p_TestSection3_q_



`,
			},
		},
		"add template from different month": {
			fileName: "2020-11-01.txt",
			templateText: `-^-[Sun] 01 Nov 2020-v-

_p_TestSection1_q_
text1
  text2



_p_TestSection2_q_



_p_TestSection3_q_



`,
			expected: map[string]string{
				"Nov2020": `ARCHIVEPREFIX Nov2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-11-01]
text1
  text2



_p_TestSection2_q_



_p_TestSection3_q_



`,
			},
		},
		"add template from different year": {
			fileName: "2019-11-02.txt",
			templateText: `-^-[Sat] 02 Nov 2019-v-

_p_TestSection1_q_
text1
  text2



_p_TestSection2_q_



_p_TestSection3_q_



`,
			expected: map[string]string{
				"Nov2019": `ARCHIVEPREFIX Nov2019 ARCHIVESUFFIX

_p_TestSection1_q_
[2019-11-02]
text1
  text2



_p_TestSection2_q_



_p_TestSection3_q_



`,
			},
		},
		"add template with earlier date to existing archive": {
			fileName: "2020-12-01.txt",
			templateText: `-^-[Tue] 01 Dec 2020-v-

_p_TestSection1_q_
text1
  text2



_p_TestSection2_q_



_p_TestSection3_q_



`,
			existing: map[string]string{
				"Dec2020": `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-02]
existingText1
  existingText2
existingText3

_p_TestSection2_q_



_p_TestSection3_q_



`,
			},
			expected: map[string]string{
				"Dec2020": `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-01]
text1
  text2
[2020-12-02]
existingText1
  existingText2
existingText3



_p_TestSection2_q_



_p_TestSection3_q_



`,
			},
		},
		"add template with later date to existing archive": {
			fileName: "2020-12-02.txt",
			templateText: `-^-[Wed] 02 Dec 2020-v-

_p_TestSection1_q_
text1
  text2



_p_TestSection2_q_



_p_TestSection3_q_



`,
			existing: map[string]string{
				"Dec2020": `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-01]
existingText1
  existingText2
existingText3

_p_TestSection2_q_



_p_TestSection3_q_



`,
			},
			expected: map[string]string{
				"Dec2020": `ARCHIVEPREFIX Dec2020 ARCHIVESUFFIX

_p_TestSection1_q_
[2020-12-01]
existingText1
  existingText2
existingText3
[2020-12-02]
text1
  text2



_p_TestSection2_q_



_p_TestSection3_q_



`,
			},
		},
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

			err := a.Add(test.fileName)
			require.NoError(t, err)

			require.Equal(t, len(test.expected), len(a.Months))
			for key, expectedText := range test.expected {
				buf := new(bytes.Buffer)
				monthArchive, found := a.Months[key]
				require.True(t, found)
				err := monthArchive.Write(buf)
				require.NoError(t, err)
				require.Equal(t, expectedText, buf.String())
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
