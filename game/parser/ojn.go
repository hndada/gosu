package parser

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/text/encoding"
	"io/ioutil"
)

const (
	ballad = iota
	rock
	dance
	techno
	hiphop
	soulRNB
	jazz
	funk
	classical
	traditional
	etc
)

// 본섭곡, 아무 최근 자작곡 한개로 테스트
// byte 초과 길이로 저장하려 하면 31바이트 이후는 잘림
type OJN struct {
	SongID           int32
	Signature        [4]byte // "OJN\0"
	EncodeVersion    float32 // Encoder value (9A 99 39 40)
	Genre            int32
	BPM              float32
	Level            [4]int16 // Last 2 bytes are unused
	EventCount       [3]int32 // including bg music notes
	NoteCount        [3]int32 // without bg music notes
	MeasureCount     [3]int32
	BlockCount       [3]int32
	OldEncodeVersion int16 // 29
	OldSongID        int16
	OldGenre         [20]byte
	BMPSize          int32
	OldFileVersion   int32
	Title            [64]byte
	Artist           [32]byte
	Noter            [32]byte
	OJMFile          [32]byte
	CoverSize        int32
	Time             [3]int32
	NoteOffset       [3]int32
	CoverOffset      int32
	_                int32
}

func ParseOJN(path string) *OJN {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var ojn OJN
	// r1 := bytes.NewReader(b)
	// r2 := encoding.Reader(r1)
	var r encoding.Decoder
	// korean 인코딩 추가
	binary.Read(r, binary.LittleEndian, &ojn)
	return &ojn
}

// wrong format, don't use this.
type ojn2 struct { // another header specification
	SongID           int32
	Signature        [4]byte // "OJN\0"
	EncodeVersion    float32 // Encoder value (9A 99 39 40)
	Genre            int32
	BPM              float32
	Level            [4]int16 // Last 2 bytes are unused
	EventCount       [3]int32 // including bg music notes
	NoteCount        [3]int32 // without bg music notes
	MeasureCount     [3]int32
	BlockCount       [3]int32
	OldEncodeVersion int16 // 29
	OldSongID        int32
	// OldGenre         [20]byte
	// BMPSize          int32
	// OldFileVersion   int32
	Title       [58]byte
	Artist      [64]byte
	Noter       [32]byte
	OJMFile     [32]byte
	CoverSize   int32
	Time        [3]int32
	NoteOffset  [4]int32
	CoverOffset int32
	_           int32
}

func readOjn2(path string) *ojn2 {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var ojn ojn2
	r := bytes.NewReader(b)
	binary.Read(r, binary.LittleEndian, &ojn)
	return &ojn
}
