package config

import "time"

const (
	DEV_MODE = true

	VidKeychron   = 0x3434
	BatteryOffset = 20
	QueryTimeout  = 150 * time.Millisecond

	DefaultDPI       = 96.0
	DefaultFontSize  = 18.0
	IconCornerRadius = 8
)
