package main

import (
	"fmt"
	"log"

	"github.com/glowfi/ghkd/internal/hotkey"
	evdev "github.com/gvalkov/golang-evdev"
)

const (
	EV_KEY = 1 // Event type for key events

	// Key states
	KEY_RELEASED = 0
	KEY_PRESSED  = 1
	KEY_REPEAT   = 2
)

func main() {
	devicePath := "/dev/input/event3" // Change this!

	device, err := evdev.Open(devicePath)
	if err != nil {
		log.Fatalf("Cannot open %s: %v", devicePath, err)
	}
	defer device.File.Close()

	fmt.Printf("Opened: %s\n", device.Name)
	fmt.Println("Press some keys (Ctrl+C to exit)...")
	fmt.Println()

	// Track currently pressed keys
	pressed := make(map[uint16]bool)

	for {
		events, err := device.Read()
		if err != nil {
			log.Fatalf("Read error: %v", err)
		}

		for _, ev := range events {
			// Only process key events
			if ev.Type != EV_KEY {
				continue
			}

			// Get state name
			state := "unknown"
			switch ev.Value {
			case KEY_RELEASED:
				state = "released"
				pressed[ev.Code] = true
			case KEY_PRESSED:
				state = "pressed"
				delete(pressed, ev.Code)
			case KEY_REPEAT:
				state = "repeat"
			}

			val, f := hotkey.LookupKeyName(ev.Code)
			if !f {
				continue
			}
			fmt.Printf("Key code=%s %s\n", val, state)
		}
	}
}
