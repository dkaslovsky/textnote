# textnote
Simple tool for creating and organizing daily notes on the command line

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dkaslovsky/textnote/blob/main/LICENSE)

## Overview
textnote is a command line tool built to quickly create and manage dated, plain text notes.
It is designed for ease of use in order to encourage the practice of daily, organized note taking.
textnote intentionally facilitates only the management (creation, opening, organizing, and consolidated archiving) of notes, following the philosophy that notes are best written in a text editor and not via a CLI.

Key features include
- configurable, sectioned note template
- easily bring content forward to the next day's note (for those to-dos that didn't quite get done...)
- simple command to consolidate daily notes into monthly archive files
- create or open today's note with a single, default command

Currently, Vim is the only supported text editor.

## Quick Start
1. [Install](#installation) textnote
2. Set a single environment variable `TEXTNOTE_DIR` to specify the directory for textnote's files
3. Start writing notes for today with a single command
```
$ textnote
```
That's it! The directory specified by `TEXTNOTE_DIR` and the default configuration file will be automatically created.

## Installation
textnote can be installed using Go's built in tooling:
```
$ go get -u github.com/dkaslovsky/textnote
```
and can of course be built from source by cloning the repository and running `go build`

## Configuration
While textnote is intended to be extremely lightweight, it is also designed to be configurable.
In particular, the template (sections, headers, date formats, and whitespace) for generating notes can be customized as desired.
Configuration is read from the `$TEXTNOTE_DIR/.config.yml` file and individual configuration parameters can be overridden with [environment variables](#environment-variable-overrides).

### Defaults
The default configuration file is automatically written the first time textnote is run:
```
header:
  prefix: ""                              # prefix to attach to header
  suffix: ""                              # suffix to attach to header
  trailingNewlines: 1                     # number of newlines after header
  timeFormat: '[Mon] 02 Jan 2006'         # Golang format for header dates
section:
  prefix: ___                             # prefix to attach to section name
  suffix: ___                             # suffix to attach to section name
  trailingNewlines: 3                     # number of newlines for empty section
  names:                                  # section names
  - TODO
  - DONE
  - NOTES
file:
  ext: txt                                # extension to use for note files
  timeFormat: "2006-01-02"                # Golang format for note file names
  cursorLine: 4                           # line to place cursor when opening a note
archive:
  afterDays: 14                           # number of days after which a note can be archived
  filePrefix: archive-                    # prefix to attach to archive file names
  headerPrefix: 'ARCHIVE '                # prefix to attach to header of archive notes
  headerSuffix: ""                        # suffix to attach to header of archive notes
  sectionContentPrefix: '['               # prefix to attach to section content date
  sectionContentSuffix: ']'               # suffix to attach to section content date
  sectionContentTimeFormat: "2006-01-02"  # Golang format for section content dates
  monthTimeFormat: Jan2006                # Golang format for month archive file and header dates
cli:
  timeFormat: "2006-01-02"                # Golang format for CLI date input
```

The default configuration produces the following note template
```
[Sun] 24 Jan 2021

___TODO___



___DONE___



___NOTES___



```
and the following archive template (ellipses added to represent content)
```
ARCHIVE Jan2021

___TODO___
[2021-01-03]
...
[2021-01-04]
...



___DONE___
[2021-01-03]
...
[2021-01-04]
...
[2021-01-06]
...


___NOTES___
[2021-01-06]
...



```

### Environment Variable Overrides
Any configuration paramter can be overridden by setting a corresponding environment variable.
Note that setting an environment variable does not change the value specified in the configuration file.
The environment variables follow the convention of upper case, underscore separators, and are prefixed with "TEXTNOTE".
The full list is always available from the CLI help `textnote --help`:
```
  TEXTNOTE_HEADER_PREFIX string
    	prefix to attach to header
  TEXTNOTE_HEADER_SUFFIX string
    	suffix to attach to header
  TEXTNOTE_HEADER_TRAILING_NEWLINES int
    	number of newlines to attach to end of header
  TEXTNOTE_HEADER_TIME_FORMAT string
    	formatting string to form headers from timestamps
  TEXTNOTE_SECTION_PREFIX string
    	prefix to attach to section names
  TEXTNOTE_SECTION_SUFFIX string
    	suffix to attach to section names
  TEXTNOTE_SECTION_TRAILING_NEWLINES int
    	number of newlines to attach to end of each section
  TEXTNOTE_SECTION_NAMES slice
    	section names
  TEXTNOTE_FILE_EXT string
    	extension for all files written
  TEXTNOTE_FILE_TIME_FORMAT string
    	formatting string to form file names from timestamps
  TEXTNOTE_FILE_CURSOR_LINE int
    	line to place cursor when opening
  TEXTNOTE_ARCHIVE_AFTER_DAYS int
    	number of days after which to archive a file
  TEXTNOTE_ARCHIVE_FILE_PREFIX string
    	prefix attached to the file name of all archive files
  TEXTNOTE_ARCHIVE_HEADER_PREFIX string
    	override header prefix for archive files
  TEXTNOTE_ARCHIVE_HEADER_SUFFIX string
    	override header suffix for archive files
  TEXTNOTE_ARCHIVE_SECTION_CONTENT_PREFIX string
    	prefix to attach to section content date
  TEXTNOTE_ARCHIVE_SECTION_CONTENT_SUFFIX string
    	suffix to attach to section content date
  TEXTNOTE_ARCHIVE_SECTION_CONTENT_TIME_FORMAT string
    	formatting string dated section content
  TEXTNOTE_ARCHIVE_MONTH_TIME_FORMAT string
    	formatting string for month archive timestamps
  TEXTNOTE_CLI_TIME_FORMAT string
    	formatting string for timestamp CLI flags
```

## Usage
### open
### archive
