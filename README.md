# go-kerbal (WIP)
 
This will eventually be a Go clone of the Kerbal [CKAN Mod Manager](https://github.com/KSP-CKAN/CKAN). This uses the metadata files from [CKAN-Meta](https://github.com/KSP-CKAN/CKAN-meta).

I use CKAN every time I play KSP. This project is mostly to get comfortable with Go. I would like to have this functional and ready to port for KSP 2 whenever that is released.

## To run
```
go build && ./go-kerbal
```
## Features so far:
 * Automatically keeps metadata up to date
 * Compiles metadata and displays info in the TUI
 * Sorts by:
   * name
 * Filters by:
   * compatibility
   * if installed (TODO)
   * tag (TODO)
 * Displays logs in-app
 * Downloads and installs mods
 * Finds Kerbal game directory (prompts for directory if not found)
 * Finds installed mods

## Images
![Main View](https://github.com/jedwards1230/go-kerbal/blob/main/screenshots/main.png?raw=true)
![Mod Selected](https://github.com/jedwards1230/go-kerbal/blob/main/screenshots/modInfo.png?raw=true)
![Log View](https://github.com/jedwards1230/go-kerbal/blob/main/screenshots/logs.png?raw=true)
![Input Directory View](https://github.com/jedwards1230/go-kerbal/blob/main/screenshots/inputDir.png?raw=true)

## TODO:
 * KSP repair script
 * Reliable way to find current KSP version
 * More precise install functions/validation
 * More metadata cleaning
 * Check mod conflicts/dependencies 
 * Implement search box
 * Tweak the TUI
   * More colors
   * Multi-select mods
   * Progress bars
   * Live logging? (currently updates log view on TUI event, not when logs are called)
