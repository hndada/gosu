package osr

import (
	"bytes"
	"encoding/binary"
	"github.com/ulikunitz/xz/lzma"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

// 키배치조합: x값, 1 2 4 8 16 32 64
type OsuReplayAction struct {
	W int64
	X float64
	Y float64
	Z int64
}

type OsuReplay struct {
	GameMode    int8
	GameVersion int32
	BeatmapMD5  string // [16]byte
	PlayerName  string
	ReplayMD5   string   // [16]byte
	Nums        [6]int16 // 300, 100, 50, 320, 200, Miss
	Score       int32
	Combo       int16
	FullCombo   int8 // bool
	ModsBits    int32
	LifeBar     string
	TimeStamp   int64
	ReplayData  []OsuReplayAction
	OnlineID    int64
	// AddMods
}

func ParseOsuReplay(path string) *OsuReplay {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var rp OsuReplay
	r := bytes.NewReader(b)
	binary.Read(r, binary.LittleEndian, &rp.GameMode)
	binary.Read(r, binary.LittleEndian, &rp.GameVersion)
	rp.BeatmapMD5 = ReadString(r)
	rp.PlayerName = ReadString(r)
	rp.ReplayMD5 = ReadString(r)
	binary.Read(r, binary.LittleEndian, &rp.Nums)
	binary.Read(r, binary.LittleEndian, &rp.Score)
	binary.Read(r, binary.LittleEndian, &rp.Combo)
	binary.Read(r, binary.LittleEndian, &rp.FullCombo)
	binary.Read(r, binary.LittleEndian, &rp.ModsBits)
	rp.LifeBar = ReadString(r)
	binary.Read(r, binary.LittleEndian, &rp.TimeStamp)
	rp.ReplayData = parseOsuReplayData(r)
	binary.Read(r, binary.LittleEndian, &rp.OnlineID)
	return &rp
}

func ReadString(r *bytes.Reader) string {
	first, err := r.ReadByte()
	if err != nil {
		panic(err)
	}
	switch first {
	case 0x00:
		return ""
	case 0x0b:
		strlen, err := binary.ReadUvarint(r)
		if err != nil {
			panic(err)
		}

		b := make([]byte, strlen)
		_, err = r.Read(b)
		if err != nil {
			panic(err)
		}
		return string(b)
	default:
		panic("not reach")
	}
}

func parseOsuReplayData(r io.Reader) (replayData []OsuReplayAction) {
	var replayDataLen int32
	binary.Read(r, binary.LittleEndian, &replayDataLen)

	var compReplayData = make([]byte, replayDataLen)
	n, err := r.Read(compReplayData)
	if err != nil {
		panic(err)
	}
	if int32(n) != replayDataLen {
		panic("error at parsing replay data")
	}

	cr, err := lzma.NewReader(bytes.NewReader(compReplayData))
	if err != nil {
		panic(err)
	}
	b := bytes.NewBuffer(make([]byte, 0, 10240))
	_, err = io.Copy(b, cr) // most stable way
	if err != nil {
		panic(err)
	}

	actions := strings.Split(b.String(), ",")
	replayData = make([]OsuReplayAction, 0, len(actions))
	for _, f := range actions[:len(actions)-1] { // the stream ended with sep letter ","
		var ra OsuReplayAction
		vs := strings.Split(f, "|")
		if len(vs) != 4 {
			panic(vs)
		}

		ra.W, err = strconv.ParseInt(vs[0], 10, 64)
		if err != nil {
			panic(err)
		}
		ra.X, err = strconv.ParseFloat(vs[1], 64)
		if err != nil {
			panic(err)
		}
		ra.Y, err = strconv.ParseFloat(vs[2], 64)
		if err != nil {
			panic(err)
		}
		ra.Z, err = strconv.ParseInt(vs[3], 10, 64)
		if err != nil {
			panic(err)
		}
		if ra.W == -12345 {
			continue
		}
		replayData = append(replayData, ra)
	}
	return
}
