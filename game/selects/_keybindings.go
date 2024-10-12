// choose key bindings from finite selections.
// Todo: KeySettings -> KeyBinding?
// if inpututil.IsKeyJustPressed(input.KeyF5) {
// 	if s.Focus != FocusKeySettings {
// 		s.lastFocus = s.Focus
// 	}
// 	s.keySettings = make([]string, 0)
// 	s.Focus = FocusKeySettings
// 	scene.UserSkin.Swipe.Play(*s.volumeSound)
// }

func setKeySettings() {
	for k := input.Key(0); k < input.KeyReserved0; k++ {
		if input.IsKeyJustPressed(k) {
			name := input.KeyToName(k)
			if name[0] == 'F' && name != "F" {
				continue
			}
			s.keySettings = append(s.keySettings, name)
		}
	}
	switch s.mode {
	case 0:
		if len(s.keySettings) >= 4 {
			s.keySettings = s.keySettings[:4]
			s.keySettings = mode.NormalizeKeys(s.keySettings)
			s.Focus = s.lastFocus
		}
	}
}
