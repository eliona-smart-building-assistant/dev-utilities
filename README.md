# Utilities for Eliona app structures

This package contains utilities to help with Eliona app development.

## Asset from struct

Utility `asset-from-struct` generates asset type JSON from provided structure.

The structure provided must have `eliona` and `subtype` field tags defined by the Eliona field tag syntax (TODO: link to doc).

### Usage

Run the utility and pipe or paste the struct to `stdin`. The resulting JSON will be printed out to `stdout`.

## Struct from asset

This utility was developed originally when there was a need to pass to Eliona specific structures separately by subtype. This method is deprecated now.

`struct-from-asset` generates the subtype-separated structs for all app's assets and prints them out to `stdout`

### Usage

Pass the path to the folder with asset type definitions (typically `app-root/eliona`) as the first parameter. The utility will take as input all files like `asset_type_*.json` from the folder.
