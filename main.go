//go:generate goversioninfo -64
package main

import (
	"log"
	"time"

	"keychron-tray/internal/config"
	"keychron-tray/internal/keychron"
	"keychron-tray/internal/logger"
	"keychron-tray/internal/tray"

	"github.com/getlantern/systray"
)

func main() {
	logger.SetupLogging()
	keychron.Init()
	if keychron.InitError != nil {
		log.Fatalf("HID init: %v", keychron.InitError)
	}
	log.Printf("%s", keychron.Devices)
	systray.Run(onReady, nil)
}

func onReady() {
	devs, err := keychron.GetValidDevices()
	if err != nil {
		log.Printf("Enumeration error: %v", err)
	}

	tray.DrawMenu(devs)

	if active := tray.GetActiveDevice(); active.Product != "" {
		tray.UpdateAllUI(active)
	}

	go func() {
		ticker := time.NewTicker(config.PollIntervalSeconds * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			pollAndUpdate()
		}
	}()

	for {
		select {
		case idx := <-tray.ChDeviceClick:
			tray.SetActiveDevice(idx)
			tray.UpdateAllUI(tray.GetActiveDevice())
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
	log.Printf("🔄 Starting poll cycle")
	for i := range keychron.Devices {
		keychron.RefreshBattery(&keychron.Devices[i])
	}
	// Sync the updated devices to the tray module
	tray.SyncDevices(keychron.Devices)
	activeDev := tray.GetActiveDevice()
	log.Printf("📱 Active device after sync: %+v", activeDev)
	tray.UpdateAllUI(activeDev)
}

func refreshAndRebuild() {
	devs, _ := keychron.GetValidDevices()
	if len(devs) > 0 {
		keychron.Devices = devs
		if active := tray.GetActiveDevice(); active.Product == "" {
			tray.SetActiveDevice(0)
		}
	}
	tray.UpdateAllUI(tray.GetActiveDevice())
}
