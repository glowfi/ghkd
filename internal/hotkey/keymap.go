package hotkey

import (
	"strings"
)

var KeyNameToCode = map[string]uint16{
	// Modifiers
	"ctrl":       KEY_LEFTCTRL,
	"control":    KEY_LEFTCTRL,
	"leftctrl":   KEY_LEFTCTRL,
	"rightctrl":  KEY_RIGHTCTRL,
	"alt":        KEY_LEFTALT,
	"leftalt":    KEY_LEFTALT,
	"rightalt":   KEY_RIGHTALT,
	"shift":      KEY_LEFTSHIFT,
	"leftshift":  KEY_LEFTSHIFT,
	"rightshift": KEY_RIGHTSHIFT,
	"super":      KEY_LEFTMETA,
	"meta":       KEY_LEFTMETA,
	"win":        KEY_LEFTMETA,
	"leftmeta":   KEY_LEFTMETA,
	"rightmeta":  KEY_RIGHTMETA,

	// Letters
	"a": KEY_A,
	"b": KEY_B,
	"c": KEY_C,
	"d": KEY_D,
	"e": KEY_E,
	"f": KEY_F,
	"g": KEY_G,
	"h": KEY_H,
	"i": KEY_I,
	"j": KEY_J,
	"k": KEY_K,
	"l": KEY_L,
	"m": KEY_M,
	"n": KEY_N,
	"o": KEY_O,
	"p": KEY_P,
	"q": KEY_Q,
	"r": KEY_R,
	"s": KEY_S,
	"t": KEY_T,
	"u": KEY_U,
	"v": KEY_V,
	"w": KEY_W,
	"x": KEY_X,
	"y": KEY_Y,
	"z": KEY_Z,

	// Numbers
	"0": KEY_0,
	"1": KEY_1,
	"2": KEY_2,
	"3": KEY_3,
	"4": KEY_4,
	"5": KEY_5,
	"6": KEY_6,
	"7": KEY_7,
	"8": KEY_8,
	"9": KEY_9,

	// Function keys
	"f1":  KEY_F1,
	"f2":  KEY_F2,
	"f3":  KEY_F3,
	"f4":  KEY_F4,
	"f5":  KEY_F5,
	"f6":  KEY_F6,
	"f7":  KEY_F7,
	"f8":  KEY_F8,
	"f9":  KEY_F9,
	"f10": KEY_F10,
	"f11": KEY_F11,
	"f12": KEY_F12,

	// Special keys
	"esc":         KEY_ESC,
	"escape":      KEY_ESC,
	"tab":         KEY_TAB,
	"space":       KEY_SPACE,
	"enter":       KEY_ENTER,
	"return":      KEY_ENTER,
	"backspace":   KEY_BACKSPACE,
	"delete":      KEY_DELETE,
	"insert":      KEY_INSERT,
	"home":        KEY_HOME,
	"end":         KEY_END,
	"pageup":      KEY_PAGEUP,
	"pagedown":    KEY_PAGEDOWN,
	"up":          KEY_UP,
	"down":        KEY_DOWN,
	"left":        KEY_LEFT,
	"right":       KEY_RIGHT,
	"capslock":    KEY_CAPSLOCK,
	"print":       KEY_PRINT,
	"printscreen": KEY_PRINT,
	"scrolllock":  KEY_SCROLLLOCK,
	"pause":       KEY_PAUSE,

	// Punctuation
	"minus":      KEY_MINUS,
	"equal":      KEY_EQUAL,
	"leftbrace":  KEY_LEFTBRACE,
	"rightbrace": KEY_RIGHTBRACE,
	"semicolon":  KEY_SEMICOLON,
	"apostrophe": KEY_APOSTROPHE,
	"grave":      KEY_GRAVE,
	"backslash":  KEY_BACKSLASH,
	"comma":      KEY_COMMA,
	"dot":        KEY_DOT,
	"slash":      KEY_SLASH,

	// Media keys (XF86 names)
	"xf86audiomute":         KEY_MUTE,
	"xf86audiolowervolume":  KEY_VOLUMEDOWN,
	"xf86audioraisevolume":  KEY_VOLUMEUP,
	"xf86audioplay":         KEY_PLAYPAUSE,
	"xf86audiopause":        KEY_PLAYPAUSE,
	"xf86audionext":         KEY_NEXTSONG,
	"xf86audioprev":         KEY_PREVIOUSSONG,
	"xf86audiostop":         KEY_STOPCD,
	"xf86monbrightnessup":   KEY_BRIGHTNESSUP,
	"xf86monbrightnessdown": KEY_BRIGHTNESSDOWN,

	// Short media key names (convenience)
	"mute":           KEY_MUTE,
	"volumedown":     KEY_VOLUMEDOWN,
	"volumeup":       KEY_VOLUMEUP,
	"playpause":      KEY_PLAYPAUSE,
	"nextsong":       KEY_NEXTSONG,
	"previoussong":   KEY_PREVIOUSSONG,
	"brightnessup":   KEY_BRIGHTNESSUP,
	"brightnessdown": KEY_BRIGHTNESSDOWN,
}

