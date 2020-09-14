package game

// time rate has been applied in advance
type BaseNote struct {
	Type  NoteType
	Time  int64
	Time2 int64 // ex) ln end time

	SampleVolume   uint8
	SampleFilename string
}

type NoteType int16
