package open

import (
	"fmt"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/config"
	"github.com/dkaslovsky/TextNote/pkg/file"
	"github.com/dkaslovsky/TextNote/pkg/template"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type commandOptions struct {
	Copy   []string
	Delete bool
}

func attachOpts(cmd *cobra.Command, cmdOpts *commandOptions) {
	flags := cmd.Flags()
	flags.StringSliceVarP(&cmdOpts.Copy, "copy", "c", []string{}, "section names to copy")
	flags.BoolVarP(&cmdOpts.Delete, "delete", "d", false, "delete previous day's section after copy (no-op without copy)")
}

func run(templateOpts config.Opts, cmdOpts commandOptions, date time.Time) error {
	t := template.NewTemplate(templateOpts)
	t.SetDate(date)

	// open file if no further operations (copy/move)
	if len(cmdOpts.Copy) == 0 {
		err := file.WriteIfNotExists(t)
		if err != nil {
			return err
		}
		return file.OpenInEditor(t)
	}

	// load t from file if it exists
	if file.Exists(t) {
		err := file.Read(t)
		if err != nil {
			return errors.Wrap(err, "cannot load file")
		}
	}

	// load source file
	src := template.NewTemplate(templateOpts)
	src.SetDate(date.Add(-24 * time.Hour))
	err := file.Read(src)
	if err != nil {
		return errors.Wrap(err, "cannot read source file for copy")
	}

	if cmdOpts.Delete {
		err := moveSections(src, t, cmdOpts.Copy)
		if err != nil {
			return err
		}
		err = file.Overwrite(src)
		if err != nil {
			return errors.Wrap(err, "failed to save changes to source file")
		}
	} else {
		err = copySections(src, t, cmdOpts.Copy)
		if err != nil {
			return err
		}
	}

	err = file.Overwrite(t)
	if err != nil {
		return errors.Wrap(err, "failed to write file")
	}
	return file.OpenInEditor(t)
}

func open(templateOpts config.Opts, date time.Time) error {
	t := template.NewTemplate(templateOpts)
	t.SetDate(date)
	if !file.Exists(t) {
		return fmt.Errorf("file [%s] for template does not exist", t.GetFilePath())
	}
	return file.OpenInEditor(t)
}

func copySections(src *template.Template, tgt *template.Template, sectionNames []string) error {
	for _, sectionName := range sectionNames {
		err := template.CopySectionContents(src, tgt, sectionName)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("cannot copy section [%s] from source to target", sectionName))
		}
	}
	return nil
}

func moveSections(src *template.Template, tgt *template.Template, sectionNames []string) error {
	for _, sectionName := range sectionNames {
		err := template.MoveSectionContents(src, tgt, sectionName)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("cannot copy section [%s] from source to target", sectionName))
		}
	}
	return nil
}
