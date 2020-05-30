package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"time"
)

const sampleRate   = 44100
// Player represents the current audio state.
type Player struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	seBytes      []byte
	seCh         chan []byte
	volume128    int
	musicType    musicType
}
func NewPlayer(audioContext *audio.Context, musicType musicType) (*Player, error) {
	type audioStream interface {
		audio.ReadSeekCloser
		Length() int64
	}

	const bytesPerSample = 4 // TODO: This should be defined in audio package


	var s audioStream
	var err error

	f, err:= ebitenutil.OpenFile("beatmap/audio.mp3")
	if err!=nil { return nil, err}

	s, err = mp3.Decode(audioContext, f)
	if err != nil {
		return nil, err
	}

	p, err := audio.NewPlayer(audioContext, s)
	if err != nil {
		return nil, err
	}
	player := &Player{
		audioContext: audioContext,
		audioPlayer:  p,
		total:        time.Second * time.Duration(s.Length()) / bytesPerSample / sampleRate,
		volume128:    128,
		seCh:         make(chan []byte),
		musicType:    musicType,
	}
	if player.total == 0 {
		player.total = 1
	}
	player.audioPlayer.Play()
	return player, nil
}

func (p *Player) Close() error {
	return p.audioPlayer.Close()
}

func (p *Player) update() error {
	select {
	case p.seBytes = <-p.seCh:
		close(p.seCh)
		p.seCh = nil
	default:
	}

	if p.audioPlayer.IsPlaying() {
		p.current = p.audioPlayer.Current()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		b := ebiten.IsRunnableOnUnfocused()
		ebiten.SetRunnableOnUnfocused(!b)
	}
	return nil
}

func (p *Player) draw(screen *ebiten.Image) {
}
