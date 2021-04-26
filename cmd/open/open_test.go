package open

import (
	"testing"
	"time"

	"github.com/dkaslovsky/textnote/pkg/template/templatetest"
	"github.com/stretchr/testify/require"
)

func TestGetLatestTemplateFile(t *testing.T) {
	opts := templatetest.GetOpts()

	type testCase struct {
		files            []string
		now              time.Time
		expectedLatest   string
		expectedNumFound int
	}

	tests := map[string]testCase{
		"empty directory": {
			files:            []string{},
			now:              time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expectedLatest:   "",
			expectedNumFound: 0,
		},
		"no timestamped template files": {
			files: []string{
				"archive-Dec2019.txt",
				"archive-2019-11-01.txt",
			},
			now:              time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expectedLatest:   "",
			expectedNumFound: 0,
		},
		"single template file in future": {
			files: []string{
				"2020-04-13.txt",
			},
			now:              time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expectedLatest:   "",
			expectedNumFound: 1,
		},
		"single template file": {
			files: []string{
				"2020-03-11.txt",
			},
			now:              time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expectedLatest:   "2020-03-11.txt",
			expectedNumFound: 1,
		},
		"multiple template files": {
			files: []string{
				"2020-03-11.txt",
				"2020-03-12.txt",
				"2020-03-13.txt",
			},
			now:              time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expectedLatest:   "2020-03-13.txt",
			expectedNumFound: 3,
		},
		"multiple template files with one in future": {
			files: []string{
				"2020-04-11.txt",
				"2020-04-12.txt",
				"2020-04-13.txt",
			},
			now:              time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expectedLatest:   "2020-04-12.txt",
			expectedNumFound: 3,
		},
		"mix of timestamped template files and other files": {
			files: []string{
				".config",
				"foobar",
				"2020-03-11.txt",
				"2020-03-12.txt",
				"2020-03-13.txt",
				"archive_April2020",
			},
			now:              time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expectedLatest:   "2020-03-13.txt",
			expectedNumFound: 3,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			latest, numFound := getLatestTemplateFile(test.files, test.now, opts.File)
			require.Equal(t, test.expectedLatest, latest)
			require.Equal(t, test.expectedNumFound, numFound)
		})
	}
}

