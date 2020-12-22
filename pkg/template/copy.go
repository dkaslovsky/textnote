package template

import (
	"time"

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
	tgtSec.contents = insert(tgtSec.contents, srcSec.contents)
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

// ArchiveSectionContents archives the contents of the specified section from a source template by
// appending to the contents of the target's section and prepending each element of the contents
// with the date of the source template
func ArchiveSectionContents(src *Template, tgt archiveTarget, sectionName string) error {
	tgtSec, err := tgt.getSection(sectionName)
	if err != nil {
		return errors.Wrap(err, "failed to find section in target")
	}
	srcSec, err := src.getSection(sectionName)
	if err != nil {
		return errors.Wrap(err, "failed to find section in source")
	}

	contents := []string{}
	for _, content := range srcSec.contents {
		if content == "" || content == "\n" {
			continue
		}
		contents = append(contents, content)
	}
	if len(contents) > 0 {
		dateStr := tgt.makeSectionContentPrefix(src.date)
		contents = append([]string{dateStr}, contents...)
		tgtSec.contents = insert(tgtSec.contents, contents)
	}
	return nil
}

// archiveTarget is the interface accepted for the target template of ArchiveSectionContents
type archiveTarget interface {
	getSection(string) (*section, error)
	makeSectionContentPrefix(time.Time) string
}

// insert inserts contents into tgt before any trailing empty elements, omitting trailing empty
// elements of contents
func insert(tgt []string, contents []string) []string {
	if len(contents) == 0 {
		return tgt
	}
	if len(tgt) == 0 {
		return contents
	}

	contentsIdx := getLastPopulatedIndex(contents) + 1
	insertIdx := getLastPopulatedIndex(tgt) + 1

	updated := []string{}
	updated = append(updated, tgt[:insertIdx]...)
	updated = append(updated, contents[:contentsIdx]...)
	updated = append(updated, tgt[insertIdx:]...)
	return updated
}

func getLastPopulatedIndex(s []string) int {
	ln := len(s)
	for i := 0; i < ln; i++ {
		idx := ln - i - 1
		if s[idx] != "\n" && s[idx] != "" {
			return idx
		}
	}
	return -1
}
