package config

import (
	"errors"
	"log"
	"os"
	"runtime"
	"time"
)

func StopDisplayAnalyzer(started *AnalyzerStart) {
	if started == nil || started.cmd == nil || started.cmd.Process == nil {
		return
	}
	if runtime.GOOS == "windows" {
		if err := started.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			log.Printf("Could not stop display analyzer: %v", err)
		}
		return
	}

	if err := started.cmd.Process.Signal(os.Interrupt); err != nil && !errors.Is(err, os.ErrProcessDone) {
		if killErr := started.cmd.Process.Kill(); killErr != nil && !errors.Is(killErr, os.ErrProcessDone) {
			log.Printf("Could not stop display analyzer: %v", killErr)
		}
		return
	}

	time.Sleep(750 * time.Millisecond)
	if err := started.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
		log.Printf("Could not force-stop display analyzer: %v", err)
	}
}
