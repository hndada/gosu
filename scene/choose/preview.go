package choose

import (
	"bytes"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/mode"
)

type PreviewPlayer struct {
	*audio.Player
	Closer func() error
	Tick   int
	Volume float64
	volume *float64
}

const wait = -500

func NewPreviewPlayer(rc io.ReadCloser) (p PreviewPlayer, err error) {
	b, err := io.ReadAll(rc)
	if err != nil {
		return
	}
	r := bytes.NewReader(b)
	streamer, err := mp3.DecodeWithSampleRate(audios.SampleRate, r)
	if err != nil {
		// fmt.Println("decode:", err)
		return
	}
	_p, err := audios.Context.NewPlayer(streamer)
	if err != nil {
		return
	}
	return PreviewPlayer{
		Player: _p,
		Closer: rc.Close,
		Tick:   wait,
		Volume: mode.S.VolumeMusic,
		volume: &mode.S.VolumeMusic,
	}, nil
}
func (p *PreviewPlayer) Update() {
	if !p.IsValid() {
		return
	}
	if vol := *p.volume; p.Volume != vol {
		p.SetVolume(p.Volume)
		p.Volume = vol
	}
	p.Tick++
	if p.Tick == 0 {
		p.Play()
	}
	if p.Tick > 0 && p.Tick <= 1000 {
		age := float64(p.Tick) / 1000
		p.SetVolume(p.Volume * age)
	}
	if p.Tick > 9000 && p.Tick <= 10000 {
		age := float64(p.Tick-9000) / 1000
		p.SetVolume(p.Volume * (1 - age))
	}
	if p.Tick >= 10000 {
		p.Tick = wait
		p.Rewind()
		p.Pause()
	}
}
func (p PreviewPlayer) IsValid() bool { return p.Player != nil }
func (p PreviewPlayer) Close() {
	p.Player.Close()
	p.Closer()
}
