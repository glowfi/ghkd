package listener

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/glowfi/ghkd/internal/hotkey"
	evdev "github.com/gvalkov/golang-evdev"
)

const inputDir = "/dev/input"

type Listener struct {
	pressed []uint16
	eventsC chan uint16
	mu      sync.RWMutex
	wg      sync.WaitGroup
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
			continue
		}
		if strings.Contains(strings.ToLower(device.Name), "mouse") {
			continue
		}
		fmt.Printf("Listening: %s\n", device.Name)
		l.wg.Add(1)
		go func() {
			defer l.wg.Done()
			l.readDevice(ctx, device)
		}()
	}

	return nil
}

func (l *Listener) read(device *evdev.InputDevice) {
	events, err := device.Read()
	if err != nil {
		return
	}

	for _, ev := range events {
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
				return code == ev.Code
			})
			if idx != -1 {
				l.pressed = append(l.pressed[:idx], l.pressed[idx+1:]...)
			}
		}
		l.mu.Unlock()
	}
}

func (l *Listener) readDevice(ctx context.Context, device *evdev.InputDevice) {
	defer device.File.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			l.read(device)
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
	defer device.File.Close()

	for capType, codes := range device.Capabilities {
		if capType.Type != evdev.EV_KEY {
			continue
		}
		for _, code := range codes {
			if code.Code >= 16 && code.Code <= 25 {
				return true
			}
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
	close(l.eventsC)
	l.wg.Wait()
}
