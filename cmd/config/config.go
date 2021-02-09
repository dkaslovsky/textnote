package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// CreateConfigCmd creates the config subcommand
func CreateConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "show configuration",
		Long:  "displays the application's configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := filepath.Join(config.AppDir, config.FileName)
			_, err := os.Stat(configPath)
			if os.IsNotExist(err) {
				return fmt.Errorf("cannot find configuration file [%s]", configPath)
			}
			f, err := os.Open(configPath)
			if err != nil {
				return errors.Wrapf(err, "unable to open configuration file [%s]", configPath)
			}
			c, err := ioutil.ReadAll(f)
			if err != nil {
				return errors.Wrapf(err, "unable to read configuration file [%s]", configPath)
			}
			fmt.Printf("%s", c)
			return nil
		},
	}
	return cmd
}
