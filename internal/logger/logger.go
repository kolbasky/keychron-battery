package logger

import (
	"io"
	"keychron-tray/internal/config"
	"log"
	"os"
	"path/filepath"
)

func SetupLogging() {
	exe, err := os.Executable()
	if err != nil {
		exe = os.Args[0]
	}
	logFile := filepath.Join(filepath.Dir(exe), "keybat.log")
	// logFile := "C:\\Users\\mtrag\\git\\keychron-battery\\bin\\keybat.log"

	_, err = os.Stat(logFile)
	fileExists := !os.IsNotExist(err)

	var output io.Writer
	if fileExists {
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			output = f
		} else {
			output = io.Discard
		}
	} else {
		output = io.Discard
	}

	if config.DEV_MODE {
		if output != nil {
			log.SetOutput(io.MultiWriter(os.Stdout, output))
		} else {
			log.SetOutput(os.Stdout)
		}
	} else {
		// if not in dev mode - log to file or do not log at all
		if output != nil {
			log.SetOutput(output)
		} else {
			log.SetOutput(io.Discard)
		}
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if fileExists {
		log.Println("📝 Logging to keybat.log")
	} else {
		log.Println("🔇 Logging to file disabled (no keybat.log found)")
	}
}
