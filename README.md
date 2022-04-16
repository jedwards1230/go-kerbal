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
 * Downloads and installs mods
 * Finds Kerbal game directory (prompts for directory if not found)
 * Sorts and filters mod list (by name and compatibility)
 * Displays logs in-app

## Images
![Main View](https://github.com/jedwards1230/go-kerbal/blob/main/screenshots/main.png?raw=true)
![Mod Selected](https://github.com/jedwards1230/go-kerbal/blob/main/screenshots/modInfo.png?raw=true)
![Log View](https://github.com/jedwards1230/go-kerbal/blob/main/screenshots/logs.png?raw=true)
![Input Directory View](https://github.com/jedwards1230/go-kerbal/blob/main/screenshots/inputDir.png?raw=true)

## TODO:
 * More metadata cleaning
 * Find mods currently installed on system
 * Check mod conflicts/dependencies 
 * Implement search box
 * Make the TUI prettier
 * Tweak the TUI
   * More colors
   * Add buttons
   * Multi-select mods
   * Progress bars
   * Live logging? (currently updates log view on TUI event, not when logs are called)
