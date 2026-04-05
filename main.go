package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"fyne.io/systray"
	"github.com/karalabe/hid"
)

const (
	vidKeychron  = 0x3434
	pollInterval = 60 * time.Second
	maxMenuSlots = 10
)

type DeviceState struct {
	Path    string
	Product string
	Battery int
	LastOK  time.Time
}

var (
	stateMu    sync.Mutex
	devices    = make(map[string]*DeviceState)
	activePath string

	mQuit      *systray.MenuItem
	mSeparator *systray.MenuItem
	mDevices   [maxMenuSlots]*systray.MenuItem
)

func getBattery(info hid.DeviceInfo) (int, []byte, error) {
	dev, err := info.Open()
	if err != nil {
		return 0, nil, err
	}
	defer dev.Close()

	cmd := make([]byte, 64)
	cmd[0] = 0xB3
	cmd[1] = 0x06

	if _, err := dev.Write(cmd); err != nil {
		return 0, nil, err
	}

	time.Sleep(100 * time.Millisecond)

	resp := make([]byte, 65)
	n, err := dev.Read(resp)
	if err != nil {
		return 0, nil, err
	}
	if n == 0 {
		return 0, nil, fmt.Errorf("no response")
	}

	log.Printf("🔍 Raw response (%d bytes): % x", n, resp[:min(40, n)])

	// Skip leading Report ID if present (common with karalabe/hid)
	payload := resp
	if n > 0 && resp[0] == 0xB3 {
		payload = resp[1:]
		log.Printf("   → Skipped leading Report ID 0xB3, payload starts with: % x", payload[:min(20, len(payload))])
	}

	// Try several possible offsets in the payload
	// Python offset 20 usually becomes ~19 or 20 after skipping Report ID
	candidates := []int{18, 19, 20, 21, 22}
	for _, off := range candidates {
		if off < len(payload) {
			bat := int(payload[off])
			if bat >= 1 && bat <= 100 { // 0 is suspicious for battery
				log.Printf("🔋 Found plausible battery at payload offset %d → %d%%", off, bat)
				return bat, resp[:n], nil
			}
		}
	}

	// Fallback: try original Python logic on the raw buffer too
	if len(resp) > 20 {
		bat := int(resp[20])
		if bat >= 1 && bat <= 100 {
			log.Printf("🔋 Fallback raw[20] = %d%%", bat)
			return bat, resp[:n], nil
		}
	}

	return 0, resp[:n], fmt.Errorf("could not find battery value in response")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func pollDevices() {
	devs := hid.Enumerate(vidKeychron, 0)
	now := time.Now()
	currentPaths := make(map[string]bool)

	for _, d := range devs {
		bat, raw, err := getBattery(d)
		if d.ProductID != 0xD028 {
			continue
		}
		if err != nil {
			if raw != nil {
				log.Printf("⚠️  %s: %v (raw: %x)", d.Path, err, raw[:min(24, len(raw))])
			} else {
				log.Printf("⚠️  %s: %v", d.Path, err)
			}
			continue
		}

		stateMu.Lock()
		if state, ok := devices[d.Path]; ok {
			state.Battery = bat
			state.LastOK = now
		} else {
			prod := d.Product
			if prod == "" {
				prod = "Keychron Device"
			}
			devices[d.Path] = &DeviceState{
				Path:    d.Path,
				Product: prod,
				Battery: bat,
				LastOK:  now,
			}
		}
		currentPaths[d.Path] = true
		stateMu.Unlock()
	}

	stateMu.Lock()
	for path, state := range devices {
		if !currentPaths[path] || now.Sub(state.LastOK) > 120*time.Second {
			delete(devices, path)
			if activePath == path {
				activePath = ""
			}
		}
	}
	if activePath == "" {
		for path := range devices {
			activePath = path
			break
		}
	}
	stateMu.Unlock()

	updateTray()
}

func updateTray() {
	stateMu.Lock()
	defer stateMu.Unlock()

	if activePath == "" {
		systray.SetTitle("No Keychron")
		systray.SetTooltip("No responding devices")
		for _, m := range mDevices {
			if m != nil {
				m.SetTitle("---")
				m.Disable()
				m.Uncheck()
			}
		}
		return
	}

	info := devices[activePath]
	bat := info.Battery

	// TEXT-ONLY TRAY (no icons)
	systray.SetTitle(fmt.Sprintf("🔋%d%%", bat))
	systray.SetTooltip(fmt.Sprintf("Keychron %s\nBattery: %d%%\nPath: %s", info.Product, bat, info.Path))

	// Sync menu
	i := 0
	for path, state := range devices {
		if i >= maxMenuSlots {
			break
		}
		m := mDevices[i]
		if m != nil {
			m.SetTitle(fmt.Sprintf("%s (%d%%)", state.Product, state.Battery))
			m.Enable()
			if path == activePath {
				m.Check()
			} else {
				m.Uncheck()
			}
			i++
		}
	}
	for i < maxMenuSlots {
		if mDevices[i] != nil {
			mDevices[i].SetTitle("---")
			mDevices[i].Disable()
			mDevices[i].Uncheck()
		}
		i++
	}
}

func onReady() {
	// Redirect logs to console (remove -H=windowsgui to see output)
	log.SetOutput(os.Stdout)
	log.Println("🚀 Keychron Battery Tray starting...")

	systray.SetTitle("🔋--%")
	systray.SetTooltip("Querying devices...")

	mQuit = systray.AddMenuItem("Quit", "Exit application")
	systray.AddSeparator()

	for i := 0; i < maxMenuSlots; i++ {
		mDevices[i] = systray.AddMenuItem("---", "")
		mDevices[i].Disable()
	}

	// Menu click handlers
	for i := 0; i < maxMenuSlots; i++ {
		go func(idx int) {
			for {
				select {
				case <-mDevices[idx].ClickedCh:
					stateMu.Lock()
					count := 0
					for path := range devices {
						if count == idx {
							activePath = path
							stateMu.Unlock()
							updateTray()
							break
						}
						count++
					}
					if count <= idx {
						stateMu.Unlock()
					}
				}
			}
		}(i)
	}

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	// Initial poll + ticker
	pollDevices()
	go func() {
		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				pollDevices()
			}
		}
	}()
}

func onExit() {
	log.Println("👋 Keychron Tray exited")
}

func main() {
	systray.Run(onReady, onExit)
}
