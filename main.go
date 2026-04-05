package main

import (
	"log"
	"time"

	"keychron-tray/internal/keychron"
	"keychron-tray/internal/tray"

	"github.com/getlantern/systray"
)

//go:generate go-bindata -pkg main -o icon.go icon.ico

var testIcon []byte

func main() {
	keychron.Init()
	if keychron.InitError != nil {
		log.Fatalf("HID init: %v", keychron.InitError)
	}
	log.Printf("%s", keychron.Devices)
	systray.Run(onReady, nil)
}

func onReady() {
	// Load icon
	// if data, err := Asset("icon.ico"); err == nil {
	// 	systray.SetIcon(data)
	// }

	// 1. Get valid devices
	devs, err := keychron.GetValidDevices()
	if err != nil {
		log.Printf("Enumeration error: %v", err)
	}

	// 2. Build menu linearly (main thread)
	tray.DrawMenu(devs)

	time.Sleep(200 * time.Millisecond)

	// 3. Update UI with default active device
	if active := tray.GetActiveDevice(); active.Product != "" {
		tray.UpdateAllUI(active)
	}

	// 4. Polling loop (background)
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			pollAndUpdate()
		}
	}()

	// 5. Main event loop (handles channels)
	for {
		select {
		case <-tray.RefreshChan():
			refreshAndRebuild()
		case <-tray.QuitChan():
			systray.Quit()
			return
		case idx := <-tray.DeviceClickChan():
			tray.SetActiveDevice(idx)
			tray.UpdateAllUI(tray.GetActiveDevice())
		}
	}
}

func pollAndUpdate() {
	// Update all device structs
	for i := range keychron.Devices {
		keychron.RefreshBattery(&keychron.Devices[i])
		// Update menu titles safely
		// (In production, you'd map devices to menu items;
		// for brevity, we'll skip dynamic title sync here)
	}
	// Refresh tooltip/status
	tray.UpdateAllUI(tray.GetActiveDevice())
}

func refreshAndRebuild() {
	devs, _ := keychron.GetValidDevices()
	// Note: systray doesn't support removing items cleanly.
	// For a production app, you'd update existing MenuItem titles instead.
	// Here we just update tooltip to reflect fresh data.
	if len(devs) > 0 {
		keychron.Devices = devs // replace slice
		if active := tray.GetActiveDevice(); active.Product == "" {
			tray.SetActiveDevice(0)
		}
	}
	tray.UpdateAllUI(tray.GetActiveDevice())
}
