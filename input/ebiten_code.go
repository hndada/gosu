package input

import "github.com/hajimehoshi/ebiten/v2"

func ebitenKeyToCode(ek int) Code {
	switch ebiten.Key(ek) {
	case ebiten.KeyA:
		return CodeA
	case ebiten.KeyB:
		return CodeB
	case ebiten.KeyC:
		return CodeC
	case ebiten.KeyD:
		return CodeD
	case ebiten.KeyE:
		return CodeE
	case ebiten.KeyF:
		return CodeF
	case ebiten.KeyG:
		return CodeG
	case ebiten.KeyH:
		return CodeH
	case ebiten.KeyI:
		return CodeI
	case ebiten.KeyJ:
		return CodeJ
	case ebiten.KeyK:
		return CodeK
	case ebiten.KeyL:
		return CodeL
	case ebiten.KeyM:
		return CodeM
	case ebiten.KeyN:
		return CodeN
	case ebiten.KeyO:
		return CodeO
	case ebiten.KeyP:
		return CodeP
	case ebiten.KeyQ:
		return CodeQ
	case ebiten.KeyR:
		return CodeR
	case ebiten.KeyS:
		return CodeS
	case ebiten.KeyT:
		return CodeT
	case ebiten.KeyU:
		return CodeU
	case ebiten.KeyV:
		return CodeV
	case ebiten.KeyW:
		return CodeW
	case ebiten.KeyX:
		return CodeX
	case ebiten.KeyY:
		return CodeY
	case ebiten.KeyZ:
		return CodeZ

	case ebiten.Key1:
		return Code1
	case ebiten.Key2:
		return Code2
	case ebiten.Key3:
		return Code3
	case ebiten.Key4:
		return Code4
	case ebiten.Key5:
		return Code5
	case ebiten.Key6:
		return Code6
	case ebiten.Key7:
		return Code7
	case ebiten.Key8:
		return Code8
	case ebiten.Key9:
		return Code9
	case ebiten.Key0:
		return Code0

	case ebiten.KeyEnter:
		return CodeReturnEnter
	case ebiten.KeyEscape:
		return CodeEscape
	case ebiten.KeyBackspace:
		return CodeDeleteBackspace
	case ebiten.KeyTab:
		return CodeTab
	case ebiten.KeySpace:
		return CodeSpacebar
	case ebiten.KeyMinus:
		return CodeHyphenMinus // -
	case ebiten.KeyEqual:
		return CodeEqualSign // =
	case ebiten.KeyBracketLeft:
		return CodeLeftSquareBracket // [
	case ebiten.KeyRightBracket:
		return CodeRightSquareBracket // ]
	case ebiten.KeyBackslash:
		return CodeBackslash // \
	case ebiten.KeySemicolon:
		return CodeSemicolon // ;
	case ebiten.KeyApostrophe:
		return CodeApostrophe // '
	case ebiten.KeyGraveAccent:
		return CodeGraveAccent // `
	case ebiten.KeyComma:
		return CodeComma // ,
	case ebiten.KeyPeriod:
		return CodeFullStop // .
	case ebiten.KeySlash:
		return CodeSlash // /
	case ebiten.KeyCapsLock:
		return CodeCapsLock

	case ebiten.KeyF1:
		return CodeF1
	case ebiten.KeyF2:
		return CodeF2
	case ebiten.KeyF3:
		return CodeF3
	case ebiten.KeyF4:
		return CodeF4
	case ebiten.KeyF5:
		return CodeF5
	case ebiten.KeyF6:
		return CodeF6
	case ebiten.KeyF7:
		return CodeF7
	case ebiten.KeyF8:
		return CodeF8
	case ebiten.KeyF9:
		return CodeF9
	case ebiten.KeyF10:
		return CodeF10
	case ebiten.KeyF11:
		return CodeF11
	case ebiten.KeyF12:
		return CodeF12

	case ebiten.KeyPause:
		return CodePause
	case ebiten.KeyInsert:
		return CodeInsert
	case ebiten.KeyHome:
		return CodeHome
	case ebiten.KeyPageUp:
		return CodePageUp
	case ebiten.KeyDelete:
		return CodeDeleteForward
	case ebiten.KeyEnd:
		return CodeEnd
	case ebiten.KeyPageDown:
		return CodePageDown

	case ebiten.KeyArrowRight:
		return CodeRightArrow
	case ebiten.KeyArrowLeft:
		return CodeLeftArrow
	case ebiten.KeyArrowDown:
		return CodeDownArrow
	case ebiten.KeyArrowUp:
		return CodeUpArrow

	case ebiten.KeyNumLock:
		return CodeKeypadNumLock
	case ebiten.KeyNumpadDivide:
		return CodeKeypadSlash
	case ebiten.KeyNumpadMultiply:
		return CodeKeypadAsterisk
	case ebiten.KeyNumpadSubtract:
		return CodeKeypadHyphenMinus // -
	case ebiten.KeyNumpadAdd:
		return CodeKeypadPlusSign // +
	case ebiten.KeyNumpadEnter:
		return CodeKeypadEnter
	case ebiten.KeyNumpad1:
		return CodeKeypad1
	case ebiten.KeyNumpad2:
		return CodeKeypad2
	case ebiten.KeyNumpad3:
		return CodeKeypad3
	case ebiten.KeyNumpad4:
		return CodeKeypad4
	case ebiten.KeyNumpad5:
		return CodeKeypad5
	case ebiten.KeyNumpad6:
		return CodeKeypad6
	case ebiten.KeyNumpad7:
		return CodeKeypad7
	case ebiten.KeyNumpad8:
		return CodeKeypad8
	case ebiten.KeyNumpad9:
		return CodeKeypad9
	case ebiten.KeyNumpad0:
		return CodeKeypad0
	case ebiten.KeyNumpadDecimal:
		return CodeKeypadFullStop // .
	case ebiten.KeyNumpadEqual:
		return CodeKeypadEqualSign // =

		// the rest of the keys are not supported by Ebiten yet
		/*
			case ebiten.KeyF13:
				return CodeF13
			case ebiten.KeyF14:
				return CodeF14
			case ebiten.KeyF15:
				return CodeF15
			case ebiten.KeyF16:
				return CodeF16
			case ebiten.KeyF17:
				return CodeF17
			case ebiten.KeyF18:
				return CodeF18
			case ebiten.KeyF19:
				return CodeF19
			case ebiten.KeyF20:
				return CodeF20
			case ebiten.KeyF21:
				return CodeF21
			case ebiten.KeyF22:
				return CodeF22
			case ebiten.KeyF23:
				return CodeF23
			case ebiten.KeyF24:
				return CodeF24

			case ebiten.KeyHelp:
				return CodeHelp

			case ebiten.KeyMute:
				return CodeMute
			case ebiten.KeyVolumeUp:
				return CodeVolumeUp
			case ebiten.KeyVolumeDown:
				return CodeVolumeDown

			case ebiten.KeyLeftControl:
				return CodeLeftControl
			case ebiten.KeyLeftShift:
				return CodeLeftShift
			case ebiten.KeyLeftAlt:
				return CodeLeftAlt
			case ebiten.KeyLeftGUI:
				return CodeLeftGUI
			case ebiten.KeyRightControl:
				return CodeRightControl
			case ebiten.KeyRightShift:
				return CodeRightShift
			case ebiten.KeyRightAlt:
				return CodeRightAlt
			case ebiten.KeyRightGUI:
				return CodeRightGUI
		*/
	}
	return CodeUnknown
}
