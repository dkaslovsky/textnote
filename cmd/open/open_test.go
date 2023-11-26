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
	type testCase struct {
		cmdOpts          *commandOptions
		files            []string
		now              time.Time
		expectedDate     string
		expectedNumFiles int
		shouldErr        bool
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
			now:              time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expectedDate:     "2020-04-11",
			expectedNumFiles: 0,
			shouldErr:        false,
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
			now:              time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expectedDate:     "2020-04-10",
			expectedNumFiles: 0,
			shouldErr:        false,
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
			now:              time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expectedDate:     "2020-04-13",
			expectedNumFiles: 0,
			shouldErr:        false,
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
			now:              time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
			expectedDate:     "2020-04-11",
			expectedNumFiles: 3,
			shouldErr:        false,
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
			now:              time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
			expectedDate:     "2020-04-15",
			expectedNumFiles: 0,
			shouldErr:        false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			getFiles := func(dir string) ([]string, error) {
				return test.files, nil
			}
			templateOpts := templatetest.GetOpts()

			// test
			numFiles, err := setDateOpt(test.cmdOpts, templateOpts, getFiles, test.now)
			if test.shouldErr {
				require.Error(t, err)
				return
			}
			require.Equal(t, test.expectedNumFiles, numFiles)
			require.NoError(t, err)
			require.Equal(t, test.expectedDate, test.cmdOpts.date)
		})
	}
}

func TestSetCopyDateOpt(t *testing.T) {
	type testCase struct {
		cmdOpts          *commandOptions
		files            []string
		now              time.Time
		expectedDate     string
		expectedNumFiles int
		shouldErr        bool
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
			now:              time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expectedDate:     "2020-04-11",
			expectedNumFiles: 0,
			shouldErr:        false,
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
			now:              time.Date(2020, 4, 12, 0, 0, 0, 0, time.UTC),
			expectedDate:     "2020-04-10",
			expectedNumFiles: 0,
			shouldErr:        false,
		},
		"default to latest": {
			cmdOpts: &commandOptions{},
			files: []string{
				"2020-04-11.txt",
				"2020-04-10.txt",
				"2020-04-09.txt",
			},
			now:              time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
			expectedDate:     "2020-04-11",
			expectedNumFiles: 3,
			shouldErr:        false,
		},
		"no latest found": {
			cmdOpts:          &commandOptions{},
			files:            []string{},
			now:              time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
			expectedDate:     "",
			expectedNumFiles: 0,
			shouldErr:        false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			getFiles := func(dir string) ([]string, error) {
				return test.files, nil
			}
			templateOpts := templatetest.GetOpts()

			// test
			numFiles, err := setCopyDateOpt(test.cmdOpts, templateOpts, getFiles, test.now)
			if test.shouldErr {
				require.Error(t, err)
				return
			}
			require.Equal(t, test.expectedNumFiles, numFiles)
			require.NoError(t, err)
			require.Equal(t, test.expectedDate, test.cmdOpts.copyDate)
		})
	}
}

func TestSetDeleteOpts(t *testing.T) {
	type testCase struct {
		cmdOpts                *commandOptions
		expectedDeleteSections bool
		expectedDeleteEmpty    bool
	}

	tests := map[string]testCase{
		"deleteFlagVal = 0": {
			cmdOpts: &commandOptions{
				deleteFlagVal: 0,
			},
			expectedDeleteSections: false,
			expectedDeleteEmpty:    false,
		},
		"deleteFlagVal < 0": {
			cmdOpts: &commandOptions{
				deleteFlagVal: -1,
			},
			expectedDeleteSections: false,
			expectedDeleteEmpty:    false,
		},
		"deleteFlagVal = 1": {
			cmdOpts: &commandOptions{
				deleteFlagVal: 1,
			},
			expectedDeleteSections: true,
			expectedDeleteEmpty:    false,
		},
		"deleteFlagVal = 2": {
			cmdOpts: &commandOptions{
				deleteFlagVal: 2,
			},
			expectedDeleteSections: true,
			expectedDeleteEmpty:    true,
		},
		"deleteFlagVal > 2": {
			cmdOpts: &commandOptions{
				deleteFlagVal: 3,
			},
			expectedDeleteSections: true,
			expectedDeleteEmpty:    true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			setDeleteOpts(test.cmdOpts)
			require.Equal(t, test.expectedDeleteSections, test.cmdOpts.deleteSections)
			require.Equal(t, test.expectedDeleteEmpty, test.cmdOpts.deleteEmpty)
		})
	}
}
