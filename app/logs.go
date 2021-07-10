package app

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func CopyingLogs(folder, flash string) {
	log.Printf("Copying logs from %q/logs to %q/logs", flash, folder)
	files, err := ioutil.ReadDir(flash + "/logs")
	if err != nil {
		log.Printf("readDir error: %v", err)
		return
	}
	if err := os.MkdirAll(folder+"/logs", 0777); err != nil {
		log.Printf("creating log folder error: %v", err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".txt") {
			continue
		}

		source, err := ioutil.ReadFile(flash + "/logs/" + name)
		if err != nil {
			log.Printf("Reading log file %q error: %v", name, err)
			continue
		}

		f, err := os.OpenFile(folder+"/logs/"+name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Printf("open remote log file %q error: %v", name, err)
			continue
		}

		if _, err = f.Write(source); err != nil {
			log.Printf("write remote log file %q error: %v", name, err)
		}
		f.Close()
	}
}
