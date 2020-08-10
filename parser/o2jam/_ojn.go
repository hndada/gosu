package o2jam

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

