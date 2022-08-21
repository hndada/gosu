//go:build release

package db

import "github.com/vmihailenco/msgpack/v5"

var (
	Marshal   = msgpack.Marshal
	Unmarshal = msgpack.Unmarshal
)

var MarshalType = "msgpack"
