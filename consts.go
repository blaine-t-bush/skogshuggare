package main

import "github.com/gdamore/tcell"

const (
	// Game parameters
	TickRate                  = 30 // Milliseconds between ticks
	MaxIterations             = 1000
	DefaultGeneratedMapWidth  = 20
	DefaultGeneratedMapHeight = 20
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
	FireSpawnChance     = 0.005 // Chance per update for fire to randomly spawn on an available tile
	FireSpreadChance    = 0.100 // Chance per update for each fire to spread to a random adjacent tile
	FireBurnoutHalflife = 200   // The age at which the chance (but not cumulative chance) for fire to burn out becomes 50%
	// Fire and hitpoints
	MaxHitPointsPlayer   = 3
	MaxHitPointsSquirrel = 1
	DamageFire           = 1
	FireWeightedDistance = 20
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
	KeyFireType1
	KeyFireType2
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
	// Title menu states
	MainMenuPageOrder
	NewGamePageOrder
	GenerateMapPageOrder
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
		KeyFireType1:     {char: '▓', aboveActor: true, style: tcell.StyleDefault.Foreground(tcell.ColorOrange).Background(tcell.ColorOrangeRed)},
		KeyFireType2:     {char: '▓', aboveActor: true, style: tcell.StyleDefault.Foreground(tcell.ColorOrangeRed).Background(tcell.ColorOrange)},
		KeyBurnt:         {char: '▓', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorDarkSlateGray).Background(tcell.ColorDarkGray)},
		KeyFirebreak:     {char: '▓', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorSandyBrown)},
	}
)

// NOTE
// Interesting Unicode characters (e.g. arrows) start at 2190.

// Title menu header
const TitleMenuHeaderAnim1 = `SKOGSHUGGARE


`
const TitleMenuHeaderAnim2 = `sKOGSHUGGARE


`
const TitleMenuHeaderAnim3 = `SkOGSHUGGARE


`
const TitleMenuHeaderAnim4 = `SKoGSHUGGARE


`
const TitleMenuHeaderAnim5 = `SKOgSHUGGARE


`
const TitleMenuHeaderAnim6 = `SKOGsHUGGARE


`
const TitleMenuHeaderAnim7 = `SKOGShUGGARE


`
const TitleMenuHeaderAnim8 = `SKOGSHuGGARE


`
const TitleMenuHeaderAnim9 = `SKOGSHUgGARE


`
const TitleMenuHeaderAnim10 = `SKOGSHUGgARE


`
const TitleMenuHeaderAnim11 = `SKOGSHUGGaRE


`
const TitleMenuHeaderAnim12 = `SKOGSHUGGArE


`

const TitleMenuHeaderAnim13 = `SKOGSHUGGARe


`
