package tray

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	"keychron-tray/internal/keychron"

	"github.com/getlantern/systray"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// 1. Global Variables (Required for state and menu items)
var (
	devices   []keychron.Device
	activeIdx int = -1

	mStatus  *systray.MenuItem
	mRefresh *systray.MenuItem
	mQuit    *systray.MenuItem

	chRefresh     = make(chan struct{}, 1)
	chQuit        = make(chan struct{}, 1)
	chDeviceClick = make(chan int, 1)
	systemDPI     = 96.0
	fontSize      = 18.0
)

//go:embed font.ttf
var fontData []byte

var txtFace font.Face

func init() {

	f, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatalf("Failed to parse font.ttf: %v", err)
	}

	txtFace, err = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     systemDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("Failed to create font face: %v", err)
	}
}

func calcBatteryColor(bat int) color.RGBA {
	const step = 5
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

// 4. Icon Generation (Text + Shadow + ICO Wrapper)
func generateBatteryIcon(bat int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))

	bg := calcBatteryColor(bat)
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	textColor := color.Black
	text := fmt.Sprintf("%d", bat)
	if text == "100" {
		text = "99"
	} // Fit optimization

	// Subtle WHITE shadow for contrast on dark green/red backgrounds
	// shadowColor := color.RGBA{255, 255, 255, 180}

	bounds, _ := font.BoundString(txtFace, text)
	textWidth := (bounds.Max.X - bounds.Min.X).Ceil()
	x := (32 - textWidth) / 2
	if x < 0 {
		x = 0
	}

	// Shadow (offset +1px)
	// dShadow := &font.Drawer{
	// 	Dst:  img,
	// 	Src:  image.NewUniform(shadowColor),
	// 	Face: txtFace,
	// 	Dot:  fixed.P(x+1, 25), // Adjusted Y for size 20
	// }
	// dShadow.DrawString(text)

	// Main Text
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(textColor),
		Face: txtFace,
		Dot:  fixed.P(x-1, 24), // Baseline for size 20
	}
	d.DrawString(text)

	return encodeICO32(img)
}

// 5. ICO Encoder (Memory-safe, no temp files)
func encodeICO32(img *image.RGBA) []byte {
	width, height := 32, 32
	buf := new(bytes.Buffer)

	// ICONDIR (6 bytes)
	binary.Write(buf, binary.LittleEndian, uint16(0)) // Reserved
	binary.Write(buf, binary.LittleEndian, uint16(1)) // Type: Icon
	binary.Write(buf, binary.LittleEndian, uint16(1)) // Count

	// ICONDIRENTRY (16 bytes)
	buf.WriteByte(byte(width))                         // Width
	buf.WriteByte(byte(height))                        // Height
	buf.WriteByte(0)                                   // Colors (0 for 32bit)
	buf.WriteByte(0)                                   // Reserved
	binary.Write(buf, binary.LittleEndian, uint16(1))  // Planes
	binary.Write(buf, binary.LittleEndian, uint16(32)) // BitsPerPixel

	// Size of bitmap info header + image data + AND mask
	// Header=40, Pixels=32*32*4, Mask=32*32/8
	dataSize := 40 + (width * height * 4) + (width * height / 8)
	binary.Write(buf, binary.LittleEndian, uint32(dataSize))
	binary.Write(buf, binary.LittleEndian, uint32(22)) // Offset to data

	// BITMAPINFOHEADER (40 bytes)
	binary.Write(buf, binary.LittleEndian, uint32(40))
	binary.Write(buf, binary.LittleEndian, uint32(width))
	binary.Write(buf, binary.LittleEndian, uint32(height*2)) // Height*2 for XOR+AND
	binary.Write(buf, binary.LittleEndian, uint16(1))
	binary.Write(buf, binary.LittleEndian, uint16(32))
	binary.Write(buf, binary.LittleEndian, uint32(0))           // Compression (BI_RGB)
	binary.Write(buf, binary.LittleEndian, uint32(dataSize-40)) // Image Size
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(0))

	// Pixel Data (BGRA, Bottom-Up)
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			buf.WriteByte(byte(b >> 8))
			buf.WriteByte(byte(g >> 8))
			buf.WriteByte(byte(r >> 8))
			buf.WriteByte(byte(a >> 8))
		}
	}

	// AND Mask (1 bit per pixel, 0=Opaque)
	maskSize := width * height / 8
	for i := 0; i < maskSize; i++ {
		buf.WriteByte(0x00)
	}

	return buf.Bytes()
}

// 6. Public Interface Functions

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

	// Update the generated icon
	UpdateTrayIcon(dev.Battery)
}

func UpdateTrayIcon(bat int) {
	iconData := generateBatteryIcon(bat)
	if iconData != nil {
		systray.SetIcon(iconData)
	}
}

func GetActiveDevice() keychron.Device {
	if activeIdx >= 0 && activeIdx < len(devices) {
		return devices[activeIdx]
	}
	return keychron.Device{}
}

func SetActiveDevice(idx int) { activeIdx = idx }

func RefreshChan() <-chan struct{} { return chRefresh }
func QuitChan() <-chan struct{}    { return chQuit }
func DeviceClickChan() <-chan int  { return chDeviceClick }

func shorten(s string) string {
	if len(s) > 24 {
		return s[:21] + "..."
	}
	return s
}
