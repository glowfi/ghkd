package listener

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"sync"

	"github.com/glowfi/ghkd/internal/hotkey"
	"github.com/holoplot/go-evdev"
)

const inputDir = "/dev/input"

type Listener struct {
	pressed []uint16
	devices []*evdev.InputDevice
	eventsC chan uint16
	wg      sync.WaitGroup
	mu      sync.RWMutex
}

func NewListener() *Listener {
	return &Listener{
		eventsC: make(chan uint16, 100),
	}
}

func (l *Listener) Start(ctx context.Context) error {
	keyboards, err := findKeyboards()
	if err != nil {
		return err
	}

	if len(keyboards) == 0 {
		return fmt.Errorf("no keyboards found")
	}

	for _, path := range keyboards {
		device, err := evdev.Open(path)
		if err != nil {
			return err
		}

		// Set non-blocking mode so Close() can interrupt ReadOne()
		if err := device.NonBlock(); err != nil {
			device.Close()
			return err
		}

		l.devices = append(l.devices, device)

		devName, err := device.Name()
		if err != nil {
			return err
		}

		fmt.Printf("Listening: %s\n", devName)
		go l.readDevice(ctx, device)
	}

	return nil
}

func (l *Listener) read(ctx context.Context, device *evdev.InputDevice) error {
	events, err := device.ReadSlice(1)
	if err != nil {
		return err
	}

	for _, ev := range events {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if ev.Type != hotkey.EV_KEY {
			continue
		}

		code := uint16(ev.Code)

		l.mu.Lock()
		switch ev.Value {
		case hotkey.KEY_PRESSED:
			l.pressed = append(l.pressed, code)
			// Notify about key press
			select {
			case l.eventsC <- code:
			default:
			}
		case hotkey.KEY_RELEASED:
			idx := slices.IndexFunc(l.pressed, func(code uint16) bool {
				return code == uint16(ev.Code)
			})
			if idx != -1 {
				l.pressed = append(l.pressed[:idx], l.pressed[idx+1:]...)
			}
		}
		l.mu.Unlock()
	}

	return nil
}

func (l *Listener) readDevice(ctx context.Context, device *evdev.InputDevice) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := l.read(ctx, device); err != nil {
				continue
			}
		}
	}
}

func findKeyboards() ([]string, error) {
	pattern := filepath.Join(inputDir, "event*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var keyboards []string
	for _, path := range matches {
		if isKeyboard(path) {
			keyboards = append(keyboards, path)
		}
	}
	return keyboards, nil
}

func isKeyboard(path string) bool {
	device, err := evdev.Open(path)
	if err != nil {
		return false
	}
	defer device.Close()

	// Get supported keys for EV_KEY type
	codes := device.CapableEvents(evdev.EV_KEY)

	// Check for letter keys (KEY_Q=16 to KEY_P=25)
	for _, code := range codes {
		if code >= 16 && code <= 25 {
			return true
		}
	}

	return false
}

func (l *Listener) PressedKeys() []uint16 {
	l.mu.RLock()
	defer l.mu.RUnlock()

	destination := make([]uint16, len(l.pressed))
	copy(destination, l.pressed)
	return destination
}

func (l *Listener) Events() <-chan uint16 {
	return l.eventsC
}

func (l *Listener) Stop() {
	for _, dev := range l.devices {
		dev.Close()
	}

	close(l.eventsC)
}
