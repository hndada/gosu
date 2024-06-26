package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// functions
var IsKeyPressed = ebiten.IsKeyPressed
var IsKeyJustPressed = inpututil.IsKeyJustPressed

type Key = ebiten.Key

const KeyNone Key = -1

// The order is consistent with Ebiten.
const (
	KeyA              Key = ebiten.KeyA
	KeyB              Key = ebiten.KeyB
	KeyC              Key = ebiten.KeyC
	KeyD              Key = ebiten.KeyD
	KeyE              Key = ebiten.KeyE
	KeyF              Key = ebiten.KeyF
	KeyG              Key = ebiten.KeyG
	KeyH              Key = ebiten.KeyH
	KeyI              Key = ebiten.KeyI
	KeyJ              Key = ebiten.KeyJ
	KeyK              Key = ebiten.KeyK
	KeyL              Key = ebiten.KeyL
	KeyM              Key = ebiten.KeyM
	KeyN              Key = ebiten.KeyN
	KeyO              Key = ebiten.KeyO
	KeyP              Key = ebiten.KeyP
	KeyQ              Key = ebiten.KeyQ
	KeyR              Key = ebiten.KeyR
	KeyS              Key = ebiten.KeyS
	KeyT              Key = ebiten.KeyT
	KeyU              Key = ebiten.KeyU
	KeyV              Key = ebiten.KeyV
	KeyW              Key = ebiten.KeyW
	KeyX              Key = ebiten.KeyX
	KeyY              Key = ebiten.KeyY
	KeyZ              Key = ebiten.KeyZ
	KeyAltLeft        Key = ebiten.KeyAltLeft
	KeyAltRight       Key = ebiten.KeyAltRight
	KeyArrowDown      Key = ebiten.KeyArrowDown
	KeyArrowLeft      Key = ebiten.KeyArrowLeft
	KeyArrowRight     Key = ebiten.KeyArrowRight
	KeyArrowUp        Key = ebiten.KeyArrowUp
	KeyBackquote      Key = ebiten.KeyBackquote
	KeyBackslash      Key = ebiten.KeyBackslash
	KeyBackspace      Key = ebiten.KeyBackspace
	KeyBracketLeft    Key = ebiten.KeyBracketLeft
	KeyBracketRight   Key = ebiten.KeyBracketRight
	KeyCapsLock       Key = ebiten.KeyCapsLock
	KeyComma          Key = ebiten.KeyComma
	KeyContextMenu    Key = ebiten.KeyContextMenu
	KeyControlLeft    Key = ebiten.KeyControlLeft
	KeyControlRight   Key = ebiten.KeyControlRight
	KeyDelete         Key = ebiten.KeyDelete
	KeyDigit0         Key = ebiten.KeyDigit0
	KeyDigit1         Key = ebiten.KeyDigit1
	KeyDigit2         Key = ebiten.KeyDigit2
	KeyDigit3         Key = ebiten.KeyDigit3
	KeyDigit4         Key = ebiten.KeyDigit4
	KeyDigit5         Key = ebiten.KeyDigit5
	KeyDigit6         Key = ebiten.KeyDigit6
	KeyDigit7         Key = ebiten.KeyDigit7
	KeyDigit8         Key = ebiten.KeyDigit8
	KeyDigit9         Key = ebiten.KeyDigit9
	KeyEnd            Key = ebiten.KeyEnd
	KeyEnter          Key = ebiten.KeyEnter
	KeyEqual          Key = ebiten.KeyEqual
	KeyEscape         Key = ebiten.KeyEscape
	KeyF1             Key = ebiten.KeyF1
	KeyF2             Key = ebiten.KeyF2
	KeyF3             Key = ebiten.KeyF3
	KeyF4             Key = ebiten.KeyF4
	KeyF5             Key = ebiten.KeyF5
	KeyF6             Key = ebiten.KeyF6
	KeyF7             Key = ebiten.KeyF7
	KeyF8             Key = ebiten.KeyF8
	KeyF9             Key = ebiten.KeyF9
	KeyF10            Key = ebiten.KeyF10
	KeyF11            Key = ebiten.KeyF11
	KeyF12            Key = ebiten.KeyF12
	KeyHome           Key = ebiten.KeyHome
	KeyInsert         Key = ebiten.KeyInsert
	KeyMetaLeft       Key = ebiten.KeyMetaLeft
	KeyMetaRight      Key = ebiten.KeyMetaRight
	KeyMinus          Key = ebiten.KeyMinus
	KeyNumLock        Key = ebiten.KeyNumLock
	KeyNumpad0        Key = ebiten.KeyNumpad0
	KeyNumpad1        Key = ebiten.KeyNumpad1
	KeyNumpad2        Key = ebiten.KeyNumpad2
	KeyNumpad3        Key = ebiten.KeyNumpad3
	KeyNumpad4        Key = ebiten.KeyNumpad4
	KeyNumpad5        Key = ebiten.KeyNumpad5
	KeyNumpad6        Key = ebiten.KeyNumpad6
	KeyNumpad7        Key = ebiten.KeyNumpad7
	KeyNumpad8        Key = ebiten.KeyNumpad8
	KeyNumpad9        Key = ebiten.KeyNumpad9
	KeyNumpadAdd      Key = ebiten.KeyNumpadAdd
	KeyNumpadDecimal  Key = ebiten.KeyNumpadDecimal
	KeyNumpadDivide   Key = ebiten.KeyNumpadDivide
	KeyNumpadEnter    Key = ebiten.KeyNumpadEnter
	KeyNumpadEqual    Key = ebiten.KeyNumpadEqual
	KeyNumpadMultiply Key = ebiten.KeyNumpadMultiply
	KeyNumpadSubtract Key = ebiten.KeyNumpadSubtract
	KeyPageDown       Key = ebiten.KeyPageDown
	KeyPageUp         Key = ebiten.KeyPageUp
	KeyPause          Key = ebiten.KeyPause
	KeyPeriod         Key = ebiten.KeyPeriod
	KeyPrintScreen    Key = ebiten.KeyPrintScreen
	KeyQuote          Key = ebiten.KeyQuote
	KeyScrollLock     Key = ebiten.KeyScrollLock
	KeySemicolon      Key = ebiten.KeySemicolon
	KeyShiftLeft      Key = ebiten.KeyShiftLeft
	KeyShiftRight     Key = ebiten.KeyShiftRight
	KeySlash          Key = ebiten.KeySlash
	KeySpace          Key = ebiten.KeySpace
	KeyTab            Key = ebiten.KeyTab
	// KeyReserved0      Key = ebiten.KeyReserved0
	// KeyReserved1      Key = ebiten.KeyReserved1
	// KeyReserved2      Key = ebiten.KeyReserved2
	// KeyReserved3      Key = ebiten.KeyReserved3
)
const KeyFinal Key = KeyTab + 1 // For iterating keys

