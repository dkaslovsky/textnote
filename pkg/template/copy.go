package template

import (
	"github.com/pkg/errors"
)

// CopySectionContents copies the contents of the specified section from a source template by
// appending to the contents of the target's section
func CopySectionContents(src *Template, tgt *Template, sectionName string) error {
	tgtSec, err := tgt.getSection(sectionName)
	if err != nil {
		return errors.Wrap(err, "failed to find section in target")
	}
	srcSec, err := src.getSection(sectionName)
	if err != nil {
		return errors.Wrap(err, "failed to find section in source")
	}
	tgtSec.contents = appendBeforeTrailingNewlines(tgtSec.contents, srcSec.contents)
	return nil
}

// MoveSectionContents moves the contents of the specified section from a source template by
// appending to the contents of the target's section and deleting from the source
func MoveSectionContents(src *Template, tgt *Template, sectionName string) error {
	err := CopySectionContents(src, tgt, sectionName)
	if err != nil {
		return err
	}
	srcSec, err := src.getSection(sectionName)
	if err != nil {
		return errors.Wrap(err, "failed to find section in source")
	}
	srcSec.deleteContents()
	return nil
}

func appendBeforeTrailingNewlines(tgt []string, contents []string) []string {
	tgtLen := len(tgt)
	for i := 0; i < tgtLen; i++ {
		idx := tgtLen - i - 1
		if tgt[idx] != "\n" && tgt[idx] != "" {
			updated := append(tgt[:idx+1], contents...)
			updated = append(updated, tgt[idx+1:]...)
			return updated
		}
	}
	return contents
}
