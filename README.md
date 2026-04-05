# Keychron Battery Monitor

A lightweight portable Windows system tray utility that displays real-time battery percentage for Keychron wireless mice. Features a dynamic, color-coded tray icon, automatic polling, and zero background noise.

## ✅ Supported Devices
- Keychron M2 8K (tested with 8K 2.4G USB dongle)
- May work with other Keychron mice

> 💡 Battery reporting may not work in Bluetooth or Wired modes on all models. For reliable readings, keep the device in **2.4G Wireless Dongle** mode.

## 📦 Installation
1. Download `keybat.exe` from [Releases](https://github.com/kolbasky/keychron-battery/releases)
2. Copy it to any folder
3. Double-click to run — no installer, no admin rights required


### 🔍 Debug Logging (Optional)
Logging is disabled by default. To enable it, create an empty `keybat.log` file in the same folder as the executable:
```powershell
New-Item keybat.log -ItemType File
```
