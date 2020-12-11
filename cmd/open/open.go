package open

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/config"
	"github.com/dkaslovsky/TextNote/pkg/template"
	"github.com/spf13/cobra"
)

// CreateTodayCmd creates the subcommand to open today's note
func CreateTodayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "today",
		Short: "open today's note",
		Long:  "open a text based note template for today",
		RunE: func(cmd *cobra.Command, args []string) error {
			date := time.Now()
			return run(date)
		},
	}
	return cmd
}

// CreateTomorrowCmd creates the subcommand to open tomorrow's note
func CreateTomorrowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tomorrow",
		Short: "open tomorrow's note",
		Long:  "open a text based note template for tomorrow",
		RunE: func(cmd *cobra.Command, args []string) error {
			date := time.Now().Add(24 * time.Hour)
			return run(date)
		},
	}
	return cmd
}

func run(date time.Time) error {
	opts, err := config.LoadOrCreate()
	if err != nil {
		return err
	}

	t := template.NewTemplate(opts)
	t.SetDate(date)
	err = createFileIfNotExists(t)
	if err != nil {
		return err
	}
	return openInEditor(t.GetFilePath(), t.GetFirstSectionFirstLine())
}

func createFileIfNotExists(t *template.Template) error {
	fileName := t.GetFilePath()
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		fo, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		err = t.Write(fo)
		if err != nil {
			return err
		}
	}
	return nil
}

func openInEditor(fileName string, line int) error {
	lineArg := fmt.Sprintf("+%d", line)
	cmd := exec.Command("vim", lineArg, fileName)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
