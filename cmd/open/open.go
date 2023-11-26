package open

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
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
	// mutually exclusive flags for date to open
	date     string
	daysBack uint
	tomorrow bool
	latest   bool

	// mutually exclusive flags for copy date
	copyDate     string
	copyDaysBack uint

	deleteFlagVal  int  // count of number of times delete flag is passed
	deleteSections bool // delete sections on copy (deleteFlagVal > 0)
	deleteEmpty    bool // delete file if empty after deleting sections (deleteFlagVal > 1)

	sections []string
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
			shouldWarnForDate, err := setDateOpt(&cmdOpts, opts, getDirFiles, now)
			if err != nil {
				return err
			}
			shouldWarnForCopy, err := setCopyDateOpt(&cmdOpts, opts, getDirFiles, now)
			if err != nil {
				return err
			}
			if shouldWarnForDate || shouldWarnForCopy {
				log.Printf("found more than %d template files when searching for latest, consider running archive command to more efficient performance", opts.TemplateFileCountThresh)
			}
			setDeleteOpts(&cmdOpts)
			return run(opts, cmdOpts)
		},
	}
	attachOpts(cmd, &cmdOpts)
	return cmd
}

func attachOpts(cmd *cobra.Command, cmdOpts *commandOptions) {
	flags := cmd.Flags()

	// mutually exclusive flags for date to open
	flags.StringVar(&cmdOpts.date, "date", "", "date for note to be opened (defaults to today)")
	flags.UintVarP(&cmdOpts.daysBack, "days-back", "d", 0, "number of days back from today for opening a note (cannot be used with date, tomorrow, or latest flags)")
	flags.BoolVarP(&cmdOpts.tomorrow, "tomorrow", "t", false, "specify tomorrow as the date for note to be opened (cannot be used with date, days-back, or latest flags)")
	flags.BoolVarP(&cmdOpts.latest, "latest", "l", false, "specify the most recent dated note to be opened (cannot be used with date, days-back, or tomorrow flags)")

	// mutually exclusive flags for copy date
	flags.StringVar(&cmdOpts.copyDate, "copy", "", "date of note for copying sections (defaults to date of most recent note, cannot be used with copy-back flag)")
	flags.UintVarP(&cmdOpts.copyDaysBack, "copy-back", "c", 0, "number of days back from today for copying from a note (cannot be used with copy flag)")

	flags.StringSliceVarP(&cmdOpts.sections, "section", "s", []string{}, "section to copy (defaults to none)")
	flags.CountVarP(&cmdOpts.deleteFlagVal, "delete", "x", "delete sections after copy (pass flag twice to also delete empty source note)")
}

func setDateOpt(cmdOpts *commandOptions, templateOpts config.Opts, getFiles func(string) ([]string, error), now time.Time) (bool, error) {
	errMutuallyExclusive := errors.New("only one of [date, days-back, tomorrow, latest] flags may be used")
	date := ""
	shouldWarn := false

	if cmdOpts.date != "" {
		date = cmdOpts.date
	}

	if cmdOpts.daysBack != 0 {
		if date != "" {
			return shouldWarn, errMutuallyExclusive
		}
		date = now.Add(-day * time.Duration(cmdOpts.daysBack)).Format(templateOpts.Cli.TimeFormat)
	}

	if cmdOpts.tomorrow {
		if date != "" {
			return shouldWarn, errMutuallyExclusive
		}
		date = now.Add(day).Format(templateOpts.Cli.TimeFormat)
	}

	if cmdOpts.latest {
		if date != "" {
			return shouldWarn, errMutuallyExclusive
		}

		files, err := getFiles(templateOpts.AppDir)
		if err != nil {
			return shouldWarn, err
		}
		latest, numFound := getLatestTemplateFile(files, now, templateOpts.File)
		if latest == "" {
			return shouldWarn, fmt.Errorf("failed to find latest template file in [%s]", templateOpts.AppDir)
		}
		if templateOpts.File.Ext != "" {
			latest = strings.TrimSuffix(latest, fmt.Sprintf(".%s", templateOpts.File.Ext))
		}
		date = latest

		// check if warning for too many template files should be displayed
		shouldWarn = numFound > templateOpts.TemplateFileCountThresh
	}

	// default to today
	if date == "" {
		date = now.Format(templateOpts.Cli.TimeFormat)
	}

	cmdOpts.date = date
	return shouldWarn, nil
}

