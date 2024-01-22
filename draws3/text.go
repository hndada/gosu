package draws

import (
	"bytes"
	"io"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

// Font: Face source; a file which implements Face.
// Face: A style; aka TypeFace. (e.g.: Arial Bold)
// An implementation of different sizes of the same face are generated on the fly.

// type Font = *text.GoTextFaceSource
// type Face = text.Face

// TrueType vs OpenType: OpenType is an extension of TrueType.
var cachedFaceSources = make(map[string]*text.GoTextFaceSource)

// src which implements the following methods can be loaded efficiently:
// Read([]byte) (int, error)
// ReadAt([]byte, int64) (int, error)
// Seek(int64, int) (int64, error)
// bytes.NewReader(b) implements io.ReadSeeker, as well as io.ReaderAt.
func NewFaceSource(src io.ReadSeeker) (*text.GoTextFaceSource, error) {
	font, err := text.NewGoTextFaceSource(src)
	if err != nil {
		return nil, err
	}
	return font, nil
}

func NewFaceSourceFromFile(fsys fs.FS, name string) (*text.GoTextFaceSource, error) {
	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewFaceSource(f.(io.ReadSeeker))
}

var cachedFaces = make(map[FaceOptions]text.Face)

type FaceOptions struct {
	Font string // Face source name
	opentype.FaceOptions
}

func NewFaceOptions() FaceOptions {
	return FaceOptions{
		Font: "goregular",
		FaceOptions: opentype.FaceOptions{
			Size:    12,
			DPI:     72,
			Hinting: font.HintingFull,
		},
	}
}

func NewFace(opts FaceOptions) (text.Face, error) {
	face, ok := cachedFaces[opts]
	if ok {
		return face, nil
	}

	src, ok := cachedFaceSources[opts.Font]
	if !ok {
		src = cachedFaceSources["goregular"]
	}

	face = &text.GoTextFace{
		Source: src,
		Size:   opts.Size,
	}
	cachedFaces[opts] = face
	return face, nil
}

func init() {
	defaultFont, _ := NewFaceSource(bytes.NewReader(goregular.TTF))
	cachedFaceSources["goregular"] = defaultFont

	// Load default face.
	defaultFaceOptions := NewFaceOptions()
	NewFace(defaultFaceOptions)
}

type Text struct {
	Text string
	FaceOptions
	face text.Face
	Box
	LineSpacing float64
}

func NewText(txt string) Text {
	t := Text{
		Text:        txt,
		LineSpacing: 1.6,
	}
	t.SetFace(NewFaceOptions())
	t.Box = NewBox(t)
	return t
}

func (t *Text) SetFace(opts FaceOptions) error {
	face, err := NewFace(opts)
	t.FaceOptions = opts
	t.face = face
	return err
}

func (t Text) IsEmpty() bool { return len(t.Text) == 0 }

func (t Text) Size() XY {
	x, y := text.Measure(t.Text, t.face, t.LineSpacing)
	return XY{x, y}
}

func (t Text) Draw(dst Image) {
	if t.IsEmpty() {
		return
	}
	src := t.Text
	op := &text.DrawOptions{
		DrawImageOptions: *t.op(),
		LayoutOptions: text.LayoutOptions{
			LineSpacingInPixels: t.LineSpacing,
		},
	}
	text.Draw(dst.Image, src, t.face, op)
}

// // E.g.: gomonobolditalic
// func (k FontOptions) String() string {
// 	var (
// 		mono   string
// 		weight string
// 		style  string
// 	)
// 	if k.Mono {
// 		mono = "mono"
// 	}
// 	switch k.Weight {
// 	case font.WeightThin: // -3
// 		weight = "thin"
// 	case font.WeightExtraLight:
// 		weight = "ultralight"
// 	case font.WeightLight:
// 		weight = "light"
// 	case font.WeightNormal: // 0; CSS: 400
// 		weight = ""
// 	case font.WeightMedium:
// 		weight = "medium"
// 	case font.WeightSemiBold:
// 		weight = "semibold"
// 	case font.WeightBold:
// 		weight = "bold"
// 	case font.WeightExtraBold:
// 		weight = "extrabold"
// 	case font.WeightBlack:
// 		weight = "black"
// 	}
// 	switch k.Style {
// 	case font.StyleNormal:
// 		style = ""
// 	case font.StyleItalic:
// 		style = "italic"
// 	case font.StyleOblique:
// 		style = "oblique"
// 	}
// 	suffix := mono + weight + style
// 	if len(suffix) == 0 {
// 		suffix = "regular"
// 	}
// 	return strings.ToLower(k.Family) + suffix
// }
