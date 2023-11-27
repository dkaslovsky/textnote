package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"dario.cat/mergo"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	// envAppDir is the name of the environment variable specifying the application directory
	envAppDir = "TEXTNOTE_DIR"
	// fileName is the name of the configuration file
	fileName = ".config.yml"
)

// appDir is the directory in which the application stores its files
var appDir = os.Getenv(envAppDir)

// Opts are options that configure the application
type Opts struct {
	AppDir                  string      `yaml:"-"` // AppDir is always read from the environment and is not written to file
	Header                  HeaderOpts  `yaml:"header"`
	Section                 SectionOpts `yaml:"section"`
	File                    FileOpts    `yaml:"file"`
	Archive                 ArchiveOpts `yaml:"archive"`
	Cli                     CliOpts     `yaml:"cli"`
	TemplateFileCountThresh int         `yaml:"templateFileCountThresh" env:"TEXTNOTE_TEMPLATE_FILE_COUNT_THRESH" env-description:"threshold for warning too many template files"`
}

// HeaderOpts are options for configuring the header of a note
type HeaderOpts struct {
	Prefix           string `yaml:"prefix" env:"TEXTNOTE_HEADER_PREFIX" env-description:"prefix to attach to header"`
	Suffix           string `yaml:"suffix" env:"TEXTNOTE_HEADER_SUFFIX" env-description:"suffix to attach to header"`
	TrailingNewlines int    `yaml:"trailingNewlines" env:"TEXTNOTE_HEADER_TRAILING_NEWLINES" env-description:"number of newlines to attach to end of header"`
	TimeFormat       string `yaml:"timeFormat" env:"TEXTNOTE_HEADER_TIME_FORMAT" env-description:"formatting string to form headers from timestamps"`
}

// SectionOpts are options for configuring sections of a note
type SectionOpts struct {
	Prefix           string   `yaml:"prefix" env:"TEXTNOTE_SECTION_PREFIX" env-description:"prefix to attach to section names"`
	Suffix           string   `yaml:"suffix" env:"TEXTNOTE_SECTION_SUFFIX" env-description:"suffix to attach to section names"`
	TrailingNewlines int      `yaml:"trailingNewlines" env:"TEXTNOTE_SECTION_TRAILING_NEWLINES" env-description:"number of newlines to attach to end of each section"`
	Names            []string `yaml:"names" env:"TEXTNOTE_SECTION_NAMES" env-description:"section names"`
}

// FileOpts are options for configuring file outputs
type FileOpts struct {
	Ext        string `yaml:"ext" env:"TEXTNOTE_FILE_EXT" env-description:"extension for all files written"`
	TimeFormat string `yaml:"timeFormat" env:"TEXTNOTE_FILE_TIME_FORMAT" env-description:"formatting string to form file names from timestamps"`
	CursorLine int    `yaml:"cursorLine" env:"TEXTNOTE_FILE_CURSOR_LINE" env-description:"line to place cursor when opening"`
}

// ArchiveOpts are options for configuring note archives
type ArchiveOpts struct {
	AfterDays                int    `yaml:"afterDays" env:"TEXTNOTE_ARCHIVE_AFTER_DAYS" env-description:"number of days after which to archive a file"`
	FilePrefix               string `yaml:"filePrefix" env:"TEXTNOTE_ARCHIVE_FILE_PREFIX" env-description:"prefix attached to the file name of all archive files"`
	HeaderPrefix             string `yaml:"headerPrefix" env:"TEXTNOTE_ARCHIVE_HEADER_PREFIX" env-description:"override header prefix for archive files"`
	HeaderSuffix             string `yaml:"headerSuffix" env:"TEXTNOTE_ARCHIVE_HEADER_SUFFIX" env-description:"override header suffix for archive files"`
	SectionContentPrefix     string `yaml:"sectionContentPrefix" env:"TEXTNOTE_ARCHIVE_SECTION_CONTENT_PREFIX" env-description:"prefix to attach to section content date"`
	SectionContentSuffix     string `yaml:"sectionContentSuffix" env:"TEXTNOTE_ARCHIVE_SECTION_CONTENT_SUFFIX" env-description:"suffix to attach to section content date"`
	SectionContentTimeFormat string `yaml:"sectionContentTimeFormat" env:"TEXTNOTE_ARCHIVE_SECTION_CONTENT_TIME_FORMAT" env-description:"formatting string dated section content"`
	MonthTimeFormat          string `yaml:"monthTimeFormat" env:"TEXTNOTE_ARCHIVE_MONTH_TIME_FORMAT" env-description:"formatting string for month archive timestamps"`
}

// CliOpts are options for configuring the CLI
type CliOpts struct {
	TimeFormat string `yaml:"timeFormat" env:"TEXTNOTE_CLI_TIME_FORMAT" env-description:"formatting string for timestamp CLI flags"`
}

// OptsBackCompat are options maintained for backwards compatibility that will be honored in the absence (zero-value) of their
// replacements as handled in loadBackCompat()
type OptsBackCompat struct {
	// TemplateFileCountThresh holds the value of the field "templateFileCountTresh" (note the typo) in a yaml configuration file
	TemplateFileCountThresh int `yaml:"templateFileCountTresh"`
}

