# skogshuggare

Run with `go run . [vision_radius]` or build with `go build .` and then run with `./skogshuggare [vision_radius]`, e.g. `./skogshuggare` or `./skogshuggare 20`. The `vison_radius` argument specifies an integer value which defines the maximum distance from the player that is rendered on the map. If no argument is provided, a default of 100 is used.

## Maps
Maps are text files with the extension `.karta`. A map file must contain one player character and one squirrel character. Its boundaries must be defined with a rectangle of `#`. Within a map file, characters are defined as follows:

| Character  | Object           |
| :--------: | :--------------- |
| `p`        | Player           |
| `s`        | Squirrel         |
| `w`        | Water            |
| `W`        | Water, alternate |
| `f`        | Fire             |
| `#`        | Wall             |

### Example
This map file would create a 17x9 level with the player spawning at `(5, 2)`, the squirrel spawning at `(11, 5)`, a fire at `(12, 2)`, and four water tiles at `(3, 5), (4, 5), (3, 6), (4, 6)`. Coordinates are 0-indexed and the origin is in the top-left. `x` increases to the right and `y` increases toward the bottom.
```
#################
#               #
#    p      f   #
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