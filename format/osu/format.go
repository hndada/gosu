package osu

import "image/color"

type Format struct {
	FormatVersion int // delimiter:(space)
	General
	Editor
	Metadata
	Difficulty
	Events
	TimingPoints
	Colours
	HitObjects
}

type General struct { // delimiter:(space)
	AudioFilename            string
	AudioLeadIn              int
	AudioHash                string // deprecated
	PreviewTime              int
	Countdown                int // nofloat
	SampleSet                string
	StackLeniency            float64
	Mode                     int // nofloat
	LetterboxInBreaks        bool
	StoryFireInFront         bool // deprecated
	UseSkinSprites           bool
	AlwaysShowPlayfield      bool // deprecated
	OverlayPosition          string
	SkinPreference           string
	EpilepsyWarning          bool
	CountdownOffset          int
	SpecialStyle             bool
	WidescreenStoryboard     bool
	SamplesMatchPlaybackRate bool
}

type Editor struct { // delimiter:(space)
	Bookmarks       []int // delimiter,
	DistanceSpacing float64
	BeatDivisor     float64
	GridSize        int
	TimelineZoom    float64
}

type Metadata struct { // delimiter:
	Title         string
	TitleUnicode  string
	Artist        string
	ArtistUnicode string
	Creator       string
	Version       string
	Source        string
	Tags          []string // delimiter(space)
	BeatmapID     int      // nofloat
	BeatmapSetID  int      // nofloat
}

type Difficulty struct { // delimiter:
	HPDrainRate       float64
	CircleSize        float64
	OverallDifficulty float64
	ApproachRate      float64
	SliderMultiplier  float64
	SliderTickRate    float64
}

type Events []Event

type TimingPoints []TimingPoint

type Colours struct { // manual
	Combos              [8]color.RGBA
	SliderTrackOverride color.RGBA
	SliderBorder        color.RGBA
}

type HitObjects []HitObject
