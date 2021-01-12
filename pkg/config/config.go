package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	envAppDir      = "TEXTNOTE_DIR"
	configFileName = ".config.yml"
)

// AppDir is the directory in which the application stores its files
var AppDir = os.Getenv(envAppDir)

// Opts are options that configure the application
type Opts struct {
	Header  HeaderOpts  `yaml:"header"`
	Section SectionOpts `yaml:"section"`
	File    FileOpts    `yaml:"file"`
	Archive ArchiveOpts `yaml:"archive"`
}

// HeaderOpts are options for configuring the header of TextNote
type HeaderOpts struct {
	Prefix           string `yaml:"prefix" env:"TEXTNOTE_HEADER_PREFIX" env-description:"prefix to attach to header"`
	Suffix           string `yaml:"suffix" env:"TEXTNOTE_HEADER_SUFFIX" env-description:"suffix to attach to header"`
	TrailingNewlines int    `yaml:"trailingNewlines" env:"TEXTNOTE_HEADER_TRAILING_NEWLINES" env-description:"number of newlines to attach to end of header"`
	TimeFormat       string `yaml:"timeFormat" env:"TEXTNOTE_HEADER_TIME_FORMAT" env-description:"formatting string to form headers from timestamps"`
}

// SectionOpts are options for configuring sections of TextNote
type SectionOpts struct {
	Prefix           string   `yaml:"prefix" env:"TEXTNOTE_SECTION_PREFIX" env-description:"prefix to attach to section names"`
	Suffix           string   `yaml:"suffix" env:"TEXTNOTE_SECTION_SUFFIX" env-description:"suffix to attach to section names"`
	TrailingNewlines int      `yaml:"trailingNewlines" env:"TEXTNOTE_SECTION_TRAILING_NEWLINES" env-description:"number of newlines to attach to end of each section"`
	Names            []string `yaml:"names" env:"TEXTNOTE_SECTION_NAMES" env-description:"section names"`
}

// FileOpts are options for configuring files written by TextNote
type FileOpts struct {
	Ext        string `yaml:"ext" env:"TEXTNOTE_FILE_EXT" env-description:"extension for all files written"`
	TimeFormat string `yaml:"timeFormat" env:"TEXTNOTE_FILE_TIME_FORMAT" env-description:"formatting string to form file names from timestamps"`
	CursorLine int    `yaml:"cursorLine" env:"TEXTNOTE_FILE_CURSOR_LINE" env-description:"line to place cursor when opening"`
}

// ArchiveOpts are options for configuring archive files written by TextNote
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
	}
}

// LoadOrCreate loads a config file or creates it using defaults
func LoadOrCreate() (Opts, error) {
	err := EnsureAppDir()
	if err != nil {
		return Opts{}, err
	}

	configPath := filepath.Join(AppDir, configFileName)
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		defaults := getDefaultOpts()
		yml, err := yaml.Marshal(defaults)
		if err != nil {
			return Opts{}, errors.Wrap(err, "unable to generate config file")
		}
		err = ioutil.WriteFile(configPath, yml, 0644)
		if err != nil {
			return Opts{}, errors.Wrap(err, fmt.Sprintf("unable to create configuration file: [%s]", configPath))
		}
		log.Printf("created default configuration file: [%s]", configPath)
	}

	opts := Opts{}
	err = cleanenv.ReadConfig(configPath, &opts)
	if err != nil {
		return opts, errors.Wrap(err, "unable to read config file")
	}

	err = ValidateOpts(opts)
	if err != nil {
		return opts, errors.Wrapf(err, "configuration error in [%s]", configFileName)
	}

	return opts, nil
}

// EnsureAppDir validates that the application directory exists or is created
func EnsureAppDir() error {
	if AppDir == "" {
		return fmt.Errorf("required environment variable [%s] is not set", envAppDir)
	}

	finfo, err := os.Stat(AppDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(AppDir, 0755)
		if err != nil {
			return err
		}
		log.Printf("created directory [%s]", AppDir)
		return nil
	}

	if !finfo.IsDir() {
		return fmt.Errorf("[%s=%s] must be a directory", envAppDir, AppDir)
	}
	return nil
}

// ValidateOpts returns an error if the specified options are misconfigured
func ValidateOpts(opts Opts) error {
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

	// validate the file cursor line is not negative
	if opts.File.CursorLine < 0 {
		return errors.New("cursor line must not be negative")
	}

	return nil
}
