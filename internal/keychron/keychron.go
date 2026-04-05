package keychron

import (
	"fmt"
	"time"

	"github.com/sstallion/go-hid"
)

const (
	VID_KEYCHRON   = 0x3434
	BATTERY_OFFSET = 20
	QUERY_TIMEOUT  = 150 * time.Millisecond
)

type Device struct {
	Path    string
	PID     uint16
	Product string
	Battery int
}

var Devices []Device
var InitError error

// GetBatteryData queries a single device path for battery %.
func GetBatteryData(path string) (int, error) {
	type result struct {
		bat int
		err error
	}
	done := make(chan result, 1)

	// Run blocking I/O in goroutine
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
		time.Sleep(QUERY_TIMEOUT)

		resp := make([]byte, 65)
		n, err := dev.Read(resp)
		if err != nil || n <= BATTERY_OFFSET {
			done <- result{0, fmt.Errorf("read failed or short")}
			return
		}

		bat := int(resp[BATTERY_OFFSET])
		if bat < 1 || bat > 100 {
			done <- result{0, fmt.Errorf("invalid battery: %d", bat)}
			return
		}
		done <- result{bat, nil}
	}()

	// Wait for result OR timeout
	select {
	case res := <-done:
		return res.bat, res.err
	case <-time.After(500 * time.Millisecond): // ← 500ms max per device
		return 0, fmt.Errorf("timeout")
	}
}

// GetValidDevices enumerates Keychron devices, filters non-responders.
func GetValidDevices() ([]Device, error) {
	// Step 1: Collect all Keychron device info (fast, no I/O)
	type devInfo struct {
		path    string
		pid     uint16
		product string
	}
	var all []devInfo

	_ = hid.Enumerate(VID_KEYCHRON, 0, func(info *hid.DeviceInfo) error {
		all = append(all, devInfo{info.Path, info.ProductID, info.ProductStr})
		return nil
	})

	// Step 2: Query each device sequentially (normal loop, easy to debug)
	var valid []Device
	for _, d := range all {
		bat, err := GetBatteryData(d.path)
		if err != nil {
			fmt.Printf("⊗ %s: %v\n", d.path, err)
			continue
		}
		fmt.Printf("✓ %s: %d%%\n", d.product, bat)
		valid = append(valid, Device{
			Path: d.path, PID: d.pid, Product: d.product, Battery: bat,
		})
	}

	fmt.Printf("📊 Found %d valid devices\n", len(valid))
	return valid, nil
}

// RefreshBattery updates Battery field of a Device struct.
func RefreshBattery(dev *Device) {
	if bat, err := GetBatteryData(dev.Path); err == nil {
		dev.Battery = bat
	}
}

func Init() {
	Devices, InitError = GetValidDevices()
}
