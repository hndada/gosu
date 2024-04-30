package util

import (
	"crypto/md5"
	"encoding/hex"
)

// EncodeMD5 returns 32 hexadecimals, converted from [16]byte
// Decode v0:
// strconv.ParseUint(string(f.BeatmapMD5[i*2:(i+1)*2]), 16, 8) // byte
func MD5(data []byte) string {
	bs := md5.Sum(data)
	return hex.EncodeToString(bs[:])
}
