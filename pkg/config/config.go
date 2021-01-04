package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	envAppDir      = "TEXTNOTE_DIR"
	configFileName = ".config.yml"
)

var (
	appDir = os.Getenv(envAppDir)
)

// Opts are options that configure the application
type Opts struct {
	AppDir  string      `yaml:"appdir,omitempty"`
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
	TimeFormat string `yaml:"timeFormat" env:"TEXTNOTE_FILE_TIME_FORMAT" env-description:"formatting string to form file names from timestamps"`
	Ext        string `yaml:"ext" env:"TEXTNOTE_FILE_EXT" env-description:"extension for all files written"`
}

// ArchiveOpts are options for configuring archive files written by TextNote
type ArchiveOpts struct {
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
			TimeFormat: "2006-01-02",
			Ext:        "txt",
		},
		Archive: ArchiveOpts{
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
	opts := Opts{}

	if appDir == "" {
		return opts, fmt.Errorf("environment variable [%s] is not set", envAppDir)
	}

	configPath := filepath.Join(appDir, configFileName)
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		defaults := getDefaultOpts()
		yml, err := yaml.Marshal(defaults)
		if err != nil {
			return opts, errors.Wrap(err, "unable to generate config file")
		}
		err = ioutil.WriteFile(configPath, yml, 0644)
		if err != nil {
			return opts, errors.Wrap(err, fmt.Sprintf("unable to create configuration file: [%s]", configPath))
		}
		log.Printf("created default configuration file: [%s]", configPath)
	}

	err = cleanenv.ReadConfig(configPath, &opts)
	if err != nil {
		return opts, errors.Wrap(err, "unable to read config file")
	}

	opts.AppDir = appDir

	err = ValidateConfig(opts)
	if err != nil {
		return opts, errors.Wrapf(err, "configuration error in [%s]", configFileName)
	}

	return opts, nil
}

// EnsureAppDir validates that the application directory exists or is created
func EnsureAppDir() error {
	if appDir == "" {
		return fmt.Errorf("required environment variable [%s] is not set", envAppDir)
	}
	finfo, err := os.Stat(appDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(appDir, 0755)
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

// ValidateConfig returns an error if the specified options are misconfigured
func ValidateConfig(opts Opts) error {
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

	// validate header time format
	t := time.Now()
	_, err := time.Parse(opts.Header.TimeFormat, t.Format(opts.Header.TimeFormat))
	if err != nil {
		return errors.New("invalid header time format")
	}
	_, err = time.Parse(opts.File.TimeFormat, t.Format(opts.File.TimeFormat))
	if err != nil {
		return errors.New("invalid file time format")
	}
	_, err = time.Parse(opts.Archive.SectionContentTimeFormat, t.Format(opts.Archive.SectionContentTimeFormat))
	if err != nil {
		return errors.New("invalid archive section content time format")
	}
	_, err = time.Parse(opts.Archive.MonthTimeFormat, t.Format(opts.Archive.MonthTimeFormat))
	if err != nil {
		return errors.New("invalid archive month time format")
	}

	return nil
}
