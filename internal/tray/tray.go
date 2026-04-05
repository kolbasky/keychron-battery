package tray

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"path/filepath"

	"keychron-tray/internal/config"
	"keychron-tray/internal/keychron"

	"github.com/getlantern/systray"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var (
	devices   []keychron.Device
	activeIdx int = -1

	mRefresh *systray.MenuItem
	mQuit    *systray.MenuItem

	chRefresh     = make(chan struct{}, 1)
	chQuit        = make(chan struct{}, 1)
	ChDeviceClick = make(chan int, 1)

	txtFace font.Face
)

func init() {
	face, err := loadSystemFont()
	if err != nil {
		log.Fatalf("Failed to load system font: %v", err)
	}
	txtFace = face
	log.Println("✓ Using system Segoe UI font")
}

func loadSystemFont() (font.Face, error) {
	fontPaths := []string{
		filepath.Join(os.Getenv("WINDIR"), "Fonts", "segoeui.ttf"),
		filepath.Join(os.Getenv("WINDIR"), "Fonts", "seguisb.ttf"),
		"C:\\Windows\\Fonts\\segoeui.ttf",
	}

	var fontData []byte
	var err error

	for _, path := range fontPaths {
		fontData, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("could not load Segoe UI: %w", err)
	}

	f, err := opentype.Parse(fontData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	return opentype.NewFace(f, &opentype.FaceOptions{
		Size:    config.DefaultFontSize,
		DPI:     config.DefaultDPI,
		Hinting: font.HintingFull,
	})
}

// 100% is green [0,255,0], 50% is yellow [255,255,0], 0% is red [255,0,0]
// so we have to go from 0 to 255 in 50 steps = ~5 per step
// when going from 100 to 50 add +5 to red component to get more yellowish
// when going from 50 to 0 subtract -5 from green component to get more red
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

func generateBatteryIcon(bat int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))

	bg := calcBatteryColor(bat)
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	applyRoundedCorners(img, config.IconCornerRadius)

	// change 100% to 99% for better fit in small tray icon
	text := fmt.Sprintf("%d", bat)
	if text == "100" {
		text = "99"
	}
	textColor := color.Black

	bounds, _ := font.BoundString(txtFace, text)
	textWidth := (bounds.Max.X - bounds.Min.X).Ceil()
	x := (32 - textWidth) / 2
	if x < 0 {
		x = 0
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(textColor),
		Face: txtFace,
		Dot:  fixed.P(x-1, 24),
	}
	d.DrawString(text)

	return encodeICO32(img)
}

func encodeICO32(img *image.RGBA) []byte {
	width, height := 32, 32
	buf := new(bytes.Buffer)

	// ICONDIR (6 bytes)
	binary.Write(buf, binary.LittleEndian, uint16(0)) // Reserved
	binary.Write(buf, binary.LittleEndian, uint16(1)) // Type: Icon
	binary.Write(buf, binary.LittleEndian, uint16(1)) // Count

	// ICONDIRENTRY (16 bytes)
	buf.WriteByte(byte(width))
	buf.WriteByte(byte(height))
	buf.WriteByte(0)
	buf.WriteByte(0)
	binary.Write(buf, binary.LittleEndian, uint16(1))
	binary.Write(buf, binary.LittleEndian, uint16(32))

	dataSize := 40 + (width * height * 4) + (width * height / 8)
	binary.Write(buf, binary.LittleEndian, uint32(dataSize))
	binary.Write(buf, binary.LittleEndian, uint32(22)) // Offset to data

	// BITMAPINFOHEADER (40 bytes)
	binary.Write(buf, binary.LittleEndian, uint32(40))
	binary.Write(buf, binary.LittleEndian, uint32(width))
	binary.Write(buf, binary.LittleEndian, uint32(height*2))
	binary.Write(buf, binary.LittleEndian, uint16(1))
	binary.Write(buf, binary.LittleEndian, uint16(32))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	binary.Write(buf, binary.LittleEndian, uint32(dataSize-40))
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

// applyRoundedCorners makes the corners of a 32x32 RGBA image transparent
func applyRoundedCorners(img *image.RGBA, radius int) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	radiusSq := radius * radius

	// Corner centers - adjusted for pixel coordinates (0-indexed)
	// For radius R, the arc should touch the edge at pixel R-1
	tlCx, tlCy := radius-1, radius-1          // Top-left
	trCx, trCy := width-radius, radius-1      // Top-right
	blCx, blCy := radius-1, height-radius     // Bottom-left
	brCx, brCy := width-radius, height-radius // Bottom-right

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Check if pixel is in any corner region
			inCorner := false
			var cx, cy int

			if x < radius && y < radius {
				// Top-left corner region
				inCorner = true
				cx, cy = tlCx, tlCy
			} else if x >= width-radius && y < radius {
				// Top-right corner region
				inCorner = true
				cx, cy = trCx, trCy
			} else if x < radius && y >= height-radius {
				// Bottom-left corner region
				inCorner = true
				cx, cy = blCx, blCy
			} else if x >= width-radius && y >= height-radius {
				// Bottom-right corner region
				inCorner = true
				cx, cy = brCx, brCy
			}

			// If in a corner region, test against that corner's circle
			if inCorner {
				dx, dy := x-cx, y-cy
				if dx*dx+dy*dy > radiusSq {
					img.Set(x, y, color.RGBA{0, 0, 0, 0})
				}
			}
		}
	}
}

func DrawMenu(devs []keychron.Device) {
	devices = devs
	if len(devs) > 0 && activeIdx == -1 {
		activeIdx = 0
	}

	for i, dev := range devs {
		prefix := "  "
		if i == activeIdx {
			prefix = "✓ "
		}

		title := fmt.Sprintf("%s%s: %d%%", prefix, shorten(dev.Product), dev.Battery)
		item := systray.AddMenuItem(title, dev.Path)

		go func(idx int, itm *systray.MenuItem) {
			for range itm.ClickedCh {
				ChDeviceClick <- idx
			}
		}(i, item)
	}

	systray.AddSeparator()
	mRefresh = systray.AddMenuItem("🗘 Refresh", "") // ↻ ⟲ ⭮ ⭯ ↺ 🔃 🔁
	mQuit = systray.AddMenuItem("🗙 Quit", "")

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
		systray.SetTooltip("Keychron Battery Monitor\n⚠️ No device selected")
		return
	}

	tooltip := fmt.Sprintf(" %s \n🔋 %d%% ", dev.Product, dev.Battery)
	systray.SetTooltip(tooltip)

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
func DeviceClickChan() <-chan int  { return ChDeviceClick }

func shorten(s string) string {
	if len(s) > 24 {
		return s[:21] + "..."
	}
	return s
}