func NamesToKeys(names []string) []Key {
	keys := make([]Key, len(names))
	for i, name := range names {
		keys[i] = NameToKey(name)
	}
	return keys
}
func KeysToNames(keys []Key) []string {
	names := make([]string, len(keys))
	for i, key := range keys {
		names[i] = KeyToName(key)
	}
	return names
}

// https://go.dev/play/p/9Lv0u4sqwKq
func NameToKey(name string) Key {
	switch name {
	case "A":
		return KeyA
	case "B":
		return KeyB
	case "C":
		return KeyC
	case "D":
		return KeyD
	case "E":
		return KeyE
	case "F":
		return KeyF
	case "G":
		return KeyG
	case "H":
		return KeyH
	case "I":
		return KeyI
	case "J":
		return KeyJ
	case "K":
		return KeyK
	case "L":
		return KeyL
	case "M":
		return KeyM
	case "N":
		return KeyN
	case "O":
		return KeyO
	case "P":
		return KeyP
	case "Q":
		return KeyQ
	case "R":
		return KeyR
	case "S":
		return KeyS
	case "T":
		return KeyT
	case "U":
		return KeyU
	case "V":
		return KeyV
	case "W":
		return KeyW
	case "X":
		return KeyX
	case "Y":
		return KeyY
	case "Z":
		return KeyZ
	case "AltLeft":
		return KeyAltLeft
	case "AltRight":
		return KeyAltRight
	case "ArrowDown":
		return KeyArrowDown
	case "ArrowLeft":
		return KeyArrowLeft
	case "ArrowRight":
		return KeyArrowRight
	case "ArrowUp":
		return KeyArrowUp
	case "Backquote":
		return KeyBackquote
	case "Backslash":
		return KeyBackslash
	case "Backspace":
		return KeyBackspace
	case "BracketLeft":
		return KeyBracketLeft
	case "BracketRight":
		return KeyBracketRight
	case "CapsLock":
		return KeyCapsLock
	case "Comma":
		return KeyComma
	case "ContextMenu":
		return KeyContextMenu
	case "ControlLeft":
		return KeyControlLeft
	case "ControlRight":
		return KeyControlRight
	case "Delete":
		return KeyDelete
	case "Digit0":
		return KeyDigit0
	case "Digit1":
		return KeyDigit1
	case "Digit2":
		return KeyDigit2
	case "Digit3":
		return KeyDigit3
	case "Digit4":
		return KeyDigit4
	case "Digit5":
		return KeyDigit5
	case "Digit6":
		return KeyDigit6
	case "Digit7":
		return KeyDigit7
	case "Digit8":
		return KeyDigit8
	case "Digit9":
		return KeyDigit9
	case "End":
		return KeyEnd
	case "Enter":
		return KeyEnter
	case "Equal":
		return KeyEqual
	case "Escape":
		return KeyEscape
	case "F1":
		return KeyF1
	case "F2":
		return KeyF2
	case "F3":
		return KeyF3
	case "F4":
		return KeyF4
	case "F5":
		return KeyF5
	case "F6":
		return KeyF6
	case "F7":
		return KeyF7
	case "F8":
		return KeyF8
	case "F9":
		return KeyF9
	case "F10":
		return KeyF10
	case "F11":
		return KeyF11
	case "F12":
		return KeyF12
	case "Home":
		return KeyHome
	case "Insert":
		return KeyInsert
	case "MetaLeft":
		return KeyMetaLeft
	case "MetaRight":
		return KeyMetaRight
	case "Minus":
		return KeyMinus
	case "NumLock":
		return KeyNumLock
	case "Numpad0":
		return KeyNumpad0
	case "Numpad1":
		return KeyNumpad1
	case "Numpad2":
		return KeyNumpad2
	case "Numpad3":
		return KeyNumpad3
	case "Numpad4":
		return KeyNumpad4
	case "Numpad5":
		return KeyNumpad5
	case "Numpad6":
		return KeyNumpad6
	case "Numpad7":
		return KeyNumpad7
	case "Numpad8":
		return KeyNumpad8
	case "Numpad9":
		return KeyNumpad9
	case "NumpadAdd":
		return KeyNumpadAdd
	case "NumpadDecimal":
		return KeyNumpadDecimal
	case "NumpadDivide":
		return KeyNumpadDivide
	case "NumpadEnter":
		return KeyNumpadEnter
	case "NumpadEqual":
		return KeyNumpadEqual
	case "NumpadMultiply":
		return KeyNumpadMultiply
	case "NumpadSubtract":
		return KeyNumpadSubtract
	case "PageDown":
		return KeyPageDown
	case "PageUp":
		return KeyPageUp
	case "Pause":
		return KeyPause
	case "Period":
		return KeyPeriod
	case "PrintScreen":
		return KeyPrintScreen
	case "Quote":
		return KeyQuote
	case "ScrollLock":
		return KeyScrollLock
	case "Semicolon":
		return KeySemicolon
	case "ShiftLeft":
		return KeyShiftLeft
	case "ShiftRight":
		return KeyShiftRight
	case "Slash":
		return KeySlash
	case "Space":
		return KeySpace
	case "Tab":
		return KeyTab
		// case "Reserved0":
		// 	return KeyReserved0
		// case "Reserved1":
		// 	return KeyReserved1
		// case "Reserved2":
		// 	return KeyReserved2
		// case "Reserved3":
		// 	return KeyReserved3
	}
	return KeyNone
}

