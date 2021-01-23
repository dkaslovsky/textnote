package open

import (
	"fmt"
	"time"

	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/dkaslovsky/textnote/pkg/file"
	"github.com/dkaslovsky/textnote/pkg/template"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type commandOptions struct {
	format    string
	date      string
	yesterday bool
	copyDate  string
	sections  []string
	delete    bool
}

// CreateOpenCmd creates the open subcommand
func CreateOpenCmd() *cobra.Command {
	cmdOpts := commandOptions{}
	cmd := &cobra.Command{
		Use:   "open",
		Short: "open a note",
		Long:  "open or create a note template",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := config.LoadOrCreate()
			if err != nil {
				return err
			}
			applyDefaults(opts, &cmdOpts)
			return run(opts, cmdOpts)
		},
	}
	attachOpts(cmd, &cmdOpts)
	return cmd
}

func attachOpts(cmd *cobra.Command, cmdOpts *commandOptions) {
	flags := cmd.Flags()
	flags.StringVarP(&cmdOpts.format, "format", "f", "", "override for time format to parse date flags specified in configuration")
	flags.StringVarP(&cmdOpts.date, "date", "d", "", "date for note to be opened (defaults to today)")
	flags.BoolVarP(&cmdOpts.yesterday, "yesterday", "y", false, "use yesterday's date for note (ignored if date is specified)")
	flags.StringVarP(&cmdOpts.copyDate, "copyDate", "c", "", "date of note for copying sections (defaults to yesterday)")
	flags.StringSliceVarP(&cmdOpts.sections, "section", "s", []string{}, "section to copy")
	flags.BoolVarP(&cmdOpts.delete, "delete", "x", false, "delete sections after copy")
}

func applyDefaults(templateOpts config.Opts, cmdOpts *commandOptions) {
	now := time.Now()
	if cmdOpts.format == "" {
		cmdOpts.format = templateOpts.Cli.TimeFormat
	}
	if cmdOpts.date == "" {
		if cmdOpts.yesterday {
			cmdOpts.date = now.Add(-24 * time.Hour).Format(cmdOpts.format)
		} else {
			cmdOpts.date = now.Format(cmdOpts.format)
		}
	}
	if cmdOpts.copyDate == "" {
		cmdOpts.copyDate = now.Add(-24 * time.Hour).Format(cmdOpts.format)
	}
}

func run(templateOpts config.Opts, cmdOpts commandOptions) error {
	date, err := time.Parse(cmdOpts.format, cmdOpts.date)
	if err != nil {
		return errors.Wrapf(err, "cannot create note for malformed date [%s]", cmdOpts.date)
	}

	t := template.NewTemplate(templateOpts, date)
	rw := file.NewReadWriter()

	// open file if no sections to copy
	if len(cmdOpts.sections) == 0 {
		if !rw.Exists(t) {
			err := rw.Overwrite(t)
			if err != nil {
				return err
			}
		}
		return file.OpenInVim(t)
	}

	// load source for copy
	copyDate, err := time.Parse(cmdOpts.format, cmdOpts.copyDate)
	if err != nil {
		return errors.Wrapf(err, "cannot copy note from malformed date [%s]", cmdOpts.copyDate)
	}
	src := template.NewTemplate(templateOpts, copyDate)
	err = rw.Read(src)
	if err != nil {
		return errors.Wrap(err, "cannot read source file for copy")
	}
	// load template contents if it exists
	if rw.Exists(t) {
		err := rw.Read(t)
		if err != nil {
			return errors.Wrap(err, "cannot load template file")
		}
	}
	// copy from source to template
	err = copySections(src, t, cmdOpts.sections)
	if err != nil {
		return err
	}

	if cmdOpts.delete {
		err = deleteSections(src, cmdOpts.sections)
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
