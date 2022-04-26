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
	MapFire       = 'f'
	// Growth chances (per game tick)
	GrowthChanceSeed    = 0.010 // Seed to sapling
	GrowthChanceSapling = 0.005 // Sapling to adult
	FireSpreadChance    = 0.100 // Chance per update for each fire to spread to a random adjacent tile
	FireBurnoutHalflife = 200   // The age at which the chance (but not cumulative chance) for fire to burn out becomes 50%
	// Hit points
	MaxHitPointsPlayer   = 3
	MaxHitPointsSquirrel = 1
	DamageFire           = 1
	// Actors
	ActorPlayer = iota
	ActorSquirrel
	ContentCategoryTerrain
	ContentCategoryDecoration
	ContentCategoryFire
	ContentCategoryTree
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
	KeyFire
	KeyBurnt
	KeyFirebreak
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
		KeyPlayer:        {char: '@', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorIndianRed)},
		KeySquirrel:      {char: 'ơ', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorRosyBrown)},
		KeyWall:          {char: '#', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorWhite)},
		KeyTreeSeed:      {char: '.', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorKhaki)},
		KeyTreeSapling:   {char: '┃', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorDarkKhaki)},
		KeyTreeTrunk:     {char: '█', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorSaddleBrown)},
		KeyTreeLeaves:    {char: '▓', aboveActor: true, style: tcell.StyleDefault.Foreground(tcell.ColorForestGreen)},
		KeyTreeStump:     {char: '▄', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorSaddleBrown)},
		KeyTreeStumpling: {char: '╻', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorDarkKhaki)},
		KeyGrassLight:    {char: '\'', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorGreenYellow)},
		KeyGrassHeavy:    {char: '"', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorGreenYellow)},
		KeyWaterLight:    {char: ' ', aboveActor: false, style: tcell.StyleDefault.Background(tcell.ColorCornflowerBlue)},
		KeyWaterHeavy:    {char: '~', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorMediumBlue).Background(tcell.ColorCornflowerBlue)},
		KeyFire:          {char: '▓', aboveActor: true, style: tcell.StyleDefault.Foreground(tcell.ColorOrangeRed)},
		KeyBurnt:         {char: '▓', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorDarkSlateGray)},
		KeyFirebreak:     {char: '▓', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorSandyBrown)},
	}
)

// NOTE
// Interesting Unicode characters (e.g. arrows) start at 2190.