func TestSetDateOpt(t *testing.T) {
	templateOpts := templatetest.GetOpts()

	type testCase struct {
		cmdOpts    *commandOptions
		files      []string
		now        time.Time
		expected   string
		shouldErr  bool
		shouldWarn bool
		warnThresh int
	}

	tests := map[string]testCase{
		"multiple mutually exclusive flags: date and daysBack set": {
			cmdOpts: &commandOptions{
				date:     "2020-04-11",
				daysBack: 2,
			},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			shouldErr:  true,
			shouldWarn: false,
		},
		"multiple mutually exclusive flags: date and tomorrow set": {
			cmdOpts: &commandOptions{
				date:     "2020-04-11",
				tomorrow: true,
			},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			shouldErr:  true,
			shouldWarn: false,
		},
		"multiple mutually exclusive flags: date and latest set": {
			cmdOpts: &commandOptions{
				date:   "2020-04-11",
				latest: true,
			},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			shouldErr:  true,
			shouldWarn: false,
		},
		"multiple mutually exclusive flags: daysBack and tomorrow set": {
			cmdOpts: &commandOptions{
				daysBack: 2,
				tomorrow: true,
			},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			shouldErr:  true,
			shouldWarn: false,
		},
		"multiple mutually exclusive flags: daysBack and latest set": {
			cmdOpts: &commandOptions{
				daysBack: 2,
				latest:   true,
			},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			shouldErr:  true,
			shouldWarn: false,
		},
		"multiple mutually exclusive flags: tomorrow and latest set": {
			cmdOpts: &commandOptions{
				tomorrow: true,
				latest:   true,
			},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			shouldErr:  true,
			shouldWarn: false,
		},
		"use date": {
			cmdOpts: &commandOptions{
				date: "2020-04-11",
			},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:   "2020-04-11",
			shouldErr:  false,
			shouldWarn: false,
		},
		"use daysBack": {
			cmdOpts: &commandOptions{
				daysBack: 2,
			},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:   "2020-04-10",
			shouldErr:  false,
			shouldWarn: false,
		},
		"use tomorrow": {
			cmdOpts: &commandOptions{
				tomorrow: true,
			},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:   "2020-04-13",
			shouldErr:  false,
			shouldWarn: false,
		},
		"use latest": {
			cmdOpts: &commandOptions{
				latest: true,
			},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
			expected:   "2020-04-11",
			shouldErr:  false,
			shouldWarn: false,
		},
		"no latest found": {
			cmdOpts: &commandOptions{
				latest: true,
			},
			files:      []string{},
			now:        time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
			shouldErr:  true,
			shouldWarn: false,
		},
		"default to today": {
			cmdOpts: &commandOptions{},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
			expected:   "2020-04-15",
			shouldErr:  false,
			shouldWarn: false,
		},
		"should warn on latest": {
			cmdOpts: &commandOptions{
				latest: true,
			},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
			expected:   "2020-04-11",
			shouldErr:  false,
			shouldWarn: true,
			warnThresh: 2,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			getFiles := func(dir string) ([]string, error) {
				return test.files, nil
			}
			if test.shouldWarn {
				templateOpts.TemplateFileCountThresh = test.warnThresh
			}

			shouldWarn, err := setDateOpt(test.cmdOpts, templateOpts, getFiles, test.now)
			if test.shouldErr {
				require.Error(t, err)
				return
			}
			require.Equal(t, test.shouldWarn, shouldWarn)
			require.NoError(t, err)
			require.Equal(t, test.expected, test.cmdOpts.date)
		})
	}
}

func TestSetCopyDateOpt(t *testing.T) {
	templateOpts := templatetest.GetOpts()

	type testCase struct {
		cmdOpts    *commandOptions
		files      []string
		now        time.Time
		expected   string
		shouldErr  bool
		shouldWarn bool
		warnThresh int
	}

	tests := map[string]testCase{
		"multiple mutually exclusive flags: copyDate and copyDaysBack set": {
			cmdOpts: &commandOptions{
				copyDate:     "2020-04-11",
				copyDaysBack: 2,
			},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			shouldErr:  true,
			shouldWarn: false,
		},
		"use copyDate": {
			cmdOpts: &commandOptions{
				copyDate: "2020-04-11",
			},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:   "2020-04-11",
			shouldErr:  false,
			shouldWarn: false,
		},
		"use copyDaysBack": {
			cmdOpts: &commandOptions{
				copyDaysBack: 2,
			},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:   "2020-04-10",
			shouldErr:  false,
			shouldWarn: false,
		},
		"default to latest": {
			cmdOpts: &commandOptions{},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
			expected:   "2020-04-11",
			shouldErr:  false,
			shouldWarn: false,
		},
		"no latest found": {
			cmdOpts:    &commandOptions{},
			files:      []string{},
			now:        time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
			expected:   "",
			shouldErr:  false,
			shouldWarn: false,
		},
		"should warn on latest": {
			cmdOpts: &commandOptions{},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:        time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
			expected:   "2020-04-11",
			shouldErr:  false,
			shouldWarn: true,
			warnThresh: 2,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			getFiles := func(dir string) ([]string, error) {
				return test.files, nil
			}
			if test.shouldWarn {
				templateOpts.TemplateFileCountThresh = test.warnThresh
			}

			shouldWarn, err := setCopyDateOpt(test.cmdOpts, templateOpts, getFiles, test.now)
			if test.shouldErr {
				require.Error(t, err)
				return
			}
			require.Equal(t, test.shouldWarn, shouldWarn)
			require.NoError(t, err)
			require.Equal(t, test.expected, test.cmdOpts.copyDate)
		})
	}
}
