package gosu

type Key int

// The order is consistent with Ebiten.
const (
	KeyA Key = iota
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	KeyAltLeft
	KeyAltRight
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	KeyArrowUp
	KeyBackquote
	KeyBackslash
	KeyBackspace
	KeyBracketLeft
	KeyBracketRight
	KeyCapsLock
	KeyComma
	KeyContextMenu
	KeyControlLeft
	KeyControlRight
	KeyDelete
	KeyDigit0
	KeyDigit1
	KeyDigit2
	KeyDigit3
	KeyDigit4
	KeyDigit5
	KeyDigit6
	KeyDigit7
	KeyDigit8
	KeyDigit9
	KeyEnd
	KeyEnter
	KeyEqual
	KeyEscape
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyHome
	KeyInsert
	KeyMetaLeft
	KeyMetaRight
	KeyMinus
	KeyNumLock
	KeyNumpad0
	KeyNumpad1
	KeyNumpad2
	KeyNumpad3
	KeyNumpad4
	KeyNumpad5
	KeyNumpad6
	KeyNumpad7
	KeyNumpad8
	KeyNumpad9
	KeyNumpadAdd
	KeyNumpadDecimal
	KeyNumpadDivide
	KeyNumpadEnter
	KeyNumpadEqual
	KeyNumpadMultiply
	KeyNumpadSubtract
	KeyPageDown
	KeyPageUp
	KeyPause
	KeyPeriod
	KeyPrintScreen
	KeyQuote
	KeyScrollLock
	KeySemicolon
	KeyShiftLeft
	KeyShiftRight
	KeySlash
	KeySpace
	KeyTab
	KeyReserved0
	KeyReserved1
	KeyReserved2
	KeyReserved3
)

