package open

import (
	"os"
	"os/exec"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/template"
	"github.com/spf13/cobra"
)

// TODO: move this to config, calculate the +4
var defaultSectionNames = []string{
	"TODO",
	"DONE",
	"NOTES",
}

func CreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "open today's note",
		Long:  "open a text based note template for today",
		RunE: func(cmd *cobra.Command, args []string) error {
			now := time.Now()
			fileName := template.GetFileNameFromTime(now)

			_, err := os.Stat(fileName)
			if !os.IsNotExist(err) {
				return openInEditor(fileName)
			}
			if err != nil {
				return err
			}

			body := makeNewBody(now)
			fo, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				return err
			}
			err = body.Write(fo)
			if err != nil {
				return err
			}
			return openInEditor(body.GetFileName())
		},
	}
	return cmd
}

func makeNewBody(date time.Time) template.Body {
	sections := []*template.Section{}
	for _, name := range defaultSectionNames {
		sections = append(sections, template.NewSection(name))
	}
	body := template.NewBody(date, sections...)
	return body
}

func openInEditor(fileName string) error {
	cmd := exec.Command("vim", "+4", fileName)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
