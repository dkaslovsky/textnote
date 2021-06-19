## 1.3.0 / 2021-06-16

* [ADDED] Second `--delete`/`-x` flag deletes source file left empty after moving section(s)

## 1.2.0 / 2021-04-26

* [ADDED] Flag to open most recently dated ("latest") note
* [ADDED] Configurable threshold for warning user of too many template files
* [ADDED] Flags to display configuration file contents (`-f`) and active configuration (`-a`)
* [ADDED] `update` subcommand for `config` command to overwrite configuration file with active configuration
* [ADDED] `init` command to more cleanly initialize textnote application directories and files
* [FIXED] Copy command defaults to latest note instead of potentially nonexistent note from previous day
* [INTERNAL] Upgraded to Go 1.16
* [INTERNAL] Deprecated use of `io/ioutil`

## 1.1.1 / 2021-02-28

* [FIXED] Fall back on defaults for parameters missing from configuration file
* [FIXED] Warning for unsupported editor configuration for cursorLine > 1

## 1.1.0 / 2021-02-16

* [ADDED] Use $EDITOR environment variable for opening notes
* [ADDED] Add support for vi/vim, nano, neovim, and emacs for using `file.cursorLine` config parameter

## 1.0.0 / 2021-02-09

* Initial release