func getDefaultOpts() Opts {
	return Opts{
		Header: HeaderOpts{
			Prefix:           "",
			Suffix:           "",
			TrailingNewlines: 1,
			TimeFormat:       "[Mon] 02 Jan 2006",
		},
		Section: SectionOpts{
			Prefix:           "___",
			Suffix:           "___",
			TrailingNewlines: 3,
			Names: []string{
				"TODO",
				"DONE",
				"NOTES",
			},
		},
		File: FileOpts{
			Ext:        "txt",
			TimeFormat: "2006-01-02",
			CursorLine: 4,
		},
		Archive: ArchiveOpts{
			AfterDays:                14,
			FilePrefix:               "archive-",
			HeaderPrefix:             "ARCHIVE ",
			HeaderSuffix:             "",
			SectionContentPrefix:     "[",
			SectionContentSuffix:     "]",
			SectionContentTimeFormat: "2006-01-02",
			MonthTimeFormat:          "Jan2006",
		},
		Cli: CliOpts{
			TimeFormat: "2006-01-02",
		},
		TemplateFileCountThresh: 90,
	}
}

// Load loads the configuration from file and/or evironment
func Load() (Opts, error) {
	opts := Opts{}

	// parse config file allowing environment variable overrides
	err := loadFromEnv(GetConfigFilePath(), &opts)
	if err != nil {
		return opts, fmt.Errorf("unable to read config file: %w", err)
	}

	// overwrite defaults with opts from file/env
	defaults := getDefaultOpts()
	err = mergo.Merge(&opts, defaults)
	if err != nil {
		return opts, fmt.Errorf("unable to integrate configuration from file with defaults: %w", err)
	}

	// set AppDir as read from environment
	opts.AppDir = appDir

	err = ValidateOpts(opts)
	if err != nil {
		return opts, fmt.Errorf("configuration error in [%s]: %w", fileName, err)
	}

	return opts, nil
}

func loadFromEnv(path string, opts *Opts) error {
	err := cleanenv.ReadConfig(path, opts)
	if err != nil {
		return err
	}

	err = loadBackCompat(path, opts)
	if err != nil {
		return fmt.Errorf("unable to read config file for backwards compatibility fields: %w", err)
	}

	return nil
}

func loadBackCompat(path string, opts *Opts) error {
	// TemplateFileCountThresh backwards compatibility with previously typo'd field
	if opts.TemplateFileCountThresh != 0 {
		return nil
	}
	backcompat := OptsBackCompat{}
	err := cleanenv.ReadConfig(GetConfigFilePath(), &backcompat)
	if err != nil {
		return err
	}
	opts.TemplateFileCountThresh = backcompat.TemplateFileCountThresh
	return nil
}

// CreateIfNotExists writes defaults to the configuration file if it does not already exist
func CreateIfNotExists() error {
	configPath := GetConfigFilePath()
	_, err := os.Stat(configPath)
	if !os.IsNotExist(err) {
		// config file exists, nothing to do
		return nil
	}

	defaults := getDefaultOpts()
	yml, err := yaml.Marshal(defaults)
	if err != nil {
		return fmt.Errorf("unable to generate config file: %w", err)
	}
	err = os.WriteFile(configPath, yml, 0o644)
	if err != nil {
		return fmt.Errorf("unable to create configuration file [%s]: %w", configPath, err)
	}
	log.Printf("created default configuration file: [%s]", configPath)
	return nil
}

// EnsureAppDir validates that the application directory exists or is created
func EnsureAppDir() error {
	if appDir == "" {
		return fmt.Errorf("required environment variable [%s] is not set", envAppDir)
	}

	finfo, err := os.Stat(appDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(appDir, 0o755)
		if err != nil {
			return err
		}
		log.Printf("created directory [%s]", appDir)
		return nil
	}

	if !finfo.IsDir() {
		return fmt.Errorf("[%s=%s] must be a directory", envAppDir, appDir)
	}
	return nil
}

// ValidateOpts returns an error if the specified options are misconfigured
func ValidateOpts(opts Opts) error {
	// validate appDir is not empty
	if opts.AppDir == "" {
		return fmt.Errorf("must include path to application directory in %s environment variable", envAppDir)
	}

	// validate at least one section
	if len(opts.Section.Names) == 0 {
		return errors.New("must include at least one section")
	}

	// validate section names are unique
	uniq := map[string]struct{}{}
	for _, name := range opts.Section.Names {
		uniq[name] = struct{}{}
	}
	if len(uniq) != len(opts.Section.Names) {
		return errors.New("section names must be unique")
	}

	// validate file archive prefix: this is needed for determining if a file is an archive
	if opts.Archive.FilePrefix == "" || strings.ReplaceAll(opts.Archive.FilePrefix, " ", "") == "" {
		return errors.New("file prefix for archives must not be empty")
	}

	// validate archive after days is at least 1
	if opts.Archive.AfterDays < 1 {
		return errors.New("archive after days must be greater than or equal to 1")
	}

	// validate file extension does not contain leading dot
	if strings.HasPrefix(opts.File.Ext, ".") {
		return errors.New("file extension must not include leading dot")
	}

	// validate the file cursor line is not negative
	if opts.File.CursorLine < 0 {
		return errors.New("cursor line must not be negative")
	}

	// validate threshold for warning on too many template files is larger than archive after days
	if opts.TemplateFileCountThresh <= opts.Archive.AfterDays {
		return errors.New("template file count threshold must be larger than archive after days")
	}

	return nil
}

// DescribeEnvVars returns a description string for environment variables used to configure the application
func DescribeEnvVars() string {
	header := ""
	description, err := cleanenv.GetDescription(&Opts{}, &header)
	if err != nil {
		return ""
	}
	return description
}

// GetConfigFilePath constructs the full path to the configuration file
func GetConfigFilePath() string {
	return filepath.Join(appDir, fileName)
}

// InitApp initializes the application by ensuring the necessary directories and files exist
func InitApp() error {
	err := EnsureAppDir()
	if err != nil {
		return err
	}
	err = CreateIfNotExists()
	if err != nil {
		return err
	}
	return nil
}
