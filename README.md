# skogshuggare

Run with `go run . [map_name [vision_radius]]` or build with `go build .` and then run with `./skogshuggare [map_name [vision_radius]]`, e.g. `./skogshuggare flod` or `./skogshuggare ö 10`. The `map_name` argument specifies the name of a map file, not including extension, in the `./kartor/` directory to import and run. If no argument is provided, the program will attempt to use `skog` located in `./kartor/skog.karta`. The `vison_radius` argument specifies an integer value which defines the maximum distance from the player that is rendered on the map. If no argument is provided, a default of 20 is used.

## Maps
Maps are text files with the extension `.karta`. A map file must contain one player character and one squirrel character. Its boundaries must be defined with a rectangle of `#`. Within a map file, characters are defined as follows:

| Character  | Object           |
| :--------: | :--------------- |
| `p`        | Player           |
| `s`        | Squirrel         |
| `w`        | Water            |
| `W`        | Water, alternate |
| `#`        | Wall             |

### Example
This map file would create a 17x9 level with the player spawning at `(5, 2)`, the squirrel spawning at `(11, 5)`, and four water tiles at `(3, 5), (4, 5), (3, 6), (4, 6)`. Coordinates are 0-indexed and the origin is in the top-left. `x` increases to the right and `y` increases toward the bottom.
```
#################
#               #
#    p          #
#               #
#               #
#  ww      s    #
#  ww           #
#               #
#################
```

## About
### Authors
- [Blaine Bush](https://github.com/blaine-t-bush)
- [blackscalare](https://github.com/blackscalare)

### License

### Name
[Skogshuggare](https://sv.wikipedia.org/wiki/Skogshuggare) is Swedish for lumberjack or logger. *Skog* means forest, from Proto-Germanic *\*skōgaz*, and *hugga* means hew, from Proto-Germanic *\*hawwaną*, cognate with English *hew*. Thus the most literal translation to English is *forest's hewer*.