# textnote
Simple tool for creating and organizing daily notes on the command line

[![Build Status](https://travis-ci.com/dkaslovsky/textnote.svg?branch=main)](https://travis-ci.com/github/dkaslovsky/textnote)
[![Coverage Status](https://coveralls.io/repos/github/dkaslovsky/textnote/badge.svg?branch=main)](https://coveralls.io/github/dkaslovsky/textnote?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/dkaslovsky/textnote)](https://goreportcard.com/report/github.com/dkaslovsky/textnote)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dkaslovsky/textnote/blob/main/LICENSE)

## Overview
textnote is a command line tool for quickly creating and managing dated, plain text notes.
It is designed for ease of use in order to encourage the practice of daily, organized note taking.
textnote intentionally facilitates only the management (creation, opening, organizing, and consolidated archiving) of notes, following the philosophy that notes are best written in a text editor and not via a CLI.

Key features:
- Configurable, sectioned note template
- Easily bring content forward to the next day's note (for those to-dos that didn't quite get done...)
- Simple command to consolidate daily notes into monthly archive files
- Create and open today's note with the default `textnote` command

All note files are stored locally on the file system in a single directory.
Notes can easily be synced to a remote server or cloud service if so desired by ensuring the application directory is remotely synced.

Currently, Vim is the only supported text editor.

## Quick Start
1. Install textnote (see [Installation](#installation))
2. Set a single environment variable `TEXTNOTE_DIR` to specify the directory for textnote's files
3. Start writing notes for today with a single command
```
$ textnote
```
That's it! The directory specified by `TEXTNOTE_DIR` and the default configuration file will be automatically created.

## Installation
### Releases
The recommended installation method is downloading the latest released binary for your operating system.
Download the appropriate binary from this repository's [releases](https://github.com/dkaslovsky/textnote/releases/latest) page or via `curl`:

macOS
```
$ curl -OL https://github.com/dkaslovsky/textnote/releases/latest/download/textnote_darwin_amd64
```

Linux
```
$ curl -OL https://github.com/dkaslovsky/textnote/releases/latest/download/textnote_linux_amd64
```

Windows
```
$ curl -OL https://github.com/dkaslovsky/textnote/releases/latest/download/textnote_windows_amd64
```

### Installing from source

textnote can also be installed using Go's built-in tooling:
```
$ go get -u github.com/dkaslovsky/textnote
```
Build from source by cloning this repository and running `go build`.

It is recommended to build using Go 1.15.7 or greater to avoid a potential security issue when looking for the desired editor in the `$PATH` ([details](https://blog.golang.org/path-security)).

## Usage
textnote is intentionally simple to use and supports two commands: `open` for creating/opening notes and `archive` for consolidating notes into monthly archive files.

### **open**
The `open` command will open a specified note in an editor, creating it first if it does not exist.

Opening or creating a note for the current day is the default action.
Simply run the root command to open or create a note for the current day:
```
$ textnote
```
which, using the default configuration and assuming today is 2021-01-24, will create and open an empty note template:
```
[Sun] 24 Jan 2021

___TODO___



___DONE___



___NOTES___



```
To open a note for a specific date other than the current day, specify the date with the `-d` flag:
```
$ textnote open -d 2020-12-22
```
where the date format is specified in the configuration or by passing a Golang time format through the `-f` flag:
```
$ textnote open -d "Dec 12 2020" -f "Jan 02 2006"
```
For convenience, yesterday's note can be opened by passing the `-y` flag:
```
$ textnote open -y
```

Sections from previous notes can be copied or moved into a current note.
Each section to be copied is specified in separate `-s` flags.
The previous day's note is used as the source by default and a specified date for a source note can be provided through the `-c` flag.
For example,
```
$ textnote open -s TODO -s NOTES
```
will create today's note with the "TODO" and "NOTES" sections copied from yesterday's note, while
```
$ textnote open -c 2021-01-17 -s NOTES
```
creates today's note with the "NOTES" section copied from the 2021-01-17 note.

To move instead of copy, add the `-x` flag to any copy command.
For example,
```
$ textnote open -c 2021-01-17 -s NOTES -x
```
moves the "NOTES" section contents from the 2021-01-17 note into the note for today.

The flag options are summarized by the command's help:
```
$ textnote open -h

open or create a note template

Usage:
  textnote open [flags]

Flags:
  -c, --copyDate string   date of note for copying sections (defaults to yesterday)
  -d, --date string       date for note to be opened (defaults to today)
  -x, --delete            delete sections after copy
  -f, --format string     override for time format to parse date flags specified in configuration
  -h, --help              help for open
  -s, --section strings   section to copy
  -y, --yesterday         use yesterday's date for note (ignored if date is specified)
```


## **archive**
The `archive` command consolidates all daily notes into month archives, gathering together the contents for each section of a month in chronological order, labeled by the original date.
Only notes older than a number of days specified in the configuration are archived.

Running the archive command
```
$ textnote archive
```
generates an archive file for every month for which a note exists.
For example, an archive of the January 2021 notes, assuming the default configuration, will have the form
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
with ellipses representing the notes' contents.

By default, the `archive` command is non-destructive: it will create archive files and leave all archived notes in place.
To delete the individual note files and retain only the generated archives, run the command with the `-x` flag:
```
$ textnote archive -x
```
This is the intended mode of operation, as it is desirable to "clean up" notes into archives, but must be intentionally enabled with `-x` for safety.
To "clean up" *after* archives have been generated, rerun the `archive` command with the `-x` flag as well as the `-n` flag to prevent duplicating the archive content:
```
$ textnote archive -x -n
```

The flag options are summarized by the command's help:
```
$ textnote archive -h

consolidate notes into monthly archive files

Usage:
  textnote archive [flags]

Flags:
  -x, --delete    delete individual files after archiving
  -h, --help      help for archive
  -n, --nowrite   disable writing archive file (helpful for deleting previously archived files)
```

## Configuration
While textnote is intended to be extremely lightweight, it is also designed to be highly configurable.
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

### Environment Variable Overrides
Any configuration parameter can be overridden by setting a corresponding environment variable.
Note that setting an environment variable does not change the value specified in the configuration file.
The full list of environment variables is listed below and is always available by running `textnote --help`:
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
