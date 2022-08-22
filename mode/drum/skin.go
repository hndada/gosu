package drum

import "github.com/hndada/gosu/draws"

// https://osu.ppy.sh/wiki/en/Skinning/osu%21taiko
type Skin struct {
	ComboSprites []draws.Sprite
	ScoreSprites []draws.Sprite // Todo: move to /gosu

	JudgmentSprites [2][3]draws.Sprite
	KeySprites      [4]draws.Sprite
	FieldSprite     draws.Sprite
	BarLineSprite   draws.Sprite // Seperator of each bar (aka measure)

	DonSprites    [2][3]draws.Sprite
	KatSprites    [2][3]draws.Sprite
	HeadSprites   [2]draws.Sprite
	TailSprites   [2]draws.Sprite
	BodySprites   [2][]draws.Sprite // Binary-building method
	RollDotSprite draws.Sprite
	ShakeSprites  [3]draws.Sprite

	DancerSprites [4][]draws.Sprite
}

// [2] that most sprites have.
const (
	NormalNote = iota
	BigNote
)

// Significant keyword goes ahead, just as number is written: Left
const (
	KeyLeftKat  = iota
	KeyLeftDon  = iota
	KeyRightDon = iota
	KeyRightKat = iota
)
const (
	NoteGround = iota
	NoteOverlay1
	NoteOverlay2
)
const (
	ShakerNote = iota
	ShakerBottom
	ShakerTop
)
const (
	DancerIdle = iota
	DancerYes
	DancerNo
	DancerHigh
)
