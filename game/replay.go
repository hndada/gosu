package game

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/ulikunitz/xz/lzma"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

// 다 처리하면 스코어, 체력 시뮬레이터 만들기
// 키배치조합: x값, 1 2 4 8 16 32 64
type ReplayAction struct {
	w int64
	x float64
	y float64
	z int64
}

type LegacyReplay struct {
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
	ReplayData  []ReplayAction
	OnlineID    int64
	// AddMods
}

func main2() {
	const rDir = "../test/Replays/"
	rs, err := ioutil.ReadDir(rDir)
	if err != nil {
		panic(err)
	}
	for _, rp := range rs {
		func() {
			path := rDir + rp.Name()
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(path)
					fmt.Println(err)
				}
			}()
			_ = ReadLegacyReplay(path)
		}()
	}
}

func main() {
	r := ReadLegacyReplay("../test/od10 empty.osr")
	var time int64
	for _, rd := range r.ReplayData {
		time+=rd.w
		fmt.Printf("%d: %+v\n", time, rd)
	}
	// PrettyPrint(&r)
}

func ReadLegacyReplay(path string) *LegacyReplay {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var rp LegacyReplay
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
	rp.ReplayData = ReadReplayData(r)
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

func ReadReplayData(r io.Reader) (replayData []ReplayAction) {
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
	replayData = make([]ReplayAction, 0, len(actions))
	for _, f := range actions[:len(actions)-1] { // the stream ended with sep letter ","
		var ra ReplayAction
		vs := strings.Split(f, "|")
		if len(vs) != 4 {
			panic(vs)
		}

		ra.w, err = strconv.ParseInt(vs[0], 10, 64)
		if err != nil {
			panic(err)
		}
		ra.x, err = strconv.ParseFloat(vs[1], 64)
		if err != nil {
			panic(err)
		}
		ra.y, err = strconv.ParseFloat(vs[2], 64)
		if err != nil {
			panic(err)
		}
		ra.z, err = strconv.ParseInt(vs[3], 10, 64)
		if err != nil {
			panic(err)
		}
		if ra.w == -12345 {
			continue
		}
		replayData = append(replayData, ra)
	}
	return
}

// print the contents of the obj
func PrettyPrint(data interface{}) {
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}