var KeyCodeToName = map[uint16]string{
	// Modifiers
	KEY_LEFTCTRL:   "ctrl",
	KEY_RIGHTCTRL:  "rightctrl",
	KEY_LEFTALT:    "alt",
	KEY_RIGHTALT:   "rightalt",
	KEY_LEFTSHIFT:  "shift",
	KEY_RIGHTSHIFT: "rightshift",
	KEY_LEFTMETA:   "super",
	KEY_RIGHTMETA:  "rightmeta",

	// Letters
	KEY_A: "a",
	KEY_B: "b",
	KEY_C: "c",
	KEY_D: "d",
	KEY_E: "e",
	KEY_F: "f",
	KEY_G: "g",
	KEY_H: "h",
	KEY_I: "i",
	KEY_J: "j",
	KEY_K: "k",
	KEY_L: "l",
	KEY_M: "m",
	KEY_N: "n",
	KEY_O: "o",
	KEY_P: "p",
	KEY_Q: "q",
	KEY_R: "r",
	KEY_S: "s",
	KEY_T: "t",
	KEY_U: "u",
	KEY_V: "v",
	KEY_W: "w",
	KEY_X: "x",
	KEY_Y: "y",
	KEY_Z: "z",

	// Numbers
	KEY_0: "0",
	KEY_1: "1",
	KEY_2: "2",
	KEY_3: "3",
	KEY_4: "4",
	KEY_5: "5",
	KEY_6: "6",
	KEY_7: "7",
	KEY_8: "8",
	KEY_9: "9",

	// Function keys
	KEY_F1:  "f1",
	KEY_F2:  "f2",
	KEY_F3:  "f3",
	KEY_F4:  "f4",
	KEY_F5:  "f5",
	KEY_F6:  "f6",
	KEY_F7:  "f7",
	KEY_F8:  "f8",
	KEY_F9:  "f9",
	KEY_F10: "f10",
	KEY_F11: "f11",
	KEY_F12: "f12",

	// Special keys
	KEY_ESC:        "esc",
	KEY_TAB:        "tab",
	KEY_SPACE:      "space",
	KEY_ENTER:      "enter",
	KEY_BACKSPACE:  "backspace",
	KEY_DELETE:     "delete",
	KEY_INSERT:     "insert",
	KEY_HOME:       "home",
	KEY_END:        "end",
	KEY_PAGEUP:     "pageup",
	KEY_PAGEDOWN:   "pagedown",
	KEY_UP:         "up",
	KEY_DOWN:       "down",
	KEY_LEFT:       "left",
	KEY_RIGHT:      "right",
	KEY_CAPSLOCK:   "capslock",
	KEY_PRINT:      "print",
	KEY_SCROLLLOCK: "scrolllock",
	KEY_PAUSE:      "pause",

	// Punctuation
	KEY_MINUS:      "minus",
	KEY_EQUAL:      "equal",
	KEY_LEFTBRACE:  "leftbrace",
	KEY_RIGHTBRACE: "rightbrace",
	KEY_SEMICOLON:  "semicolon",
	KEY_APOSTROPHE: "apostrophe",
	KEY_GRAVE:      "grave",
	KEY_BACKSLASH:  "backslash",
	KEY_COMMA:      "comma",
	KEY_DOT:        "dot",
	KEY_SLASH:      "slash",

	// Media keys
	KEY_MUTE:           "mute",
	KEY_VOLUMEDOWN:     "volumedown",
	KEY_VOLUMEUP:       "volumeup",
	KEY_PLAYPAUSE:      "playpause",
	KEY_NEXTSONG:       "nextsong",
	KEY_PREVIOUSSONG:   "previoussong",
	KEY_STOPCD:         "stop",
	KEY_BRIGHTNESSUP:   "brightnessup",
	KEY_BRIGHTNESSDOWN: "brightnessdown",
}

// ModifierKeys for quick lookup
var ModifierKeys = map[string]uint16{
	"ctrl":       KEY_LEFTCTRL,
	"control":    KEY_LEFTCTRL,
	"leftctrl":   KEY_LEFTCTRL,
	"rightctrl":  KEY_RIGHTCTRL,
	"alt":        KEY_LEFTALT,
	"leftalt":    KEY_LEFTALT,
	"rightalt":   KEY_RIGHTALT,
	"shift":      KEY_LEFTSHIFT,
	"leftshift":  KEY_LEFTSHIFT,
	"rightshift": KEY_RIGHTSHIFT,
	"super":      KEY_LEFTMETA,
	"meta":       KEY_LEFTMETA,
	"win":        KEY_LEFTMETA,
	"leftmeta":   KEY_LEFTMETA,
	"rightmeta":  KEY_RIGHTMETA,
}

func IsModifier(keyStr string) bool {
	_, found := ModifierKeys[keyStr]
	return found
}

func LookupKeyCode(name string) (uint16, bool) {
	name = strings.ToLower(strings.TrimSpace(name))
	code, ok := KeyNameToCode[name]
	return code, ok
}

func LookupKeyName(code uint16) (string, bool) {
	if name, ok := KeyCodeToName[code]; ok {
		return name, true
	}
	return "unknown", false
}
