package open

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/template"
	"github.com/spf13/cobra"
)

// TODO: move this to config
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
			fileName := template.GetFileName(now)

			_, err := os.Stat(fileName)
			if os.IsNotExist(err) {
				body := makeNewBody(now)
				fo, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
				if err != nil {
					return err
				}
				err = body.Write(fo)
				if err != nil {
					return err
				}
			}
			return openInEditor(fileName, template.GetFirstSectionLine())
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

func openInEditor(fileName string, line int) error {
	lineArg := fmt.Sprintf("+%d", line)
	cmd := exec.Command("vim", lineArg, fileName)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
