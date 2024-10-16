package draws

import (
	"bytes"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

// Font: Face source; a file which implements Face.
// Face: A style; aka TypeFace. (e.g.: Arial Bold)
var cachedFaceSources = make(map[string]*text.GoTextFaceSource)
var cachedFaces = make(map[FaceOptions]text.Face)

// src which implements the following methods can be loaded efficiently:
// Read([]byte) (int, error)
// ReadAt([]byte, int64) (int, error)
// Seek(int64, int) (int64, error)
// bytes.NewReader(b) implements io.ReadSeeker, as well as io.ReaderAt.
func LoadFaceSource(name string, src io.ReadSeeker) error {
	font, err := text.NewGoTextFaceSource(src)
	if err != nil {
		return err
	}
	cachedFaceSources[name] = font
	return nil
}

func LoadFaceSourceFromFile(fsys fs.FS, name string) error {
	f, err := fsys.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	return LoadFaceSource(filepath.Base(name), f.(io.ReadSeeker))
}

// TrueType vs OpenType: OpenType is an extension of TrueType.
type FaceOptions struct {
	Font string // Face source name
	opentype.FaceOptions
}

func NewFaceOptions() FaceOptions {
	return FaceOptions{
		Font: "goregular",
		FaceOptions: opentype.FaceOptions{
			Size:    24,
			DPI:     72,
			Hinting: font.HintingVertical,
		},
	}
}

func LoadFace(opts FaceOptions) {
	src, ok := cachedFaceSources[opts.Font]
	if !ok {
		src = cachedFaceSources["goregular"]
	}

	cachedFaces[opts] = &text.GoTextFace{
		Source: src,
		Size:   opts.Size,
	}
}

func init() {
	_ = LoadFaceSource("goregular", bytes.NewReader(goregular.TTF))
	LoadFace(NewFaceOptions())
}

// Text does not have to have a separate source struct
// unlike Image/Sprite, as it is lightweight enough.
type Text struct {
	Text        string
	LineSpacing float64
	FaceOptions
	face text.Face
	Box
}

func NewText(txt string) Text {
	t := Text{
		Text:        txt,
		LineSpacing: 1.6,
	}
	t.SetFace(NewFaceOptions())
	t.Box = NewBox(&t)
	return t
}

func (t *Text) SetFace(opts FaceOptions) {
	LoadFace(opts)
	t.FaceOptions = opts
	t.face = cachedFaces[opts]
}

// func (t Text) IsEmpty() bool { return len(t.Text) == 0 }

func (t Text) Size() XY {
	return NewXY(text.Measure(t.Text, t.face, t.LineSpacing))
}

func (t Text) Draw(dst Image) {
	// if t.IsEmpty() {
	// 	return
	// }
	text.Draw(dst.Image, t.Text, t.face, &text.DrawOptions{
		DrawImageOptions: *t.op(),
		LayoutOptions: text.LayoutOptions{
			LineSpacing: t.LineSpacing,
		},
	})
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
