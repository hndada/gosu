package osr

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/ulikunitz/xz/lzma"
)

func Parse(dat []byte) (*Format, error) {
	var f Format
	r := bytes.NewReader(dat)
	var err error
	if err = binary.Read(r, binary.LittleEndian, &f.GameMode); err != nil {
		return &f, err
	}
	if err = binary.Read(r, binary.LittleEndian, &f.GameVersion); err != nil {
		return &f, err
	}
	if f.BeatmapMD5, err = readString(r); err != nil {
		return &f, err
	}
	if f.PlayerName, err = readString(r); err != nil {
		return &f, err
	}
	if f.ReplayMD5, err = readString(r); err != nil {
		return &f, err
	}
	if err = binary.Read(r, binary.LittleEndian, &f.Num300); err != nil {
		return &f, err
	}
	if err = binary.Read(r, binary.LittleEndian, &f.Num100); err != nil {
		return &f, err
	}
	if err = binary.Read(r, binary.LittleEndian, &f.Num50); err != nil {
		return &f, err
	}
	if err = binary.Read(r, binary.LittleEndian, &f.NumGeki); err != nil {
		return &f, err
	}
	if err = binary.Read(r, binary.LittleEndian, &f.NumKatu); err != nil {
		return &f, err
	}
	if err = binary.Read(r, binary.LittleEndian, &f.NumMiss); err != nil {
		return &f, err
	}
	if err = binary.Read(r, binary.LittleEndian, &f.Score); err != nil {
		return &f, err
	}
	if err = binary.Read(r, binary.LittleEndian, &f.Combo); err != nil {
		return &f, err
	}
	if err = binary.Read(r, binary.LittleEndian, &f.FullCombo); err != nil {
		return &f, err
	}
	if err = binary.Read(r, binary.LittleEndian, &f.ModsBits); err != nil {
		return &f, err
	}
	if f.LifeBar, err = readString(r); err != nil {
		return &f, err
	}
	if err = binary.Read(r, binary.LittleEndian, &f.TimeStamp); err != nil {
		return &f, err
	}
	if f.ReplayData, err = parseReplayData(r); err != nil {
		return &f, err
	}
	if err = binary.Read(r, binary.LittleEndian, &f.OnlineID); err != nil {
		return &f, err
	}
	return &f, nil
}

func readString(r *bytes.Reader) (string, error) {
	first, err := r.ReadByte()
	if err != nil {
		return "", err
	}
	switch first {
	case 0x00:
		return "", nil
	case 0x0b:
		length, err := binary.ReadUvarint(r)
		if err != nil {
			return "", err
		}
		b := make([]byte, length)
		if _, err = r.Read(b); err != nil {
			return "", err
		}
		return string(b), nil
	default:
		return "", errors.New("invalid replay file: corrupted string header")
	}
}

func parseReplayData(r io.Reader) ([]Action, error) {
	var length int32
	var err error
	empty := make([]Action, 0)
	if err = binary.Read(r, binary.LittleEndian, &length); err != nil {
		return empty, err
	}

	compressedData := make([]byte, length)
	n, err := r.Read(compressedData)
	if err != nil {
		return empty, err
	}
	if int32(n) != length {
		return empty, errors.New("invalid replay file: corrupted ReplayData length")
	}
	r2, err := lzma.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return empty, err
	}
	b := bytes.NewBuffer(make([]byte, 0, 10240))
	if _, err = io.Copy(b, r2); err != nil { // most stable way
		return empty, err
	}
	dat := strings.Split(b.String(), ",")

	actions := make([]Action, 0, len(dat)-1)
	for _, f := range dat[:len(dat)-1] { // the stream ended with sep letter ","
		var a Action
		vs := strings.Split(f, "|")
		if len(vs) != 4 {
			return actions, errors.New("invalid replay file: corrupted Action data; length is not 4")
		}
		if a.W == -12345 {
			continue
		}
		if a.W, err = strconv.ParseInt(vs[0], 10, 64); err != nil {
			return actions, errors.New("invalid replay file: W is corrupted")
		}
		if a.X, err = strconv.ParseFloat(vs[1], 64); err != nil {
			return actions, errors.New("invalid replay file: X is corrupted")
		}
		if a.Y, err = strconv.ParseFloat(vs[2], 64); err != nil {
			return actions, errors.New("invalid replay file: Y is corrupted")
		}
		if a.Z, err = strconv.ParseInt(vs[3], 10, 64); err != nil {
			return actions, errors.New("invalid replay file: Z is corrupted")
		}
		actions = append(actions, a)
	}
	return actions, nil
}
