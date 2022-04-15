# go-kerbal (WIP)
 
This will eventually be a Golang clone of the Kerbal [CKAN Mod Manager](https://github.com/KSP-CKAN/CKAN). This uses the metadata files from [CKAN-Meta](https://github.com/KSP-CKAN/CKAN-meta).

## To run
```
go build && ./go-kerbal
```
## Features so far:
 * Automatically keeps metadata up to date
 * Compiles metadata and displays info in the TUI
 * Finds Kerbal game directory
 * Sorts and filters mod list
 * Displays logs in-app

## Images
![Main View](https://github.com/jedwards1230/go-kerbal/blob/main/screenshots/main.png?raw=true)
![Mod Selected](https://github.com/jedwards1230/go-kerbal/blob/main/screenshots/modInfo.png?raw=true)
![Log View](https://github.com/jedwards1230/go-kerbal/blob/main/screenshots/logs.png?raw=true)

## TODO:
 * More metadata cleaning
 * Find mods currently installed on system
 * Implement downloader
 * Implement installer
 * Implement search box
 * Keep adding to the TUI
 * And a lot more...
