package input

// convVirtualKeyCode converts a Win32 virtual key code number
// into the standard keycodes used by the key package.
func convVirtualKeyCode(vKey uint32) Code {
	switch vKey {
	case 0x01: // VK_LBUTTON left mouse button
	case 0x02: // VK_RBUTTON right mouse button
	case 0x03: // VK_CANCEL control-break processing
	case 0x04: // VK_MBUTTON middle mouse button
	case 0x05: // VK_XBUTTON1 X1 mouse button
	case 0x06: // VK_XBUTTON2 X2 mouse button
	case 0x08: // VK_BACK
		return CodeDeleteBackspace
	case 0x09: // VK_TAB
		return CodeTab
	case 0x0C: // VK_CLEAR
	case 0x0D: // VK_RETURN
		return CodeReturnEnter
	case 0x10: // VK_SHIFT
		return CodeLeftShift
	case 0x11: // VK_CONTROL
		return CodeLeftControl
	case 0x12: // VK_MENU
		return CodeLeftAlt
	case 0x13: // VK_PAUSE
	case 0x14: // VK_CAPITAL
		return CodeCapsLock
	case 0x15: // VK_KANA, VK_HANGUEL, VK_HANGUL
	case 0x17: // VK_JUNJA
	case 0x18: // VK_FINA, L
	case 0x19: // VK_HANJA, VK_KANJI
	case 0x1B: // VK_ESCAPE
		return CodeEscape
	case 0x1C: // VK_CONVERT
	case 0x1D: // VK_NONCONVERT
	case 0x1E: // VK_ACCEPT
	case 0x1F: // VK_MODECHANGE
	case 0x20: // VK_SPACE
		return CodeSpacebar
	case 0x21: // VK_PRIOR
		return CodePageUp
	case 0x22: // VK_NEXT
		return CodePageDown
	case 0x23: // VK_END
		return CodeEnd
	case 0x24: // VK_HOME
		return CodeHome
	case 0x25: // VK_LEFT
		return CodeLeftArrow
	case 0x26: // VK_UP
		return CodeUpArrow
	case 0x27: // VK_RIGHT
		return CodeRightArrow
	case 0x28: // VK_DOWN
		return CodeDownArrow
	case 0x29: // VK_SELECT
	case 0x2A: // VK_PRINT
	case 0x2B: // VK_EXECUTE
	case 0x2C: // VK_SNAPSHOT
	case 0x2D: // VK_INSERT
	case 0x2E: // VK_DELETE
		return CodeDeleteForward
	case 0x2F: // VK_HELP
		return CodeHelp
	case 0x30:
		return Code0
	case 0x31:
		return Code1
	case 0x32:
		return Code2
	case 0x33:
		return Code3
	case 0x34:
		return Code4
	case 0x35:
		return Code5
	case 0x36:
		return Code6
	case 0x37:
		return Code7
	case 0x38:
		return Code8
	case 0x39:
		return Code9
	case 0x41:
		return CodeA
	case 0x42:
		return CodeB
	case 0x43:
		return CodeC
	case 0x44:
		return CodeD
	case 0x45:
		return CodeE
	case 0x46:
		return CodeF
	case 0x47:
		return CodeG
	case 0x48:
		return CodeH
	case 0x49:
		return CodeI
	case 0x4A:
		return CodeJ
	case 0x4B:
		return CodeK
	case 0x4C:
		return CodeL
	case 0x4D:
		return CodeM
	case 0x4E:
		return CodeN
	case 0x4F:
		return CodeO
	case 0x50:
		return CodeP
	case 0x51:
		return CodeQ
	case 0x52:
		return CodeR
	case 0x53:
		return CodeS
	case 0x54:
		return CodeT
	case 0x55:
		return CodeU
	case 0x56:
		return CodeV
	case 0x57:
		return CodeW
	case 0x58:
		return CodeX
	case 0x59:
		return CodeY
	case 0x5A:
		return CodeZ
	case 0x5B: // VK_LWIN
		return CodeLeftGUI
	case 0x5C: // VK_RWIN
		return CodeRightGUI
	case 0x5D: // VK_APPS
	case 0x5F: // VK_SLEEP
	case 0x60: // VK_NUMPAD0
		return CodeKeypad0
	case 0x61: // VK_NUMPAD1
		return CodeKeypad1
	case 0x62: // VK_NUMPAD2
		return CodeKeypad2
	case 0x63: // VK_NUMPAD3
		return CodeKeypad3
	case 0x64: // VK_NUMPAD4
		return CodeKeypad4
	case 0x65: // VK_NUMPAD5
		return CodeKeypad5
	case 0x66: // VK_NUMPAD6
		return CodeKeypad6
	case 0x67: // VK_NUMPAD7
		return CodeKeypad7
	case 0x68: // VK_NUMPAD8
		return CodeKeypad8
	case 0x69: // VK_NUMPAD9
		return CodeKeypad9
	case 0x6A: // VK_MULTIPLY
		return CodeKeypadAsterisk
	case 0x6B: // VK_ADD
		return CodeKeypadPlusSign
	case 0x6C: // VK_SEPARATOR
	case 0x6D: // VK_SUBTRACT
		return CodeKeypadHyphenMinus
	case 0x6E: // VK_DECIMAL
		return CodeFullStop
	case 0x6F: // VK_DIVIDE
		return CodeKeypadSlash
	case 0x70: // VK_F1
		return CodeF1
	case 0x71: // VK_F2
		return CodeF2
	case 0x72: // VK_F3
		return CodeF3
	case 0x73: // VK_F4
		return CodeF4
	case 0x74: // VK_F5
		return CodeF5
	case 0x75: // VK_F6
		return CodeF6
	case 0x76: // VK_F7
		return CodeF7
	case 0x77: // VK_F8
		return CodeF8
	case 0x78: // VK_F9
		return CodeF9
	case 0x79: // VK_F10
		return CodeF10
	case 0x7A: // VK_F11
		return CodeF11
	case 0x7B: // VK_F12
		return CodeF12
	case 0x7C: // VK_F13
		return CodeF13
	case 0x7D: // VK_F14
		return CodeF14
	case 0x7E: // VK_F15
		return CodeF15
	case 0x7F: // VK_F16
		return CodeF16
	case 0x80: // VK_F17
		return CodeF17
	case 0x81: // VK_F18
		return CodeF18
	case 0x82: // VK_F19
		return CodeF19
	case 0x83: // VK_F20
		return CodeF20
	case 0x84: // VK_F21
		return CodeF21
	case 0x85: // VK_F22
		return CodeF22
	case 0x86: // VK_F23
		return CodeF23
	case 0x87: // VK_F24
		return CodeF24
	case 0x90: // VK_NUMLOCK
		return CodeKeypadNumLock
	case 0x91: // VK_SCROLL
	case 0xA0: // VK_LSHIFT
		return CodeLeftShift
	case 0xA1: // VK_RSHIFT
		return CodeRightShift
	case 0xA2: // VK_LCONTROL
		return CodeLeftControl
	case 0xA3: // VK_RCONTROL
		return CodeRightControl
	case 0xA4: // VK_LMENU
	case 0xA5: // VK_RMENU
	case 0xA6: // VK_BROWSER_BACK
	case 0xA7: // VK_BROWSER_FORWARD
	case 0xA8: // VK_BROWSER_REFRESH
	case 0xA9: // VK_BROWSER_STOP
	case 0xAA: // VK_BROWSER_SEARCH
	case 0xAB: // VK_BROWSER_FAVORITES
	case 0xAC: // VK_BROWSER_HOME
	case 0xAD: // VK_VOLUME_MUTE
		return CodeMute
	case 0xAE: // VK_VOLUME_DOWN
		return CodeVolumeDown
	case 0xAF: // VK_VOLUME_UP
		return CodeVolumeUp
	case 0xB0: // VK_MEDIA_NEXT_TRACK
	case 0xB1: // VK_MEDIA_PREV_TRACK
	case 0xB2: // VK_MEDIA_STOP
	case 0xB3: // VK_MEDIA_PLAY_PAUSE
	case 0xB4: // VK_LAUNCH_MAIL
	case 0xB5: // VK_LAUNCH_MEDIA_SELECT
	case 0xB6: // VK_LAUNCH_APP1
	case 0xB7: // VK_LAUNCH_APP2
	case 0xBA: // VK_OEM_1 ';:'
		return CodeSemicolon
	case 0xBB: // VK_OEM_PLUS '+'
		return CodeEqualSign
	case 0xBC: // VK_OEM_COMMA ','
		return CodeComma
	case 0xBD: // VK_OEM_MINUS '-'
		return CodeHyphenMinus
	case 0xBE: // VK_OEM_PERIOD '.'
		return CodeFullStop
	case 0xBF: // VK_OEM_2 '/?'
		return CodeSlash
	case 0xC0: // VK_OEM_3 '`~'
		return CodeGraveAccent
	case 0xDB: // VK_OEM_4 '[{'
		return CodeLeftSquareBracket
	case 0xDC: // VK_OEM_5 '\|'
		return CodeBackslash
	case 0xDD: // VK_OEM_6 ']}'
		return CodeRightSquareBracket
	case 0xDE: // VK_OEM_7 'single-quote/double-quote'
		return CodeApostrophe
	case 0xDF: // VK_OEM_8
		return CodeUnknown
	case 0xE2: // VK_OEM_102
	case 0xE5: // VK_PROCESSKEY
	case 0xE7: // VK_PACKET
	case 0xF6: // VK_ATTN
	case 0xF7: // VK_CRSEL
	case 0xF8: // VK_EXSEL
	case 0xF9: // VK_EREOF
	case 0xFA: // VK_PLAY
	case 0xFB: // VK_ZOOM
	case 0xFC: // VK_NONAME
	case 0xFD: // VK_PA1
	case 0xFE: // VK_OEM_CLEAR
	}
	return CodeUnknown
}
