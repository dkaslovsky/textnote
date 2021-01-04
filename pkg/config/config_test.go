package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateConfig(t *testing.T) {
	t.Run("no section names", func(t *testing.T) {
		opts := getDefaultOpts()
		opts.Section.Names = []string{}
		err := ValidateConfig(opts)
		require.Error(t, err)
	})

	t.Run("section names are not unique", func(t *testing.T) {
		opts := getDefaultOpts()
		opts.Section.Names = []string{
			"section1",
			"section2",
			"section1",
		}
		err := ValidateConfig(opts)
		require.Error(t, err)
	})

	t.Run("archive file prefix is empty string", func(t *testing.T) {
		opts := getDefaultOpts()
		opts.Archive.FilePrefix = ""
		err := ValidateConfig(opts)
		require.Error(t, err)
	})

	t.Run("archive file prefix is blank", func(t *testing.T) {
		opts := getDefaultOpts()
		opts.Archive.FilePrefix = "     "
		err := ValidateConfig(opts)
		require.Error(t, err)
	})

	// t.Run("invalid header time format", func(t *testing.T) {
	// 	opts := getDefaultOpts()
	// 	opts.Header.TimeFormat = "xyz"
	// 	err := ValidateConfig(opts)
	// 	require.Error(t, err)
	// })
}
