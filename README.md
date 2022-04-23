# skogshuggare

Run with `go run . [map_name]` or build with `go build .` and then run with `./skogshuggare [map_name]`, e.g. `./skogshuggare flod`. The argument specifies the name of a map file, not including extension, in the `./kartor/` to import and run. If no argument is provided, the program will attempt to use `./kartor/skog.karta`.

## Maps
Maps are text files with the extension `.karta`. A map file must contain one player character and one squirrel character. Within a map file, characters are defined as follows:

| Character | Object |
| :-: | :- |
|`p` | Player |
|`s` | Squirrel |
|`w` | Water |
|`#` | Wall |