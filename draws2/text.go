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

type (
	Font = *text.GoTextFaceSource
	Face = text.Face
)

// src which implements the following methods can be loaded efficiently:
// Read([]byte) (int, error)
// ReadAt([]byte, int64) (int, error)
// Seek(int64, int) (int64, error)
func NewFont(src io.ReadSeeker) (Font, error) {
	font, err := text.NewGoTextFaceSource(src)
	if err != nil {
		return nil, err
	}
	return font, nil
}

// bytes.NewReader(b) implements io.ReadSeeker, as well as io.ReaderAt.
func NewFontFromBytes(b []byte) (Font, error) {
	return NewFont(bytes.NewReader(b))
}

func NewFontFromFile(fsys fs.FS, name string) (Font, error) {
	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewFont(f.(io.ReadSeeker))
}

// Face: A style; aka TypeFace. (e.g.: Arial Bold)
// Font: A file which implements Face. (e.g.: Arial Bold 12pt)
// An implementation of different sizes of the same face are generated on the fly.
type FontOptions struct {
	Family string
	Mono   bool
	Weight font.Weight
	Style  font.Style
}

func DefaultFontOptions() FontOptions {
	return FontOptions{
		Family: "Go",
		Mono:   false,
		Weight: font.WeightNormal,
		Style:  font.StyleNormal,
	}
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

// TrueType vs OpenType: OpenType is an extension of TrueType.
var cachedFonts = make(map[FontOptions]Font)

func init() {
	font, _ := NewFontFromBytes(goregular.TTF)
	cachedFonts[DefaultFontOptions()] = font
}

type FaceOptions struct {
	FontOptions
	opentype.FaceOptions
}

func NewFaceOptions() FaceOptions {
	return FaceOptions{
		FontOptions: DefaultFontOptions(),
		FaceOptions: opentype.FaceOptions{
			Size:    12,
			DPI:     72,
			Hinting: font.HintingFull,
		},
	}
}

var cachedFaces = make(map[FaceOptions]Face)

func NewFace(opts FaceOptions) (Face, error) {
	face, ok := cachedFaces[opts]
	if ok {
		return face, nil
	}

	font, ok := cachedFonts[opts.FontOptions]
	if !ok {
		font = cachedFonts[DefaultFontOptions()]
	}
	face = &text.GoTextFace{
		Source: font,
		Size:   opts.Size,
	}
	cachedFaces[opts] = face
	return face, nil
}

type Text struct {
	FaceOptions
	face        Face
	Text        string
	LineSpacing float64
}

func NewText(txt string) Text {
	return Text{Text: txt, LineSpacing: 1.6}
}

func (t *Text) SetFace(opts FaceOptions) error {
	face, err := NewFace(opts)
	t.FaceOptions = opts
	t.face = face
	return err
}

func (t Text) IsEmpty() bool { return len(t.Text) == 0 }

func (t Text) Size() Vector2 {
	return Vec2(text.Measure(t.Text, t.face, t.LineSpacing))
}

func (t Text) Draw(dst Image, op *Op) {
	if t.IsEmpty() {
		return
	}
	text.Draw(dst.Image, t.Text, t.face, &text.DrawOptions{
		DrawImageOptions: *op,
		LayoutOptions: text.LayoutOptions{
			LineSpacingInPixels: t.LineSpacing,
		},
	})
}
