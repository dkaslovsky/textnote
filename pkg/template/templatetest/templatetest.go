// Package templatetest provides utilities for template testing
package templatetest

import (
	"fmt"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/config"
)

// Date is a fixed date - changing this value will affect some tests
var Date = time.Date(2020, 12, 20, 1, 1, 1, 1, time.UTC)

// GetOpts returns a configuration struct for tests - changing these values will affect some tests
func GetOpts() config.Opts {
	return config.Opts{
		Header: config.HeaderOpts{
			Prefix:           "-^-",
			Suffix:           "-v-",
			TrailingNewlines: 1,
			TimeFormat:       "[Mon] 02 Jan 2006",
		},
		Section: config.SectionOpts{
			Prefix:           "_p_",
			Suffix:           "_q_",
			TrailingNewlines: 3,
			Names: []string{
				"TestSection1",
				"TestSection2",
				"TestSection3",
			},
		},
		File: config.FileOpts{
			TimeFormat: "2006-01-02",
		},
		Archive: config.ArchiveOpts{
			HeaderPrefix:             "ARCHIVEPREFIX ",
			HeaderSuffix:             " ARCHIVESUFFIX",
			SectionContentPrefix:     "[",
			SectionContentSuffix:     "]",
			SectionContentTimeFormat: "2006-01-02",
			MonthTimeFormat:          "Jan2006",
		},
		AppDir: "my/app/dir",
	}
}

// MakeItemHeader is a helper to construct an expected header property of a contentItem struct
func MakeItemHeader(date time.Time, opts config.Opts) string {
	return fmt.Sprintf("%s%s%s",
		opts.Archive.HeaderPrefix,
		date.Format(opts.Archive.SectionContentTimeFormat),
		opts.Archive.HeaderSuffix,
	)
}
