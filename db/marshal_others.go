//go:build !release

package db

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var (
	Marshal   = json.Marshal
	Unmarshal = json.Unmarshal
)

var MarshalType = "json"
