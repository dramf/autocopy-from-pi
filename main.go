package main

import (
	"etpribor.ru/autocopy/app"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	configFile = ""
	ready      = make(chan struct{})
)

func init() {
	flag.StringVar(&configFile, "config", "settings.yml", "a path to the config yaml file")
}

const currentVersion = "v0.1.10"

func getLoggerForLamp(folder, lampNumber string) (*log.Logger, error) {
	now := app.GetCurrentDateName()
	dir := fmt.Sprintf("%s/%s/%s/logs", folder, now, lampNumber)
	if err := os.MkdirAll(dir, 0x666); err != nil {
		return nil, err
	}
	fn := fmt.Sprintf("%s/lamp_%s.txt", dir, now)
	f, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	w := io.MultiWriter(os.Stdout, f)
	logger := log.New(w, "["+lampNumber+" log] ", log.LstdFlags)
	log.Printf("The specific log file will be used for %s lamp: %q", lampNumber, fn)
	return logger, nil
}

func main() {
	rand.Seed(12212112)
	flag.Parse()
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	cfg := &app.Config{}
	if err := yaml.Unmarshal(b, cfg); err != nil {
		log.Fatal(err)
	}
	runner(cfg)
}

func runner(cfg *app.Config) {
	if err := app.MountRemoteServer(cfg.UploadPath, cfg.LocalEndpoint); err != nil {
		log.Fatalf("MountRemoteServer fatal error: %v", err)
	}
	folder := fmt.Sprintf("%s/%s/+logs", cfg.LocalEndpoint, strings.TrimPrefix(cfg.Folder, "/"))
	if err := os.MkdirAll(folder, 0666); err != nil {
		log.Fatalf("creating log folder (%q) error: %v", folder, err)
	}

	fn := fmt.Sprintf("%s/sys_%s.txt", folder, app.GetCurrentDateName())
	fMainLogs, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("setMainLogFile error: %v", err)
	}
	defer fMainLogs.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, fMainLogs))

	log.Printf("Running ETP AutoCopy %s", currentVersion)
	log.Printf("File for main logs: %q", fn)
	cfg.LogConfig()

	tick := time.NewTicker(time.Millisecond * time.Duration(cfg.PollInterval))
	const checkMoreOften = 100
	checkMounter := time.NewTicker(time.Millisecond * time.Duration(cfg.PollInterval) * checkMoreOften)

	log.Print("Start listening of new USB flashes")
	for {
		select {
		case <-tick.C:
			flashes := app.FlashDetector(&cfg.MountPrefix, &cfg.LocalEndpoint)
			for _, flash := range flashes {
				log.Printf("Mounted a new flash drive %q for copy to %q", flash, cfg.UploadPath)
				lampNumber := app.GetLampNumber(flash)
				logger, err := getLoggerForLamp(folder, lampNumber)
				if err != nil {
					log.Printf("[ERROR] Can't create a specific log for %q lamp (flash: %q): %v. The default logger will be used.", lampNumber, flash, err)
					logger = log.Default()
				}
				go app.CopyingMoviesFromFlash(folder, flash, true, logger)
				go app.CopyingLogs(folder, flash, lampNumber, logger)
			}
		case <-checkMounter.C:
			if err := app.MountRemoteServer(cfg.UploadPath, cfg.LocalEndpoint); err != nil {
				log.Fatalf("MountRemoteServer fatal error: %v", err)
			}
		}
	}
}
