# textnote
Simple tool for creating and organizing daily notes on the command line

[![Build Status](https://travis-ci.com/dkaslovsky/textnote.svg?branch=main)](https://travis-ci.com/github/dkaslovsky/textnote)
[![Coverage Status](https://coveralls.io/repos/github/dkaslovsky/textnote/badge.svg?branch=main)](https://coveralls.io/github/dkaslovsky/textnote?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/dkaslovsky/textnote)](https://goreportcard.com/report/github.com/dkaslovsky/textnote)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dkaslovsky/textnote/blob/main/LICENSE)

## Overview
textnote is a command line tool for quickly creating and managing daily plain text notes.
It is designed for ease of use to encourage the practice of daily, organized note taking.
textnote intentionally facilitates only the management (creation, opening, organizing, and consolidated archiving) of notes, following the philosophy that notes are best written in a text editor and not via a CLI.

Key features:
- Configurable, sectioned note template
- Easily bring content forward to the next day's note (for those to-dos that didn't quite get done today...)
- Simple command to consolidate daily notes into monthly archive files
- Create and open today's note with the default `textnote` command

All note files are stored locally on the file system in a single directory.
Notes can easily be synced to a remote server or cloud service if so desired by ensuring the application directory is remotely synced.

textnote opens notes using the text editor specified by the environment variable `$EDITOR` and defaults to Vim if the environment variable is not set.
See the [Editor-Specific Configuration](#editor-specific-configuration) subsection for more details. 

## Quick Start
1. Install textnote (see [Installation](#installation))
2. Set a single environment variable `TEXTNOTE_DIR` to specify the directory for textnote's files

That's it, textnote is ready to go!

The directory specified by `TEXTNOTE_DIR` and the default configuration file will be automatically created the first time textnote is run.

Start writing notes for today with a single command
```
$ textnote
```

To first configure textnote before creating notes, run
```
$ textnote init
```
and then edit the configuration file found at the displayed path.

## Installation
textnote can be installed by downloading a prebuilt binary or by the `go get` command.

### Releases
The recommended installation method is downloading the latest released binary.
Download the appropriate binary for your operating system from this repository's [releases](https://github.com/dkaslovsky/textnote/releases/latest) page or via `curl`:

macOS
```
$ curl -o textnote -L https://github.com/dkaslovsky/textnote/releases/latest/download/textnote_darwin_amd64
```

Linux
```
$ curl -o textnote -L https://github.com/dkaslovsky/textnote/releases/latest/download/textnote_linux_amd64
```

Windows
```
> curl.exe -o textnote.exe -L https://github.com/dkaslovsky/textnote/releases/latest/download/textnote_windows_amd64.exe
```

### Installing from source

textnote can also be installed using Go's built-in tooling:
```
$ go get -u github.com/dkaslovsky/textnote
```
Build from source by cloning this repository and running `go build`.

It is recommended to build using Go 1.15.7 or greater to avoid a potential security issue when looking for the desired editor in the `$PATH` ([details](https://blog.golang.org/path-security)).

## Usage
textnote is intentionally simple to use and supports two main commands: `open` for creating/opening notes and `archive` for consolidating notes into monthly archive files.

### **open**
The `open` command will open a dated note in an editor, creating it first if it does not exist.

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
To open a note for a specific date other than the current day, specify the date with the `--date` flag:
```
$ textnote open --date 2020-12-22
```
where the date format is specified in the configuration.

Alternatively, a note can be opened by passing the number of days prior to the current day using the `-d` flag. For example,
```
$ textnote open -d 1
```
opens yesterday's note.

Sections from previous notes can be copied or moved into a current note.
Each section to be copied is specified in a separate `-s` flag.
The most recent dated note is used as the source by default and a specific date for a source note can be provided through the `--copy` flag.
For example,
```
$ textnote open -s TODO -s NOTES
```
will create today's note with the "TODO" and "NOTES" sections copied from the most recently dated (often yesterday's) note, while
```
$ textnote open --copy 2021-01-17 -s TODO
```
creates today's note with the "TODO" section copied from the 2021-01-17 note.
Use the `-c` flag to instead specify the source by the number of days back from the current day.
For example,
```
$ textnote open -c 3 -s TODO
```
creates today's note with the "TODO" section copied from 3 days ago.

To move instead of copy, add the `-x` flag to any copy command.
For example,
```
$ textnote open --copy 2021-01-17 -s NOTES -x
```
moves the "NOTES" section contents from the 2021-01-17 note into the note for today.

Pass two delete flags (`-xx`) to also delete the source note if moving section(s) leaves the source empty:
```
$ textnote open --copy 2021-01-17 -s NOTES -xx
```

The `--date` and `--copy` (or `-d` and `-c`) flags can be used in combination if such a workflow is desired.

For convenience, the `-t` flag can be used to open tomorrow's note:
```
$ textnote open -t
```
For example,
```
$ textnote open -t -s TODO
```
creates a note for tomorrow with a copy of today's "TODO" section contents, assuming a note for today exits.

Also for convenience, the latest (most recent) dated note can be opened using the `-l` flag:
```
$ textnote open -l
```
The most recently dated note is typically from the previous day or a few days ago, but this command will return the note for the current date if it already exists.
It will ignore notes dated in the future.

When opening/copying requires searching for the latest (most recently dated) note, textnote checks the number of template files that were required to be searched.
If this number is above a threshold (as set in the [configuration](#configuration)), a message is displayed suggesting to run the [archive](#archive) command to reduce the number of template files.
This message can be effectively disabled by configuring the `templateFileCountThresh` configuration parameter to be very large, but doing so is not recommended.

The flag options are summarized by the command's help:
```
$ textnote open -h

open or create a note template

Usage:
  textnote open [flags]

Flags:
      --copy string       date of note for copying sections (defaults to date of most recent note, cannot be used with copy-back flag)
  -c, --copy-back uint    number of days back from today for copying from a note (cannot be used with copy flag)
      --date string       date for note to be opened (defaults to today)
  -d, --days-back uint    number of days back from today for opening a note (cannot be used with date, tomorrow, or latest flags)
  -x, --delete count      delete sections after copy (pass flag twice to also delete empty source note)
  -h, --help              help for open
  -l, --latest            specify the most recent dated note to be opened (cannot be used with date, days-back, or tomorrow flags)
  -s, --section strings   section to copy (defaults to none)
  -t, --tomorrow          specify tomorrow as the date for note to be opened (cannot be used with date, days-back, or latest flags)
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
with ellipses representing the daily notes' contents.

By default, the `archive` command is non-destructive: it will create archive files and leave all notes in place.
To delete the individual note files and retain only the generated archives, run the command with the `-x` flag:
```
$ textnote archive -x
```
This is the intended mode of operation, as it is desirable to "clean up" notes into archives, but must be intentionally enabled with `-x` for safety.
Running with the `--dry-run` flag prints the file names to be deleted without performing any actions:
```
$ textnote archive --dry-run
```

If the `archive` command is run without the delete flag, archive files are written and the original notes are left in place.
To "clean up" the original notes *after* archives have been generated, rerun the `archive` command with the `-x` flag as well as the `-n` flag to prevent duplicating the archive content:
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
  -x, --delete     delete individual files after archiving
      --dry-run    print file names to be deleted instead of performing deletes (other flags are ignored)
  -h, --help       help for archive
  -n, --no-write   disable writing archive files (helpful for deleting previously archived files)
```

### **Additional Functionality**
textnote is designed for simplicity. 
Because textnote writes files to a single directory on the local filesystem, most functionality outside of the scope described above can be easily accomplished using stanard command line tools (e.g., `grep` for search).

## Configuration
While textnote is intended to be extremely lightweight, it is also designed to be highly configurable.
In particular, the template (sections, headers, date formats, and whitespace) for generating notes can be customized as desired.
One might wish to configure headers and section titles for markdown compatibility or change date formats to match regional convention.

Configuration is read from the `$TEXTNOTE_DIR/.config.yml` file.
Changes to configuration parameters can be made by updating this file.
Individual configuration parameters also can be overridden with [environment variables](#environment-variable-overrides).

Importantly, if textnote's configuration is changed, notes created using a previous configuration might be incompatible with textnote's functionality.

The configuration file can be displayed by running the `config` command with the `-f` flag:
```
$ textnote config -f
```
The configuration file path is displayed by using the `-p` flag:
```
$ textnote config -p
```
[Defaults](#defaults) are used for configuration parameters omitted from the configuration file or configuration [environment variables](#environment-variable-overrides).
The `config` command with the `-a` flag displays the full "active" configuration used when the application runs, including default and environment parameters:
```
$ textnote config -a
```
To update the configuration file to match the active configuration, run
```
$ textnote config update
```
This command overwrites the existing configuration file.
It can be used instead of manual updates to the configuration file by passing environment variables.
For example,
```
$ TEXTNOTE_ARCHIVE_FILE_PREFIX="my_archive-" textnote config update
```
The `update` command is also helpful for writing configuration parameters that have been added with new versions of textnote.

The `config` command options are summarized by the command's help:
```
$ textnote config -h

manages the application's configuration

Usage:
  textnote config [flags]
  textnote config [command]

Available Commands:
  update      update the configuration file with active configuration

Flags:
  -a, --active   display configuration the application actively uses (includes environment variable configuration)
  -f, --file     display contents of configuration file (default)
  -h, --help     help for config
  -p, --path     display path to configuration file

Use "textnote config [command] --help" for more information about a command.
```

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
templateFileCountTresh: 90                # threshold for displaying a warning for too many template files
```

### Environment Variable Overrides
Any configuration parameter can be overridden by setting a corresponding environment variable.
Note that setting an environment variable does not change the value specified in the configuration file.
The full list of environment variables is listed below and is always available by running `textnote --help`:
```
  TEXTNOTE_TEMPLATE_FILE_COUNT_THRESH int
    	threshold for warning too many template files
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

### Editor-Specific Configuration
Currently, textnote supports the `file.cusorLine` and `TEXTNOTE_FILE_CURSOR_LINE` configuration for the following editors:
* Vi/Vim
* Emacs
* Neovim
* Nano

textnote will work with all other editors but will not respect this congifuration parameter.

## License
textnote is released under the [MIT License](https://github.com/dkaslovsky/textnote/blob/main/LICENSE).
Dependency licenses are available in this repository's [CREDITS](./CREDITS) file.
