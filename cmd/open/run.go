package open

import (
	"fmt"
	"time"

	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/dkaslovsky/textnote/pkg/file"
	"github.com/dkaslovsky/textnote/pkg/template"
	"github.com/pkg/errors"
)

// Run executes the workflow for creating/opening a note
func Run(templateOpts config.Opts, cmdOpts CommandOptions, sections []string, date time.Time, copyDate time.Time) error {
	rw := file.NewReadWriter()
	t := template.NewTemplate(templateOpts, date)

	// open file if no further operations (copy/move)
	if !cmdOpts.Copy {
		if !rw.Exists(t) {
			err := rw.Overwrite(t)
			if err != nil {
				return err
			}
		}
		return file.OpenInVim(t)
	}

	if len(sections) == 0 {
		sections = templateOpts.Section.Names
	}

	// load template contents if it exists
	if rw.Exists(t) {
		err := rw.Read(t)
		if err != nil {
			return errors.Wrap(err, "cannot load template file")
		}
	}
	// load source for copy
	src := template.NewTemplate(templateOpts, copyDate)
	err := rw.Read(src)
	if err != nil {
		return errors.Wrap(err, "cannot read source file for copy")
	}
	// copy from source to template
	err = copySections(src, t, sections)
	if err != nil {
		return err
	}

	if cmdOpts.Delete {
		err = deleteSections(src, sections)
		if err != nil {
			return errors.Wrap(err, "failed to remove section content from source file")
		}
		err = rw.Overwrite(src)
		if err != nil {
			return errors.Wrap(err, "failed to save changes to source file")
		}
	}

	err = rw.Overwrite(t)
	if err != nil {
		return errors.Wrap(err, "failed to write file")
	}
	return file.OpenInVim(t)
}

func copySections(src *template.Template, tgt *template.Template, sectionNames []string) error {
	for _, sectionName := range sectionNames {
		err := tgt.CopySectionContents(src, sectionName)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("cannot copy section [%s] from source to target", sectionName))
		}
	}
	return nil
}

func deleteSections(t *template.Template, sectionNames []string) error {
	for _, sectionName := range sectionNames {
		err := t.DeleteSectionContents(sectionName)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("cannot delete section [%s] from template", sectionName))
		}
	}
	return nil
}
