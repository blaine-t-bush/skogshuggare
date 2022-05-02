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
	MapCloud    = 'c'
	MapWall     = '#'
	MapFire     = 'f'
	// Growth chances (per game tick)
	BirdSpawnChance     = 0.950
	CloudSpawnChance    = 0.005
	GrowthChanceSeed    = 0.010 // Seed to sapling
	GrowthChanceSapling = 0.005 // Sapling to adult
	FireSpawnChance     = 0.005 // Chance per update for fire to randomly spawn on an available tile
	FireSpreadChance    = 0.900 // Chance per update for each fire to spread to a random adjacent tile
	FireBurnoutHalflife = 2000  // The age at which the chance (but not cumulative chance) for fire to burn out becomes 50%
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
	KeyCloud
	KeyBird
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
	AnimationStateGrassLight1
	AnimationStateGrassLight2
	AnimationStateGrassHeavy1
	AnimationStateGrassHeavy2
	AnimationStateWater1
	AnimationStateWater2
	AnimationStateFire1
	AnimationStateFire2
	AnimationStateFire3
	AnimationStateFire4
	AnimationStateCloud1
	AnimationStateCloud2
	AnimationStateBird1
	AnimationStateBird2
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

	animationRates = map[int]int{
		KeyGrassLight: 2,
		KeyGrassHeavy: 2,
		KeyWater:      4,
		KeyFire:       1,
		KeyCloud:      4,
		KeyBird:       1,
	}

	animationMarkov = map[int]map[int]AnimationMarkovNode{
		KeyGrassLight: {
			0: {AnimationStateGrassLight1, []AnimationMarkovConnection{
				{0, 0.90},
				{1, 0.10},
			}},
			1: {AnimationStateGrassLight2, []AnimationMarkovConnection{
				{0, 0.50},
				{1, 0.50},
			}},
		},
		KeyGrassHeavy: {
			0: {AnimationStateGrassHeavy1, []AnimationMarkovConnection{
				{0, 0.90},
				{1, 0.10},
			}},
			1: {AnimationStateGrassHeavy2, []AnimationMarkovConnection{
				{0, 0.50},
				{1, 0.50},
			}},
		},
		KeyWater: {
			0: {AnimationStateWater1, []AnimationMarkovConnection{
				{0, 0.90},
				{1, 0.10},
			}},
			1: {AnimationStateWater2, []AnimationMarkovConnection{
				{0, 0.50},
				{1, 0.50},
			}},
		},
		KeyFire: {
			0: {AnimationStateFire1, []AnimationMarkovConnection{
				{0, 0.10},
				{1, 0.70},
				{2, 0.10},
				{3, 0.10},
			}},
			1: {AnimationStateFire2, []AnimationMarkovConnection{
				{0, 0.10},
				{1, 0.10},
				{2, 0.70},
				{3, 0.10},
			}},
			2: {AnimationStateFire3, []AnimationMarkovConnection{
				{0, 0.10},
				{1, 0.10},
				{2, 0.10},
				{3, 0.70},
			}},
			3: {AnimationStateFire4, []AnimationMarkovConnection{
				{0, 0.70},
				{1, 0.10},
				{2, 0.10},
				{3, 0.10},
			}},
		},
		KeyCloud: {
			0: {AnimationStateCloud1, []AnimationMarkovConnection{
				{0, 0.70},
				{1, 0.30},
			}},
			1: {AnimationStateCloud2, []AnimationMarkovConnection{
				{0, 0.30},
				{1, 0.70},
			}},
		},
		KeyBird: {
			0: {AnimationStateBird1, []AnimationMarkovConnection{
				{0, 0.00},
				{1, 1.00},
			}},
			1: {AnimationStateBird2, []AnimationMarkovConnection{
				{0, 1.00},
				{1, 0.00},
			}},
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
		KeyBurnt:         {char: ' ', aboveActor: false, style: tcell.StyleDefault.Background(tcell.ColorDarkGray)},
		KeyFirebreak:     {char: ' ', aboveActor: false, style: tcell.StyleDefault.Background(tcell.ColorSandyBrown)},
		// Animated symbols
		AnimationStateGrassLight1: {char: '´', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorGreenYellow)},
		AnimationStateGrassLight2: {char: '`', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorGreenYellow)},
		AnimationStateGrassHeavy1: {char: '"', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorGreenYellow)},
		AnimationStateGrassHeavy2: {char: '”', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorGreenYellow)},
		AnimationStateWater1:      {char: ' ', aboveActor: false, style: tcell.StyleDefault.Background(tcell.ColorCornflowerBlue)},
		AnimationStateWater2:      {char: '~', aboveActor: false, style: tcell.StyleDefault.Foreground(tcell.ColorMediumBlue).Background(tcell.ColorCornflowerBlue)},
		AnimationStateFire1:       {char: '▓', aboveActor: true, style: tcell.StyleDefault.Foreground(tcell.ColorOrange).Background(tcell.ColorOrange)},
		AnimationStateFire2:       {char: '▓', aboveActor: true, style: tcell.StyleDefault.Foreground(tcell.ColorOrangeRed).Background(tcell.ColorOrange)},
		AnimationStateFire3:       {char: '▓', aboveActor: true, style: tcell.StyleDefault.Foreground(tcell.ColorOrange).Background(tcell.ColorOrange)},
		AnimationStateFire4:       {char: '▓', aboveActor: true, style: tcell.StyleDefault.Foreground(tcell.ColorOrange).Background(tcell.ColorOrangeRed)},
		AnimationStateCloud1:      {char: '▓', aboveActor: true, style: tcell.StyleDefault.Foreground(tcell.ColorLightGray).Background(tcell.ColorWhite)},
		AnimationStateCloud2:      {char: '▓', aboveActor: true, style: tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorLightGray)},
		AnimationStateBird1:       {char: '^', aboveActor: true, style: tcell.StyleDefault.Foreground(tcell.ColorWhite)},
		AnimationStateBird2:       {char: 'v', aboveActor: true, style: tcell.StyleDefault.Foreground(tcell.ColorWhite)},
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
