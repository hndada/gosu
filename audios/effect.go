package audios

type EffectPlayer struct {
	MainVolume *float64
}

func (ep EffectPlayer) Play(src []byte, vol float64) {
	player := Context.NewPlayerFromBytes(src)
	player.SetVolume(*ep.MainVolume * vol)
	player.Play()
}
