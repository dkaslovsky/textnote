package open

import (
	"testing"
	"time"

	"github.com/dkaslovsky/textnote/pkg/template/templatetest"
	"github.com/stretchr/testify/require"
)

func TestGetLatestFile(t *testing.T) {
	opts := templatetest.GetOpts()

	type testCase struct {
		files     []string
		now       time.Time
		expected  string
		shouldErr bool
	}

	tests := map[string]testCase{
		"empty directory": {
			files:     []string{},
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:  "",
			shouldErr: true,
		},
		"no timestamped template files": {
			files: []string{
				"archive-Dec2019.txt",
				"archive-2019-11-01.txt",
			},
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:  "",
			shouldErr: true,
		},
		"single template file in future": {
			files: []string{
				"2020-04-13.txt",
			},
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:  "",
			shouldErr: true,
		},
		"single template file": {
			files: []string{
				"2020-03-11.txt",
			},
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:  "2020-03-11.txt",
			shouldErr: false,
		},
		"multiple template files": {
			files: []string{
				"2020-03-11.txt",
				"2020-03-12.txt",
				"2020-03-13.txt",
			},
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:  "2020-03-13.txt",
			shouldErr: false,
		},
		"multiple template files with one in future": {
			files: []string{
				"2020-04-11.txt",
				"2020-04-12.txt",
				"2020-04-13.txt",
			},
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:  "2020-04-12.txt",
			shouldErr: false,
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
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:  "2020-03-13.txt",
			shouldErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			latest, err := getLatestFile(test.files, test.now, opts.File)
			if test.shouldErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.expected, latest)
		})
	}
}

func TestSetDateOpt(t *testing.T) {
	templateOpts := templatetest.GetOpts()

	type testCase struct {
		cmdOpts   *commandOptions
		files     []string
		now       time.Time
		expected  string
		shouldErr bool
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
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			shouldErr: true,
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
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			shouldErr: true,
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
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			shouldErr: true,
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
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			shouldErr: true,
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
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			shouldErr: true,
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
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			shouldErr: true,
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
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:  "2020-04-11",
			shouldErr: false,
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
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:  "2020-04-10",
			shouldErr: false,
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
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:  "2020-04-13",
			shouldErr: false,
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
			now:       time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
			expected:  "2020-04-11",
			shouldErr: false,
		},
		"no latest found": {
			cmdOpts: &commandOptions{
				latest: true,
			},
			files:     []string{},
			now:       time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
			shouldErr: true,
		},
		"default to today": {
			cmdOpts: &commandOptions{},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:       time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
			expected:  "2020-04-15",
			shouldErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			getFiles := func(dir string) ([]string, error) {
				return test.files, nil
			}
			err := setDateOpt(test.cmdOpts, templateOpts, getFiles, test.now)
			if test.shouldErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.expected, test.cmdOpts.date)
		})
	}
}

func TestSetCopyDateOpt(t *testing.T) {
	templateOpts := templatetest.GetOpts()

	type testCase struct {
		cmdOpts   *commandOptions
		files     []string
		now       time.Time
		expected  string
		shouldErr bool
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
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			shouldErr: true,
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
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:  "2020-04-11",
			shouldErr: false,
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
			now:       time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expected:  "2020-04-10",
			shouldErr: false,
		},
		"default to latest": {
			cmdOpts: &commandOptions{},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:       time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
			expected:  "2020-04-11",
			shouldErr: false,
		},
		"no latest found": {
			cmdOpts:   &commandOptions{},
			files:     []string{},
			now:       time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
			shouldErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			getFiles := func(dir string) ([]string, error) {
				return test.files, nil
			}
			err := setCopyDateOpt(test.cmdOpts, templateOpts, getFiles, test.now)
			if test.shouldErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.expected, test.cmdOpts.copyDate)
		})
	}
}
