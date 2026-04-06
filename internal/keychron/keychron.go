package keychron

import (
	"fmt"
	"log"
	"time"

	"keychron-tray/internal/config"

	"github.com/sstallion/go-hid"
)

type Device struct {
	Path      string
	PID       uint16
	Product   string
	Battery   int
	Connected bool
}

var Devices []Device
var InitError error

func GetBatteryData(path string) (int, error) {
	type result struct {
		bat int
		err error
	}
	done := make(chan result, 1)

	go func() {
		dev, err := hid.OpenPath(path)
		if err != nil {
			done <- result{0, err}
			return
		}
		defer dev.Close()

		cmd := [64]byte{0xB3, 0x06}
		if _, err := dev.Write(cmd[:]); err != nil {
			done <- result{0, err}
			return
		}
		time.Sleep(config.QueryTimeout)

		resp := make([]byte, 65)
		n, err := dev.Read(resp)
		if err != nil || n <= config.BatteryOffset {
			done <- result{0, fmt.Errorf("read failed or short")}
			return
		}

		bat := int(resp[config.BatteryOffset])
		if bat < 1 || bat > 100 {
			done <- result{0, fmt.Errorf("invalid battery: %d", bat)}
			return
		}
		done <- result{bat, nil}
	}()

	select {
	case res := <-done:
		return res.bat, res.err
	case <-time.After(500 * time.Millisecond): // ← 500ms max per device
		return 0, fmt.Errorf("timeout")
	}
}

func GetValidDevices() ([]Device, error) {
	type devInfo struct {
		path    string
		pid     uint16
		product string
	}
	var all []devInfo

	_ = hid.Enumerate(config.VidKeychron, 0, func(info *hid.DeviceInfo) error {
		all = append(all, devInfo{info.Path, info.ProductID, info.ProductStr})
		return nil
	})

	var valid []Device
	for _, d := range all {
		bat, err := GetBatteryData(d.path)
		if err != nil {
			fmt.Printf("🗙 %s: %v\n", d.path, err)
			continue
		}
		fmt.Printf("✓ %s: %d%%\n", d.product, bat)
		valid = append(valid, Device{
			Path: d.path, PID: d.pid, Product: d.product, Battery: bat, Connected: true,
		})
	}

	fmt.Printf("📊 Found %d valid devices\n", len(valid))
	return valid, nil
}

func RefreshBattery(dev *Device) {
	log.Printf("🔍 Polling device: %s (%s)", dev.Product, dev.Path)
	if bat, err := GetBatteryData(dev.Path); err == nil {
		log.Printf("✅ Poll success: %d%%", bat)
		dev.Battery = bat
		dev.Connected = true
	} else {
		log.Printf("❌ Poll failed: %v", err)
		dev.Connected = false
	}
}

func Init() {
	Devices, InitError = GetValidDevices()
}
