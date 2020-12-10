package open

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

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
			fileName := template.GetFileName(date)
			err := createIfNotExists(date, fileName)
			if err != nil {
				return err
			}
			return openInEditor(fileName, template.FirstSectionFirstLine)
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
			fileName := template.GetFileName(date)
			err := createIfNotExists(date, fileName)
			if err != nil {
				return err
			}
			return openInEditor(fileName, template.FirstSectionFirstLine)
		},
	}
	return cmd
}

func createIfNotExists(date time.Time, fileName string) error {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		body := makeBody(date)
		fo, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		err = body.Write(fo)
		if err != nil {
			return err
		}
	}
	return nil
}

func makeBody(date time.Time) template.Body {
	sections := []*template.Section{}
	sectionNames := strings.Split(template.SectionNames, ",")
	for _, name := range sectionNames {
		sections = append(sections, template.NewSection(name))
	}
	body := template.NewBody(date, sections...)
	return body
}

func openInEditor(fileName string, line int) error {
	lineArg := fmt.Sprintf("+%d", line)
	cmd := exec.Command("vim", lineArg, fileName)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
