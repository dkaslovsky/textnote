package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateOpts(t *testing.T) {
	t.Run("no section names", func(t *testing.T) {
		opts := getDefaultOpts()
		opts.Section.Names = []string{}
		err := ValidateOpts(opts)
		require.Error(t, err)
	})

	t.Run("section names are not unique", func(t *testing.T) {
		opts := getDefaultOpts()
		opts.Section.Names = []string{
			"section1",
			"section2",
			"section1",
		}
		err := ValidateOpts(opts)
		require.Error(t, err)
	})

	t.Run("section names are unique", func(t *testing.T) {
		opts := getDefaultOpts()
		opts.Section.Names = []string{
			"section1",
			"section2",
			"section3",
		}
		err := ValidateOpts(opts)
		require.NoError(t, err)
	})

	t.Run("archive file prefix is empty string", func(t *testing.T) {
		opts := getDefaultOpts()
		opts.Archive.FilePrefix = ""
		err := ValidateOpts(opts)
		require.Error(t, err)
	})

	t.Run("archive file prefix is blank", func(t *testing.T) {
		opts := getDefaultOpts()
		opts.Archive.FilePrefix = "     "
		err := ValidateOpts(opts)
		require.Error(t, err)
	})

	t.Run("archive file prefix is not empty or blank", func(t *testing.T) {
		opts := getDefaultOpts()
		opts.Archive.FilePrefix = "xyzarchivexyz"
		err := ValidateOpts(opts)
		require.NoError(t, err)
	})

	t.Run("archive after days is negative", func(t *testing.T) {
		opts := getDefaultOpts()
		opts.Archive.AfterDays = -1
		err := ValidateOpts(opts)
		require.Error(t, err)
	})

	t.Run("archive after days is zero", func(t *testing.T) {
		opts := getDefaultOpts()
		opts.Archive.AfterDays = 0
		err := ValidateOpts(opts)
		require.Error(t, err)
	})

	t.Run("archive after days is one", func(t *testing.T) {
		opts := getDefaultOpts()
		opts.Archive.AfterDays = 1
		err := ValidateOpts(opts)
		require.NoError(t, err)
	})

	t.Run("file cursor line is negative", func(t *testing.T) {
		opts := getDefaultOpts()
		opts.File.CursorLine = -2
		err := ValidateOpts(opts)
		require.Error(t, err)
	})

	t.Run("file cursor line is zero", func(t *testing.T) {
		opts := getDefaultOpts()
		opts.File.CursorLine = 0
		err := ValidateOpts(opts)
		require.NoError(t, err)
	})
}
