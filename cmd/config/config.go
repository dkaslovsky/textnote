package config

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dkaslovsky/textnote/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type commandOptions struct {
	path   bool
	active bool
	file   bool
}

// CreateConfigCmd creates the config subcommand
func CreateConfigCmd() *cobra.Command {
	cmdOpts := commandOptions{}
	cmd := &cobra.Command{
		Use:   "config",
		Short: "manage configuration",
		Long:  "manages the application's configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := config.GetConfigFilePath()

			if cmdOpts.path {
				log.Printf("configuration file path: [%s]", configPath)
				return nil
			}

			if cmdOpts.active {
				return displayActiveConfig()
			}

			// default
			return displayConfigFile(configPath)
		},
	}
	attachOpts(cmd, &cmdOpts)
	cmd.AddCommand(CreateConfigUpdateCmd())
	return cmd
}

func attachOpts(cmd *cobra.Command, cmdOpts *commandOptions) {
	flags := cmd.Flags()
	flags.BoolVarP(&cmdOpts.path, "path", "p", false, "display path to configuration file")
	flags.BoolVarP(&cmdOpts.active, "active", "a", false, "display configuration the application actively uses (includes environment variable configuration)")
	flags.BoolVarP(&cmdOpts.file, "file", "f", false, "display contents of configuration file (default)")
}

// CreateConfigUpdateCmd creates the config update subcommand
func CreateConfigUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update the configuration file with active configuration",
		Long:  "update the configuration file to match the active configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			active, err := getActiveConfigYaml()
			if err != nil {
				return err
			}
			return os.WriteFile(config.GetConfigFilePath(), active, 0o644)
		},
	}
	return cmd
}

func displayConfigFile(configPath string) error {
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("cannot find configuration file [%s]", configPath)
	}
	f, err := os.Open(configPath)
	if err != nil {
		return errors.Wrapf(err, "unable to open configuration file [%s]", configPath)
	}
	c, err := io.ReadAll(f)
	if err != nil {
		return errors.Wrapf(err, "unable to read configuration file [%s]", configPath)
	}
	log.Print(string(c))
	return nil
}

func displayActiveConfig() error {
	yml, err := getActiveConfigYaml()
	if err != nil {
		return err
	}
	log.Print(string(yml))
	return nil
}

func getActiveConfigYaml() ([]byte, error) {
	opts, err := config.Load()
	if err != nil {
		return []byte{}, err
	}
	return yaml.Marshal(opts)
}