// https://go.dev/play/p/H8IQTm5BEBp
func KeyToName(k Key) string {
	switch k {
	case KeyA:
		return "A"
	case KeyB:
		return "B"
	case KeyC:
		return "C"
	case KeyD:
		return "D"
	case KeyE:
		return "E"
	case KeyF:
		return "F"
	case KeyG:
		return "G"
	case KeyH:
		return "H"
	case KeyI:
		return "I"
	case KeyJ:
		return "J"
	case KeyK:
		return "K"
	case KeyL:
		return "L"
	case KeyM:
		return "M"
	case KeyN:
		return "N"
	case KeyO:
		return "O"
	case KeyP:
		return "P"
	case KeyQ:
		return "Q"
	case KeyR:
		return "R"
	case KeyS:
		return "S"
	case KeyT:
		return "T"
	case KeyU:
		return "U"
	case KeyV:
		return "V"
	case KeyW:
		return "W"
	case KeyX:
		return "X"
	case KeyY:
		return "Y"
	case KeyZ:
		return "Z"
	case KeyAltLeft:
		return "AltLeft"
	case KeyAltRight:
		return "AltRight"
	case KeyArrowDown:
		return "ArrowDown"
	case KeyArrowLeft:
		return "ArrowLeft"
	case KeyArrowRight:
		return "ArrowRight"
	case KeyArrowUp:
		return "ArrowUp"
	case KeyBackquote:
		return "Backquote"
	case KeyBackslash:
		return "Backslash"
	case KeyBackspace:
		return "Backspace"
	case KeyBracketLeft:
		return "BracketLeft"
	case KeyBracketRight:
		return "BracketRight"
	case KeyCapsLock:
		return "CapsLock"
	case KeyComma:
		return "Comma"
	case KeyContextMenu:
		return "ContextMenu"
	case KeyControlLeft:
		return "ControlLeft"
	case KeyControlRight:
		return "ControlRight"
	case KeyDelete:
		return "Delete"
	case KeyDigit0:
		return "Digit0"
	case KeyDigit1:
		return "Digit1"
	case KeyDigit2:
		return "Digit2"
	case KeyDigit3:
		return "Digit3"
	case KeyDigit4:
		return "Digit4"
	case KeyDigit5:
		return "Digit5"
	case KeyDigit6:
		return "Digit6"
	case KeyDigit7:
		return "Digit7"
	case KeyDigit8:
		return "Digit8"
	case KeyDigit9:
		return "Digit9"
	case KeyEnd:
		return "End"
	case KeyEnter:
		return "Enter"
	case KeyEqual:
		return "Equal"
	case KeyEscape:
		return "Escape"
	case KeyF1:
		return "F1"
	case KeyF2:
		return "F2"
	case KeyF3:
		return "F3"
	case KeyF4:
		return "F4"
	case KeyF5:
		return "F5"
	case KeyF6:
		return "F6"
	case KeyF7:
		return "F7"
	case KeyF8:
		return "F8"
	case KeyF9:
		return "F9"
	case KeyF10:
		return "F10"
	case KeyF11:
		return "F11"
	case KeyF12:
		return "F12"
	case KeyHome:
		return "Home"
	case KeyInsert:
		return "Insert"
	case KeyMetaLeft:
		return "MetaLeft"
	case KeyMetaRight:
		return "MetaRight"
	case KeyMinus:
		return "Minus"
	case KeyNumLock:
		return "NumLock"
	case KeyNumpad0:
		return "Numpad0"
	case KeyNumpad1:
		return "Numpad1"
	case KeyNumpad2:
		return "Numpad2"
	case KeyNumpad3:
		return "Numpad3"
	case KeyNumpad4:
		return "Numpad4"
	case KeyNumpad5:
		return "Numpad5"
	case KeyNumpad6:
		return "Numpad6"
	case KeyNumpad7:
		return "Numpad7"
	case KeyNumpad8:
		return "Numpad8"
	case KeyNumpad9:
		return "Numpad9"
	case KeyNumpadAdd:
		return "NumpadAdd"
	case KeyNumpadDecimal:
		return "NumpadDecimal"
	case KeyNumpadDivide:
		return "NumpadDivide"
	case KeyNumpadEnter:
		return "NumpadEnter"
	case KeyNumpadEqual:
		return "NumpadEqual"
	case KeyNumpadMultiply:
		return "NumpadMultiply"
	case KeyNumpadSubtract:
		return "NumpadSubtract"
	case KeyPageDown:
		return "PageDown"
	case KeyPageUp:
		return "PageUp"
	case KeyPause:
		return "Pause"
	case KeyPeriod:
		return "Period"
	case KeyPrintScreen:
		return "PrintScreen"
	case KeyQuote:
		return "Quote"
	case KeyScrollLock:
		return "ScrollLock"
	case KeySemicolon:
		return "Semicolon"
	case KeyShiftLeft:
		return "ShiftLeft"
	case KeyShiftRight:
		return "ShiftRight"
	case KeySlash:
		return "Slash"
	case KeySpace:
		return "Space"
	case KeyTab:
		return "Tab"
		// case KeyReserved0:
		// 	return "Reserved0"
		// case KeyReserved1:
		// 	return "Reserved1"
		// case KeyReserved2:
		// 	return "Reserved2"
		// case KeyReserved3:
		// 	return "Reserved3"
	}
	return "(None)"
}

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
		// case KeyReserved0:
		// case KeyReserved1:
		// case KeyReserved2:
		// case KeyReserved3:
	}
	return 0x00 // Unknown
}

// for i, s := range strings.Split(text, "\n") {
// 	fmt.Printf("case %s: return 0x%02X\n", s, i+65)
// }
