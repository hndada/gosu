package db

// Header contains non-play information.
// Changing Header's data will not affect integrity of the chart.
// Mode-specific fields are located to each Chart struct.
type Header struct {
	ChartSetID    int64 // Compatibility for osu.
	ChartID       int64 // Todo: ChartID -> ID
	MusicName     string
	MusicUnicode  string
	Artist        string
	ArtistUnicode string
	MusicSource   string
	ChartName     string
	Charter       string
	HolderID      int64

	PreviewTime     int64
	MusicFilename   string // Filename is fine to use (cf. Filepath)
	ImageFilename   string
	VideoFilename   string
	VideoTimeOffset int64

	Mode    int
	SubMode int
}

// func (c Header) MusicPath(cpath string) (string, bool) {
// 	if name := c.MusicFilename; name == "virtual" || name == "" {
// 		return "", false
// 	}
// 	return filepath.Join(filepath.Dir(cpath), c.MusicFilename), true
// }
// func (c Header) BackgroundPath(cpath string) string {
// 	return filepath.Join(filepath.Dir(cpath), c.ImageFilename)
// }
