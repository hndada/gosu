package osr

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"strconv"
	"strings"

	"github.com/ulikunitz/xz/lzma"
)

type Format struct {
	GameMode    int8
	GameVersion int32
	BeatmapMD5  string
	PlayerName  string
	ReplayMD5   string
	Num300      int16
	Num100      int16
	Num50       int16
	NumGeki     int16
	NumKatu     int16
	NumMiss     int16
	Score       int32
	Combo       int16
	FullCombo   bool
	ModsBits    int32
	LifeBar     string
	TimeStamp   int64
	ReplayData  []Action
	OnlineID    int64

	// AddMods is indirect data of accuracy at Target Practice.
	// It exists only when the mod is on.
	// AddMods     float64
}

type Action struct {
	// W is elapsed time since last action.
	W int64
	// X is x-axis mouse cursor position or pressed keys at osu!mania.
	// The least bit refers to state of the leftmost column and so on.
	X float64
	// Y is y-axis mouse cursor position.
	Y float64
	// Z is pressed keys at standard.
	Z int64
}

// NewFormat requires struct which implements the two following methods:
// func Read(b []byte) (n int, err error)
// func ReadByte() (byte, error)
// The simplest way is doing io.ReadAll then bytes.NewReader.
func NewFormat(file fs.File) (f *Format, err error) {
	f = &Format{}

	dat, err := io.ReadAll(file)
	if err != nil {
		return
	}
	r := bytes.NewReader(dat)

	if err = read(r, &f.GameMode); err != nil {
		return
	}
	if err = read(r, &f.GameVersion); err != nil {
		return
	}
	if err = readString(r, &f.BeatmapMD5); err != nil {
		return
	}
	if err = readString(r, &f.PlayerName); err != nil {
		return
	}
	if err = readString(r, &f.ReplayMD5); err != nil {
		return
	}
	if err = read(r, &f.Num300); err != nil {
		return
	}
	if err = read(r, &f.Num100); err != nil {
		return
	}
	if err = read(r, &f.Num50); err != nil {
		return
	}
	if err = read(r, &f.NumGeki); err != nil {
		return
	}
	if err = read(r, &f.NumKatu); err != nil {
		return
	}
	if err = read(r, &f.NumMiss); err != nil {
		return
	}
	if err = read(r, &f.Score); err != nil {
		return
	}
	if err = read(r, &f.Combo); err != nil {
		return
	}
	if err = read(r, &f.FullCombo); err != nil {
		return
	}
	if err = read(r, &f.ModsBits); err != nil {
		return
	}
	if err = readString(r, &f.LifeBar); err != nil {
		return
	}
	if err = read(r, &f.TimeStamp); err != nil {
		return
	}
	if err = readReplayData(r, &f.ReplayData); err != nil {
		return
	}
	if err = read(r, &f.OnlineID); err != nil {
		return
	}
	return f, nil
}

func read(r *bytes.Reader, dst any) error {
	return binary.Read(r, binary.LittleEndian, dst)
}

func readString(r *bytes.Reader, dst *string) error {
	first, err := r.ReadByte()
	if err != nil {
		return err
	}
	switch first {
	case 0x00:
		return nil
	case 0x0b:
		length, err := binary.ReadUvarint(r)
		if err != nil {
			return err
		}

		b := make([]byte, length)
		if _, err = r.Read(b); err != nil {
			return err
		}

		*dst = string(b)
		return nil
	default:
		return errors.New("invalid string header")
	}
}

func readReplayData(r io.Reader, dst *[]Action) error {
	var err error

	var length int32
	if err = binary.Read(r, binary.LittleEndian, &length); err != nil {
		return err
	}

	compressedData := make([]byte, length)
	n, err := r.Read(compressedData)
	if err != nil {
		return err
	}
	if int32(n) != length {
		return errors.New("replay data length mismatch")
	}

	r2, err := lzma.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(make([]byte, 0, 10240))
	if _, err = io.Copy(b, r2); err != nil { // most stable way
		return err
	}
	dat := strings.Split(b.String(), ",")

	actions := make([]Action, 0, len(dat)-1)
	for _, f := range dat[:len(dat)-1] { // the stream ended with sep letter ","
		var a Action
		vs := strings.Split(f, "|")
		if len(vs) != 4 {
			return errors.New("action data length is not 4")
		}
		if a.W, err = strconv.ParseInt(vs[0], 10, 64); err != nil {
			return errors.New("failed to parse w")
		}
		// if a.W == -12345 {
		// 	continue
		// }
		if a.X, err = strconv.ParseFloat(vs[1], 64); err != nil {
			return errors.New("failed to parse x")
		}
		if a.Y, err = strconv.ParseFloat(vs[2], 64); err != nil {
			return errors.New("failed to parse y")
		}
		if a.Z, err = strconv.ParseInt(vs[3], 10, 64); err != nil {
			return errors.New("failed to parse z")
		}
		actions = append(actions, a)
	}
	*dst = actions
	return nil
}

// for i := range md5 {...}
// strconv.ParseUint(string(f.BeatmapMD5[i*2:(i+1)*2]), 16, 8)
func (f Format) MD5() (hash [16]byte, err error) {
	var hashBytes []byte
	hashBytes, err = hex.DecodeString(f.BeatmapMD5)
	if err != nil {
		return
	}
	copy(hash[:], hashBytes)
	return
}

// In normal replay, first 2 data are dummy with x = 256 and y = -500
// I assume it is for setting time offset: -1.
// In auto replay, first data is blank action.
// func (f Format) IsAuto() bool {
// 	const (
// 		x = 256
// 		y = -500
// 	)
// 	if len(f.ReplayData) < 2 {
// 		return true
// 	}
// 	a0, a1 := f.ReplayData[0], f.ReplayData[1]
// 	if a0.X == 0 && a0.Y == 0 {
// 		return true
// 	} else if a0.X == x && a0.Y == y && a1.X == x && a1.Y == y {
// 		return false
// 	}
// 	panic("no reach")
// }

// // Last action data is dummy which is for random seed.
// func (f Format) TrimmedActions() []Action {
// 	if f.IsAuto() {
// 		return f.ReplayData
// 	}
// 	return f.ReplayData[2 : len(f.ReplayData)-1]
// }

// BufferTime returns the amount of time of waiting before music start when playing a chart.
// func (f Format) BufferTime() int64 {
// 	if f.IsAuto() {
// 		return 0
// 	}
// 	a0, a1, a2 := f.ReplayData[0], f.ReplayData[1], f.ReplayData[2]
// 	return a0.W + a1.W + a2.W // Must be 0 - 1 - (actual buffer time)
// }
