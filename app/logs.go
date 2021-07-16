package app

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CopyingLogs(folder, flash string) {
	log.Printf("Copying logs from %q/logs to %q/logs", flash, folder)
	files, err := ioutil.ReadDir(flash + "/logs")
	if err != nil {
		log.Printf("readDir error: %v", err)
		return
	}

	cmcFiles := []string{}
	lampNumber := "00000"
	err = filepath.Walk(flash, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasPrefix(path, flash+"/CMC15_") {
			cmcFiles = append(cmcFiles, path)
		}
		return nil
	})
	if err != nil {
		log.Printf("Search of CMC15_xxxxx files error: %v", err)
	}
	if len(cmcFiles) > 0 {
		lampNumber = getLampNumber(cmcFiles[0])
	}
	log.Printf("Using next lamp number: %q", lampNumber)

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
			log.Printf("creating log folder error: %v", err)
			continue
		}

		logFile := flash + "/logs/" + name
		source, err := ioutil.ReadFile(logFile)
		if err != nil {
			log.Printf("Reading log file %q error: %v", name, err)
			continue
		}

		f, err := os.OpenFile(logFolder+name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Printf("open remote log file %q error: %v", name, err)
			continue
		}

		if _, err = f.Write(source); err != nil {
			log.Printf("write remote log file %q error: %v", name, err)
		}
		f.Close()
		if err := os.Remove(logFile); err != nil {
			log.Printf("remove log file %q error: %v", logFile, err)
			continue
		}
		log.Printf("log %q was copied to %q and removed", name, logFolder)
	}
}
