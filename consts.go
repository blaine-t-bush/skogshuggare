package main

import "github.com/gdamore/tcell"

const (
	// Game parameters
	TickRate = 30 // Milliseconds between ticks
	// Map characters
	MapPlayer     = 'p'
	MapSquirrel   = 's'
	MapWaterLight = 'w'
	MapWaterHeavy = 'W'
	MapWall       = '#'
	// Growth chances (per game tick)
	GrowthChanceSeed    = 0.010 // Seed to sapling
	GrowthChanceSapling = 0.005 // Sapling to adult
	SeedCreationChance  = 0.005 // Seed spawning
	SeedCreationMax     = 3     // Maximum number of seeds to create per tick
	// Actors
	ActorPlayer = iota
	ActorSquirrel
	// Keys for accessing properties of various symbols
	KeyPlayer
	KeySquirrel
	KeyWall
	KeyTreeSeed
	KeyTreeSapling
	KeyTreeTrunk
	KeyTreeLeaves
	KeyTreeStump
	KeyTreeStumpling
	KeyGrassLight
	KeyGrassHeavy
	KeyWaterLight
	KeyWaterHeavy
	// Directions
	DirUp
	DirRight
	DirDown
	DirLeft
	DirOmni
	DirRandom
	DirNone
	// Living tree states
	TreeStateSeed
	TreeStateSapling
	TreeStateAdult
	// Harvested tree states
	TreeStateRemoved
	TreeStateStump
	TreeStateTrunk
	TreeStateStumpling
	// Border states
	TopBorder
	RightBorder
	BottomBorder
	LeftBorder
	TopLeftCorner
	TopRightCorner
	BottomRightCorner
	BottomLeftCorner
	// Title menu states
	MainMenuPageOrder
	NewGamePageOrder
	// DifficultyPageOrder
)

var (
	treeGrowingStages = map[int]GrowthInfo{ // For a given state (key), gives the next growth state (value.newState) and chance [0-1] of the next growth stage (value.chance)
		TreeStateSeed:    {newState: TreeStateSapling, chance: GrowthChanceSeed},
		TreeStateSapling: {newState: TreeStateAdult, chance: GrowthChanceSeed},
	}

	treeHarvestingStages = map[int]int{ // For a given state (key), gives the next harvesting state (value)
		TreeStateSeed:      TreeStateRemoved,
		TreeStateSapling:   TreeStateStumpling,
		TreeStateStumpling: TreeStateRemoved,
		TreeStateAdult:     TreeStateTrunk,
		TreeStateTrunk:     TreeStateStump,
		TreeStateStump:     TreeStateRemoved,
	}

	symbols = map[int]Symbol{ // Color options are listed at https://github.com/gdamore/tcell/blob/master/color.go
		KeyPlayer:        {char: '@', style: tcell.StyleDefault.Foreground(tcell.ColorIndianRed)},
		KeySquirrel:      {char: 's', style: tcell.StyleDefault.Foreground(tcell.ColorRosyBrown)},
		KeyWall:          {char: '#', style: tcell.StyleDefault.Foreground(tcell.ColorWhite)},
		KeyTreeSeed:      {char: '.', style: tcell.StyleDefault.Foreground(tcell.ColorKhaki)},
		KeyTreeSapling:   {char: '┃', style: tcell.StyleDefault.Foreground(tcell.ColorDarkKhaki)},
		KeyTreeTrunk:     {char: '█', style: tcell.StyleDefault.Foreground(tcell.ColorSaddleBrown)},
		KeyTreeLeaves:    {char: '▓', style: tcell.StyleDefault.Foreground(tcell.ColorForestGreen)},
		KeyTreeStump:     {char: '▄', style: tcell.StyleDefault.Foreground(tcell.ColorSaddleBrown)},
		KeyTreeStumpling: {char: '╻', style: tcell.StyleDefault.Foreground(tcell.ColorDarkKhaki)},
		KeyGrassLight:    {char: '\'', style: tcell.StyleDefault.Foreground(tcell.ColorGreenYellow)},
		KeyGrassHeavy:    {char: '"', style: tcell.StyleDefault.Foreground(tcell.ColorGreenYellow)},
		KeyWaterLight:    {char: ' ', style: tcell.StyleDefault.Background(tcell.ColorCornflowerBlue)},
		KeyWaterHeavy:    {char: '~', style: tcell.StyleDefault.Foreground(tcell.ColorMediumBlue).Background(tcell.ColorCornflowerBlue)},
	}
)

// NOTE
// Interesting Unicode characters (e.g. arrows) start at 2190.

// Title menu header

const TitleMenuHeaderAnim1 = `
SKOGSHUGGARE


`
const TitleMenuHeaderAnim2 = `
sKOGSHUGGARE


`
const TitleMenuHeaderAnim3 = `
SkOGSHUGGARE


`
const TitleMenuHeaderAnim4 = `
SKoGSHUGGARE


`
const TitleMenuHeaderAnim5 = `
SKOgSHUGGARE


`
const TitleMenuHeaderAnim6 = `
SKOGsHUGGARE


`
