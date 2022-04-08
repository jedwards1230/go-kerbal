# go-kerbal
 
Basically a TUI based Golang clone of the Kerbal [CKAN Mod Manager](https://github.com/KSP-CKAN/CKAN)

This uses the metadata files from [CKAN-Meta](https://github.com/KSP-CKAN/CKAN-meta)

## To run
```
go build && ./go-kerbal
```

## Currently:
 * Automatically keeps metadata up to date
 * Compiles metadata and displays info in the TUI
 * Finds Kerbal game directory
 * Sorts and filters mod list


## TODO:
 * Better .ckan data cleaning
 * Find currently installed mods
 * Implement installer
 * Search through list of mods
   * Add filtering for display
 * Make the TUI pretty
 * And a lot more...