// See https://docs.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes for reference.
// Alt, Shift, and Control has 3 codes: overall, left and right.
// Order of ArrowKeys are different between Ebiten and Windows VK.
// Order of Page Up/Down are different between Ebiten and Windows VK.
// KeyNumpadEnter returns VK_RETURN, while main Enter returns the same.
// KeyReserved0 ~ KeyReserved3 are skipped.
// Supposed KeyContextMenu stands for Applications key, which is next to Right control
// Supposed KeyMeta Left and Right stands for Left and Right Windows key.
// Supposed KeyNumpadEqual returns VK_OEM_PLUS. I guess this is derived from Apple Keyboard.
func ToVirtualKey(k Key) uint32 {
	switch k {
	case KeyA:
		return 0x41
	case KeyB:
		return 0x42
	case KeyC:
		return 0x43
	case KeyD:
		return 0x44
	case KeyE:
		return 0x45
	case KeyF:
		return 0x46
	case KeyG:
		return 0x47
	case KeyH:
		return 0x48
	case KeyI:
		return 0x49
	case KeyJ:
		return 0x4A
	case KeyK:
		return 0x4B
	case KeyL:
		return 0x4C
	case KeyM:
		return 0x4D
	case KeyN:
		return 0x4E
	case KeyO:
		return 0x4F
	case KeyP:
		return 0x50
	case KeyQ:
		return 0x51
	case KeyR:
		return 0x52
	case KeyS:
		return 0x53
	case KeyT:
		return 0x54
	case KeyU:
		return 0x55
	case KeyV:
		return 0x56
	case KeyW:
		return 0x57
	case KeyX:
		return 0x58
	case KeyY:
		return 0x59
	case KeyZ:
		return 0x5A
	case KeyAltLeft:
		return 0xA4 // VK_LMENU
	case KeyAltRight:
		return 0xA5 // VK_RMENU
	case KeyArrowDown:
		return 0x28
	case KeyArrowLeft:
		return 0x25
	case KeyArrowRight:
		return 0x27
	case KeyArrowUp:
		return 0x26
	case KeyBackquote: // "`"
		return 0xC0 // VK_OEM_3
	case KeyBackslash: // "\"
		return 0xDC // VK_OEM_5
	case KeyBackspace:
		return 0x08
	case KeyBracketLeft:
		return 0xDB // VK_OEM_4
	case KeyBracketRight:
		return 0xDD // VK_OEM_6
	case KeyCapsLock:
		return 0x14
	case KeyComma:
		return 0xBC // VK_OEM_COMMA
	case KeyContextMenu:
		return 0x5D // VK_APPS
	case KeyControlLeft:
		return 0xA2 // VK_LCONTROL
	case KeyControlRight:
		return 0xA3 // VK_RCONTROL
	case KeyDelete:
		return 0x2E
	case KeyDigit0:
		return 0x30
	case KeyDigit1:
		return 0x31
	case KeyDigit2:
		return 0x32
	case KeyDigit3:
		return 0x33
	case KeyDigit4:
		return 0x34
	case KeyDigit5:
		return 0x35
	case KeyDigit6:
		return 0x36
	case KeyDigit7:
		return 0x37
	case KeyDigit8:
		return 0x38
	case KeyDigit9:
		return 0x39
	case KeyEnd:
		return 0x23
	case KeyEnter:
		return 0x0D
	case KeyEqual:
		return 0xBB // VK_OEM_PLUS
	case KeyEscape:
		return 0x1B
	case KeyF1:
		return 0x70
	case KeyF2:
		return 0x71
	case KeyF3:
		return 0x72
	case KeyF4:
		return 0x73
	case KeyF5:
		return 0x74
	case KeyF6:
		return 0x75
	case KeyF7:
		return 0x76
	case KeyF8:
		return 0x77
	case KeyF9:
		return 0x78
	case KeyF10:
		return 0x79
	case KeyF11:
		return 0x7A
	case KeyF12:
		return 0x7B
	case KeyHome:
		return 0x24
	case KeyInsert:
		return 0x2D
	case KeyMetaLeft:
		return 0x5B // VK_LWIN
	case KeyMetaRight:
		return 0x5C // VK_RWIN
	case KeyMinus:
		return 0xBD // VK_OEM_MINUS
	case KeyNumLock:
		return 0x90
	case KeyNumpad0:
		return 0x60
	case KeyNumpad1:
		return 0x61
	case KeyNumpad2:
		return 0x62
	case KeyNumpad3:
		return 0x63
	case KeyNumpad4:
		return 0x64
	case KeyNumpad5:
		return 0x65
	case KeyNumpad6:
		return 0x66
	case KeyNumpad7:
		return 0x67
	case KeyNumpad8:
		return 0x68
	case KeyNumpad9:
		return 0x69
	case KeyNumpadAdd:
		return 0x6B // VK_ADD
	case KeyNumpadDecimal:
		return 0x6E
	case KeyNumpadDivide:
		return 0x6F
	case KeyNumpadEnter:
		return 0x0D // VK_RETURN
	case KeyNumpadEqual:
		return 0xBB // VK_OEM_PLUS
	case KeyNumpadMultiply:
		return 0x6A
	case KeyNumpadSubtract:
		return 0x6D
	case KeyPageDown:
		return 0x22 // VK_NEXT
	case KeyPageUp:
		return 0x21 // VK_PRIOR
	case KeyPause:
		return 0x13
	case KeyPeriod:
		return 0xBE // VK_OEM_PERIOD
	case KeyPrintScreen:
		return 0x2C // VK_SNAPSHOT
	case KeyQuote:
		return 0xDE // VK_OEM_7
	case KeyScrollLock:
		return 0x91
	case KeySemicolon:
		return 0xBA // VK_OEM_1
	case KeyShiftLeft:
		return 0xA0 // VK_LSHIFT
	case KeyShiftRight:
		return 0xA1 // // VK_RSHIFT
	case KeySlash:
		return 0xBF // VK_OEM_2
	case KeySpace:
		return 0x20
	case KeyTab:
		return 0x09
	case KeyReserved0:
	case KeyReserved1:
	case KeyReserved2:
	case KeyReserved3:
	}
	return 0x00 // Unknown
}

// for i, s := range strings.Split(text, "\n") {
// 	fmt.Printf("case %s: return 0x%02X\n", s, i+65)
// }
