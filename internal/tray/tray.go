package tray

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	"keychron-tray/internal/keychron"

	"github.com/getlantern/systray"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var (
	devices       []keychron.Device
	activeIdx     int = -1
	mStatus       *systray.MenuItem
	mRefresh      *systray.MenuItem
	mQuit         *systray.MenuItem
	chRefresh     = make(chan struct{}, 1)
	chQuit        = make(chan struct{}, 1)
	chDeviceClick = make(chan int, 1)
)

// calcBatteryColor returns RGB per your spec:
// 100->50: Red increases, Green=255
// 50->1: Red=255, Green decreases
func calcBatteryColor(bat int) color.RGBA {
	const step = 5 // 255/50 = 5.1
	if bat >= 50 {
		r := uint8((100 - bat) * step)
		if r > 255 {
			r = 255
		}
		return color.RGBA{R: r, G: 255, B: 0, A: 255}
	}
	g := uint8(bat * step)
	if g > 255 {
		g = 255
	}
	return color.RGBA{R: 255, G: g, B: 0, A: 255}
}

// encodeICO32 creates a valid Windows ICO file in memory from an image.
// It includes the required BITMAPINFOHEADER which was missing before.
func encodeICO32(img *image.RGBA) []byte {
	width, height := 32, 32
	buf := new(bytes.Buffer)

	// 1. ICONDIR (6 bytes)
	binary.Write(buf, binary.LittleEndian, uint16(0)) // Reserved
	binary.Write(buf, binary.LittleEndian, uint16(1)) // Type: 1 (Icon)
	binary.Write(buf, binary.LittleEndian, uint16(1)) // Count: 1 image

	// 2. ICONDIRENTRY (16 bytes)
	buf.WriteByte(byte(width))                         // Width
	buf.WriteByte(byte(height))                        // Height
	buf.WriteByte(0)                                   // Color count (0 for 32bit)
	buf.WriteByte(0)                                   // Reserved
	binary.Write(buf, binary.LittleEndian, uint16(1))  // Color Planes
	binary.Write(buf, binary.LittleEndian, uint16(32)) // Bits Per Pixel

	// Image Data Size: 40 (Header) + Pixels + Mask
	dataSize := 40 + (width * height * 4) + (width * height / 8)
	binary.Write(buf, binary.LittleEndian, uint32(dataSize))

	// Offset to Image Data (6 bytes ICONDIR + 16 bytes Entry = 22)
	binary.Write(buf, binary.LittleEndian, uint32(22))

	// 3. BITMAPINFOHEADER (40 bytes) - CRITICAL FOR WINDOWS
	binary.Write(buf, binary.LittleEndian, uint32(40)) // Header Size
	binary.Write(buf, binary.LittleEndian, uint32(width))
	binary.Write(buf, binary.LittleEndian, uint32(height*2))                      // Double height for XOR+AND
	binary.Write(buf, binary.LittleEndian, uint16(1))                             // Planes
	binary.Write(buf, binary.LittleEndian, uint16(32))                            // Bits
	binary.Write(buf, binary.LittleEndian, uint32(0))                             // Compression (BI_RGB)
	binary.Write(buf, binary.LittleEndian, uint32(width*height*4+width*height/8)) // Image Size
	binary.Write(buf, binary.LittleEndian, uint32(0))                             // X PPM
	binary.Write(buf, binary.LittleEndian, uint32(0))                             // Y PPM
	binary.Write(buf, binary.LittleEndian, uint32(0))                             // Colors Used
	binary.Write(buf, binary.LittleEndian, uint32(0))                             // Colors Important

	// 4. XOR Pixels (BGRA, Bottom-Up)
	// Go images are Top-Down, Windows ICO is Bottom-Up
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// Convert 0-65535 to 0-255
			buf.WriteByte(byte(b >> 8)) // B
			buf.WriteByte(byte(g >> 8)) // G
			buf.WriteByte(byte(r >> 8)) // R
			buf.WriteByte(byte(a >> 8)) // A
		}
	}

	// 5. AND Mask (1 bit per pixel, 0=Opaque, 1=Transparent)
	// We want opaque background, so we write all 0x00.
	maskSize := width * height / 8
	for i := 0; i < maskSize; i++ {
		buf.WriteByte(0x00)
	}

	return buf.Bytes()
}

// generateBatteryIcon creates the image and encodes it to ICO
func generateBatteryIcon(bat int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))

	// Background
	bg := calcBatteryColor(bat)
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	// Text
	textColor := color.Black
	if bat >= 45 && bat <= 55 {
		textColor = color.Black
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(textColor),
		Face: basicfont.Face7x13,
	}
	text := fmt.Sprintf("%d%%", bat)
	textWidth := len(text) * 7
	x := (32 - textWidth) / 2
	if x < 0 {
		x = 0
	}
	d.Dot = fixed.P(x, 22)
	d.DrawString(text)

	return encodeICO32(img)
}

// UpdateTrayIcon sets the icon directly from memory (No Temp Files!)
func UpdateTrayIcon(bat int) {
	iconData := generateBatteryIcon(bat)
	if iconData != nil {
		systray.SetIcon(iconData)
	} else {
		log.Printf("⚠️ Failed to generate icon")
	}
}

func DrawMenu(devs []keychron.Device) {
	devices = devs
	if len(devs) > 0 && activeIdx == -1 {
		activeIdx = 0
	}

	mStatus = systray.AddMenuItem("🔋 --", "")
	systray.AddSeparator()
	systray.AddMenuItem("🖱️ Devices", "").Disable()
	systray.AddSeparator()

	for i, dev := range devs {
		title := fmt.Sprintf("  %s: %d%%", shorten(dev.Product), dev.Battery)
		item := systray.AddMenuItem(title, dev.Path)

		go func(idx int, itm *systray.MenuItem) {
			for range itm.ClickedCh {
				chDeviceClick <- idx
			}
		}(i, item)
	}

	systray.AddSeparator()
	mRefresh = systray.AddMenuItem("🔄 Refresh", "")
	mQuit = systray.AddMenuItem("❌ Quit", "")

	go func() {
		for range mRefresh.ClickedCh {
			chRefresh <- struct{}{}
		}
	}()
	go func() {
		for range mQuit.ClickedCh {
			chQuit <- struct{}{}
		}
	}()
}

func UpdateAllUI(dev keychron.Device) {
	if dev.Product == "" {
		mStatus.SetTitle("🔋 Show battery")
		systray.SetTooltip("Keychron Battery Monitor\n⚠️ No device selected")
		return
	}

	title := fmt.Sprintf("🔋 %s: %d%%", shorten(dev.Product), dev.Battery)
	tooltip := fmt.Sprintf("%s\n🔋 %d%%", dev.Product, dev.Battery)
	mStatus.SetTitle(title)
	systray.SetTooltip(tooltip)

	// Update the dynamic icon
	UpdateTrayIcon(dev.Battery)
}

func GetActiveDevice() keychron.Device {
	if activeIdx >= 0 && activeIdx < len(devices) {
		return devices[activeIdx]
	}
	return keychron.Device{}
}

func SetActiveDevice(idx int)      { activeIdx = idx }
func RefreshChan() <-chan struct{} { return chRefresh }
func QuitChan() <-chan struct{}    { return chQuit }
func DeviceClickChan() <-chan int  { return chDeviceClick }

func shorten(s string) string {
	if len(s) > 24 {
		return s[:21] + "..."
	}
	return s
}