func setCopyDateOpt(cmdOpts *commandOptions, templateOpts config.Opts, getFiles func(string) ([]string, error), now time.Time) (bool, error) {
	if cmdOpts.copyDate != "" && cmdOpts.copyDaysBack != 0 {
		return false, errors.New("only one of [copy, copy-back] flags may be used")
	}

	if cmdOpts.copyDate != "" {
		return false, nil
	}
	if cmdOpts.copyDaysBack != 0 {
		cmdOpts.copyDate = now.Add(-day * time.Duration(cmdOpts.copyDaysBack)).Format(templateOpts.Cli.TimeFormat)
		return false, nil
	}

	// default to latest
	files, err := getFiles(templateOpts.AppDir)
	if err != nil {
		return false, err
	}
	latest, numFound := getLatestTemplateFile(files, now, templateOpts.File)
	if templateOpts.File.Ext != "" {
		latest = strings.TrimSuffix(latest, fmt.Sprintf(".%s", templateOpts.File.Ext))
	}
	cmdOpts.copyDate = latest

	// check if warning for too many template files should be displayed
	shouldWarn := numFound > templateOpts.TemplateFileCountThresh

	return shouldWarn, nil
}

func setDeleteOpts(cmdOpts *commandOptions) {
	cmdOpts.deleteSections = cmdOpts.deleteFlagVal > 0
	cmdOpts.deleteEmpty = cmdOpts.deleteFlagVal > 1
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
	if cmdOpts.copyDate == "" {
		return fmt.Errorf("cannot find note to copy, [%s] might be empty", templateOpts.AppDir)
	}
	if cmdOpts.copyDate == cmdOpts.date {
		return fmt.Errorf("copying from note dated [%s] not allowed when writing to note for date [%s]", cmdOpts.copyDate, cmdOpts.date)
	}
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

	if cmdOpts.deleteSections {
		err = deleteSections(src, cmdOpts.sections)
		if err != nil {
			return errors.Wrap(err, "failed to remove section content from source file")
		}

		if cmdOpts.deleteEmpty && src.IsEmpty() {
			err = os.Remove(src.GetFilePath())
			if err != nil {
				return errors.Wrapf(err, "failed to remove empty source file")
			}
		} else {
			err = rw.Overwrite(src)
			if err != nil {
				return errors.Wrap(err, "failed to save changes to source file")
			}

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
	if t.GetFileCursorLine() > 1 && !ed.Supported {
		log.Printf("Editor [%s] only supported with its default arguments, additional configuration ignored", ed.Cmd)
	}
	if ed.Default {
		log.Printf("Environment variable [%s] not set, attempting to use default editor [%s]", editor.EnvEditor, ed.Cmd)
	}
	return ed.Open(t)
}

func getLatestTemplateFile(files []string, now time.Time, opts config.FileOpts) (string, int) {
	latest := ""
	delta := math.Inf(1)
	numTemplateFiles := 0

	for _, f := range files {
		fileTime, ok := template.ParseTemplateFileName(f, opts)
		if !ok {
			// skip archive files and other non-template files that cannot be parsed
			continue
		}
		numTemplateFiles++
		curdelta := now.Sub(fileTime).Hours()
		if curdelta < 0 {
			continue
		}
		if curdelta < delta {
			delta = curdelta
			latest = f
		}
	}

	return latest, numTemplateFiles
}

func getDirFiles(dir string) ([]string, error) {
	fileNames := []string{}

	dirItems, err := os.ReadDir(dir)
	if err != nil {
		return fileNames, err
	}

	for _, item := range dirItems {
		if item.IsDir() {
			continue
		}
		fileNames = append(fileNames, item.Name())
	}

	return fileNames, nil
}

// func warnIfTooManyTemplateFiles(n int, thresh int) {
// 	if n > thresh {
// 		log.Printf("Found %d template files, consider running archive command for more efficient performance", n)
// 	}
// }
