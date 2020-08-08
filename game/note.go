package game

type BaseNote struct {
	Type int16
	Time     int64
	Time2    int64

	SampleVolume   uint8
	SampleFilename string
}

// mania.Note 다시 짜기
// mania.beatmap->chart 다시 짜기
// lv 위한 Note 어떻게하지
// play 코드 구축
// lv 대강