package osu

import (
	"image/color"
)

type Format struct {
	FormatVersion int `json:"formatVersion"`
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
	AudioFilename            string  `json:"audioFilename"`
	AudioLeadIn              int     `json:"audioLeadIn"`
	AudioHash                string  `json:"audioHash"` // deprecated
	PreviewTime              int     `json:"previewTime"`
	Countdown                int     `json:"countdown"` // nofloat
	SampleSet                string  `json:"sampleSet"`
	StackLeniency            float64 `json:"stackLeniency"`
	Mode                     int     `json:"mode"` // nofloat
	LetterboxInBreaks        bool    `json:"letterboxInBreaks"`
	StoryFireInFront         bool    `json:"storyFireInFront"` // deprecated
	UseSkinSprites           bool    `json:"useSkinSprites"`
	AlwaysShowPlayfield      bool    `json:"alwaysShowPlayfield"` // deprecated
	OverlayPosition          string  `json:"overlayPosition"`
	SkinPreference           string  `json:"skinPreference"`
	EpilepsyWarning          bool    `json:"epilepsyWarning"`
	CountdownOffset          int     `json:"countdownOffset"`
	SpecialStyle             bool    `json:"specialStyle"`
	WidescreenStoryboard     bool    `json:"widescreenStoryboard"`
	SamplesMatchPlaybackRate bool    `json:"samplesMatchPlaybackRate"`
}
type Editor struct { // delimiter:(space)
	Bookmarks       []int   // delimiter,
	DistanceSpacing float64 `json:"distanceSpacing"`
	BeatDivisor     float64 `json:"beatDivisor"`
	GridSize        int     `json:"gridSize"`
	TimelineZoom    float64 `json:"timelineZoom"`
}
type Metadata struct { // delimiter:
	Title         string   `json:"title"`
	TitleUnicode  string   `json:"titleUnicode"`
	Artist        string   `json:"artist"`
	ArtistUnicode string   `json:"artistUnicode"`
	Creator       string   `json:"creator"`
	Version       string   `json:"version"`
	Source        string   `json:"source"`
	Tags          []string // delimiter(space)
	BeatmapID     int      `json:"beatmapID"`    // nofloat
	BeatmapSetID  int      `json:"beatmapSetID"` // nofloat
}
type Difficulty struct { // delimiter:
	HPDrainRate       float64 `json:"hpDrainRate"`
	CircleSize        float64 `json:"circleSize"`
	OverallDifficulty float64 `json:"overallDifficulty"`
	ApproachRate      float64 `json:"approachRate"`
	SliderMultiplier  float64 `json:"sliderMultiplier"`
	SliderTickRate    float64 `json:"sliderTickRate"`
}
type Events []Event
type TimingPoints []TimingPoint
type Colours struct { // manual
	Combos              [8]color.RGBA
	SliderTrackOverride color.RGBA
	SliderBorder        color.RGBA
}
type HitObjects []HitObject
