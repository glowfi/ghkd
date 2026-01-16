package hotkey

import (
	evdev "github.com/gvalkov/golang-evdev"
)

const (
	EV_KEY       = 1
	KEY_RELEASED = 0
	KEY_PRESSED  = 1
	KEY_REPEAT   = 2
)

const (
	// Standard keys
	KEY_ESC        = evdev.KEY_ESC
	KEY_1          = evdev.KEY_1
	KEY_2          = evdev.KEY_2
	KEY_3          = evdev.KEY_3
	KEY_4          = evdev.KEY_4
	KEY_5          = evdev.KEY_5
	KEY_6          = evdev.KEY_6
	KEY_7          = evdev.KEY_7
	KEY_8          = evdev.KEY_8
	KEY_9          = evdev.KEY_9
	KEY_0          = evdev.KEY_0
	KEY_MINUS      = evdev.KEY_MINUS
	KEY_EQUAL      = evdev.KEY_EQUAL
	KEY_BACKSPACE  = evdev.KEY_BACKSPACE
	KEY_TAB        = evdev.KEY_TAB
	KEY_Q          = evdev.KEY_Q
	KEY_W          = evdev.KEY_W
	KEY_E          = evdev.KEY_E
	KEY_R          = evdev.KEY_R
	KEY_T          = evdev.KEY_T
	KEY_Y          = evdev.KEY_Y
	KEY_U          = evdev.KEY_U
	KEY_I          = evdev.KEY_I
	KEY_O          = evdev.KEY_O
	KEY_P          = evdev.KEY_P
	KEY_LEFTBRACE  = evdev.KEY_LEFTBRACE
	KEY_RIGHTBRACE = evdev.KEY_RIGHTBRACE
	KEY_ENTER      = evdev.KEY_ENTER
	KEY_LEFTCTRL   = evdev.KEY_LEFTCTRL
	KEY_A          = evdev.KEY_A
	KEY_S          = evdev.KEY_S
	KEY_D          = evdev.KEY_D
	KEY_F          = evdev.KEY_F
	KEY_G          = evdev.KEY_G
	KEY_H          = evdev.KEY_H
	KEY_J          = evdev.KEY_J
	KEY_K          = evdev.KEY_K
	KEY_L          = evdev.KEY_L
	KEY_SEMICOLON  = evdev.KEY_SEMICOLON
	KEY_APOSTROPHE = evdev.KEY_APOSTROPHE
	KEY_GRAVE      = evdev.KEY_GRAVE
	KEY_LEFTSHIFT  = evdev.KEY_LEFTSHIFT
	KEY_BACKSLASH  = evdev.KEY_BACKSLASH
	KEY_Z          = evdev.KEY_Z
	KEY_X          = evdev.KEY_X
	KEY_C          = evdev.KEY_C
	KEY_V          = evdev.KEY_V
	KEY_B          = evdev.KEY_B
	KEY_N          = evdev.KEY_N
	KEY_M          = evdev.KEY_M
	KEY_COMMA      = evdev.KEY_COMMA
	KEY_DOT        = evdev.KEY_DOT
	KEY_SLASH      = evdev.KEY_SLASH
	KEY_RIGHTSHIFT = evdev.KEY_RIGHTSHIFT
	KEY_LEFTALT    = evdev.KEY_LEFTALT
	KEY_SPACE      = evdev.KEY_SPACE
	KEY_CAPSLOCK   = evdev.KEY_CAPSLOCK
	KEY_F1         = evdev.KEY_F1
	KEY_F2         = evdev.KEY_F2
	KEY_F3         = evdev.KEY_F3
	KEY_F4         = evdev.KEY_F4
	KEY_F5         = evdev.KEY_F5
	KEY_F6         = evdev.KEY_F6
	KEY_F7         = evdev.KEY_F7
	KEY_F8         = evdev.KEY_F8
	KEY_F9         = evdev.KEY_F9
	KEY_F10        = evdev.KEY_F10
	KEY_F11        = evdev.KEY_F11
	KEY_F12        = evdev.KEY_F12
	KEY_RIGHTCTRL  = evdev.KEY_RIGHTCTRL
	KEY_RIGHTALT   = evdev.KEY_RIGHTALT
	KEY_HOME       = evdev.KEY_HOME
	KEY_UP         = evdev.KEY_UP
	KEY_PAGEUP     = evdev.KEY_PAGEUP
	KEY_LEFT       = evdev.KEY_LEFT
	KEY_RIGHT      = evdev.KEY_RIGHT
	KEY_END        = evdev.KEY_END
	KEY_DOWN       = evdev.KEY_DOWN
	KEY_PAGEDOWN   = evdev.KEY_PAGEDOWN
	KEY_INSERT     = evdev.KEY_INSERT
	KEY_DELETE     = evdev.KEY_DELETE
	KEY_LEFTMETA   = evdev.KEY_LEFTMETA
	KEY_RIGHTMETA  = evdev.KEY_RIGHTMETA
	KEY_SCROLLLOCK = evdev.KEY_SCROLLLOCK
	KEY_NUMLOCK    = evdev.KEY_NUMLOCK
	KEY_PRINT      = evdev.KEY_SYSRQ
	KEY_PAUSE      = evdev.KEY_PAUSE

	// Media keys (kernel names)
	KEY_MUTE           = evdev.KEY_MUTE
	KEY_VOLUMEDOWN     = evdev.KEY_VOLUMEDOWN
	KEY_VOLUMEUP       = evdev.KEY_VOLUMEUP
	KEY_NEXTSONG       = evdev.KEY_NEXTSONG
	KEY_PREVIOUSSONG   = evdev.KEY_PREVIOUSSONG
	KEY_PLAYPAUSE      = evdev.KEY_PLAYPAUSE
	KEY_PAUSECD        = evdev.KEY_PAUSECD
	KEY_STOPCD         = evdev.KEY_STOPCD
	KEY_BRIGHTNESSDOWN = evdev.KEY_BRIGHTNESSDOWN
	KEY_BRIGHTNESSUP   = evdev.KEY_BRIGHTNESSUP

	// Additional media/launcher keys
	KEY_CALC   = evdev.KEY_CALC
	KEY_MAIL   = evdev.KEY_MAIL
	KEY_SEARCH = evdev.KEY_SEARCH
	KEY_FILE   = evdev.KEY_FILE
	KEY_WWW    = evdev.KEY_WWW
	KEY_COFFEE = evdev.KEY_COFFEE // Screen lock
	KEY_SLEEP  = evdev.KEY_SLEEP
	KEY_WAKEUP = evdev.KEY_WAKEUP
	KEY_POWER  = evdev.KEY_POWER
)
