package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func CopyingLogs(folder, flash, lampNumber string, logger *log.Logger) {
	logger.Printf("Copying logs from %q/logs to %q/logs", flash, folder)
	files, err := ioutil.ReadDir(flash + "/logs")
	if err != nil {
		logger.Printf("readDir error: %v", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := strings.ToLower(file.Name())
		if !strings.HasSuffix(name, ".txt") {
			continue
		}

		logFolder := fmt.Sprintf("%s/20%s/%s/", folder, strings.TrimSuffix(name, ".txt"), lampNumber)
		if err := os.MkdirAll(logFolder, 0666); err != nil {
			logger.Printf("creating log folder error: %v", err)
			continue
		}

		logFile := flash + "/logs/" + name
		source, err := ioutil.ReadFile(logFile)
		if err != nil {
			logger.Printf("Reading log file %q error: %v", name, err)
			continue
		}

		f, err := os.OpenFile(logFolder+name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			logger.Printf("open remote log file %q error: %v", name, err)
			continue
		}

		if _, err = f.Write(source); err != nil {
			logger.Printf("write remote log file %q error: %v", name, err)
		}
		f.Close()
		if err := os.Remove(logFile); err != nil {
			logger.Printf("remove log file %q error: %v", logFile, err)
			continue
		}
		logger.Printf("log %q was copied to %q and removed", name, logFolder)
	}
}
