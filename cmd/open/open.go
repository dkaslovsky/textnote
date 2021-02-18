package open

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/dkaslovsky/textnote/pkg/editor"
	"github.com/dkaslovsky/textnote/pkg/file"
	"github.com/dkaslovsky/textnote/pkg/template"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const day = 24 * time.Hour

type commandOptions struct {
	date         string
	copyDate     string
	daysBack     uint
	copyDaysBack uint
	tomorrow     bool
	sections     []string
	delete       bool
}

// CreateOpenCmd creates the open subcommand
func CreateOpenCmd() *cobra.Command {
	cmdOpts := commandOptions{}
	cmd := &cobra.Command{
		Use:          "open",
		Short:        "open a note",
		Long:         "open or create a note template",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := config.Load()
			if err != nil {
				return err
			}
			now := time.Now()
			setDateOpt(&cmdOpts, now, opts)
			setCopyDateOpt(&cmdOpts, now, opts)
			return run(opts, cmdOpts)
		},
	}
	attachOpts(cmd, &cmdOpts)
	return cmd
}

func attachOpts(cmd *cobra.Command, cmdOpts *commandOptions) {
	flags := cmd.Flags()
	flags.StringVar(&cmdOpts.date, "date", "", "date for note to be opened (defaults to today)")
	flags.StringVar(&cmdOpts.copyDate, "copy", "", "date of note for copying sections (defaults to yesterday)")
	flags.UintVarP(&cmdOpts.daysBack, "days-back", "d", 0, "number of days back from today for opening a note (ignored if date or tomorrow flags are used)")
	flags.UintVarP(&cmdOpts.copyDaysBack, "copy-back", "c", 0, "number of days back from today for copying from a note (ignored if copy flag is used)")
	flags.BoolVarP(&cmdOpts.tomorrow, "tomorrow", "t", false, "specify tomorrow as the date for note to be opened (ignored if date flag is used)")
	flags.StringSliceVarP(&cmdOpts.sections, "section", "s", []string{}, "section to copy (defaults to none)")
	flags.BoolVarP(&cmdOpts.delete, "delete", "x", false, "delete sections after copy")
}

func setDateOpt(cmdOpts *commandOptions, now time.Time, templateOpts config.Opts) {
	if cmdOpts.date != "" {
		return
	}
	if cmdOpts.tomorrow {
		// set date as tomorrow
		cmdOpts.date = now.Add(day).Format(templateOpts.Cli.TimeFormat)
		return
	}
	if cmdOpts.daysBack != 0 {
		// use daysBack if specified
		cmdOpts.date = now.Add(-day * time.Duration(cmdOpts.daysBack)).Format(templateOpts.Cli.TimeFormat)
		return
	}
	// default is today
	cmdOpts.date = now.Format(templateOpts.Cli.TimeFormat)
}

func setCopyDateOpt(cmdOpts *commandOptions, now time.Time, templateOpts config.Opts) {
	if cmdOpts.copyDate != "" {
		return
	}
	if cmdOpts.tomorrow {
		// default to today if copying to tomorrow's note
		cmdOpts.copyDate = now.Format(templateOpts.Cli.TimeFormat)
		return
	}
	if cmdOpts.copyDaysBack != 0 {
		// use copyDaysBack if specified
		cmdOpts.copyDate = now.Add(-day * time.Duration(cmdOpts.copyDaysBack)).Format(templateOpts.Cli.TimeFormat)
		return
	}
	// default is yesterday
	cmdOpts.copyDate = now.Add(-day).Format(templateOpts.Cli.TimeFormat)
}

func run(templateOpts config.Opts, cmdOpts commandOptions) error {
	date, err := time.Parse(templateOpts.Cli.TimeFormat, cmdOpts.date)
	if err != nil {
		return errors.Wrapf(err, "cannot create note for malformed date [%s]", cmdOpts.date)
	}

	t := template.NewTemplate(templateOpts, date)
	rw := file.NewReadWriter()
	ed := editor.GetEditor(os.Getenv(editor.EnvEditor))

	// open file if no sections to copy
	if len(cmdOpts.sections) == 0 {
		if !rw.Exists(t) {
			err := rw.Overwrite(t)
			if err != nil {
				return err
			}
		}
		return openInEditor(t, ed)
	}

	// load source for copy
	copyDate, err := time.Parse(templateOpts.Cli.TimeFormat, cmdOpts.copyDate)
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
	return openInEditor(t, ed)
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

func openInEditor(t *template.Template, ed *editor.Editor) error {
	if t.GetFileCursorLine() > 0 && !ed.Supported {
		log.Printf("Editor [%s] only supported with its default arguments, additional configuration ignored", ed.Cmd)
	}
	if ed.Default {
		log.Printf("Environment variable [%s] not set, attempting to use default editor [%s]", editor.EnvEditor, ed.Cmd)
	}
	return file.Open(t, ed)
}
