# go-kerbal
 
Basically a TUI based Golang clone of the Kerbal [CKAN Mod Manager](https://github.com/KSP-CKAN/CKAN)

This uses the metadata files from [CKAN-Meta](https://github.com/KSP-CKAN/CKAN-meta)

## To run
```
./go-kerbal
```

## Currently:
 * Pulls metadata and reads .ckan files
 * Compiles mods into a list for display
 * Displays a GUI with the list of mods


## TODO:
 * Optimize how .ckan files are saved into Modules and the Registry
 * Handle old versions of mods
 * Find kerbal game directory
   * Find currently installed mods
 * Implement mod installer
 * Search through list of mods
   * Add filtering for display
 * Efficiently save metadata repo and only update as needed
 * Make the GUI pretty
 * And a lot more...
