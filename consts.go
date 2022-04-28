package main

import "github.com/gdamore/tcell"

const (
	// Game parameters
	StateTickDuration     = 30  // Milliseconds between game state update ticks
	AnimationTickDuration = 200 // Milliseconds between animation update ticks
	MaxIterations         = 1000
	// Map characters
	MapPlayer   = 'p'
	MapSquirrel = 's'
	MapWater    = 'w'
	MapWall     = '#'
	MapFire     = 'f'
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
	KeyWater
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
	// Title menu states
	MainMenuPageOrder
	NewGamePageOrder
	// DifficultyPageOrder
	// Animation states
	AnimationStateWater1
	AnimationStateWater2
	AnimationStateFire1
	AnimationStateFire2
	AnimationStateFire3
	AnimationStateFire4
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

	animationStates = map[int]map[int]int{
		KeyWater: { // animationStage: animationState
			0: AnimationStateWater1,
			1: AnimationStateWater1,
			2: AnimationStateWater1,
			3: AnimationStateWater1,
			4: AnimationStateWater2,
		},
		KeyFire: { // animationStage: animationState
			0: AnimationStateFire1,
			1: AnimationStateFire2,
			2: AnimationStateFire3,
			3: AnimationStateFire4,
		},
	}

	symbols = map[int]Symbol{ // Color options are listed at https://github.com/gdamore/tcell/blob/master/color.go
		// Static symbols
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
		KeyWater:         {char: ' ', aboveActor: false, style: tcell.StyleDefault.Background(tcell.ColorCornflowerBlue)},
		KeyBurnt:         {char: '▓', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorDarkSlateGray).Background(tcell.ColorDarkGray)},
		KeyFirebreak:     {char: '▓', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorSandyBrown)},
		// Animated symbols
		AnimationStateWater1: {char: ' ', aboveActor: false, style: tcell.StyleDefault.Background(tcell.ColorCornflowerBlue)},
		AnimationStateWater2: {char: '~', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorMediumBlue).Background(tcell.ColorCornflowerBlue)},
		AnimationStateFire1:  {char: '▓', aboveActor: true, style: tcell.StyleDefault.Foreground(tcell.ColorOrange).Background(tcell.ColorOrange)},
		AnimationStateFire2:  {char: '▓', aboveActor: true, style: tcell.StyleDefault.Foreground(tcell.ColorOrangeRed).Background(tcell.ColorOrange)},
		AnimationStateFire3:  {char: '▓', aboveActor: true, style: tcell.StyleDefault.Foreground(tcell.ColorOrange).Background(tcell.ColorOrange)},
		AnimationStateFire4:  {char: '▓', aboveActor: true, style: tcell.StyleDefault.Foreground(tcell.ColorOrange).Background(tcell.ColorOrangeRed)},
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
