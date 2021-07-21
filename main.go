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

const currentVersion = "v0.1.8"

func setMainLogFile(dir string) error {
	fn := fmt.Sprintf("%s/sys_%s.txt", dir, app.GetCurrentDateName())
	fMainLogs, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0x666)
	if err != nil {
		return err
	}
	log.SetOutput(io.MultiWriter(os.Stdout, fMainLogs))
	log.Printf("File for main logs: %q", fn)
	return nil
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
	cfg.LogConfig()

	runner(cfg)
}

func runner(cfg *app.Config) {
	if err := app.MountRemoteServer(cfg.UploadPath, cfg.LocalEndpoint); err != nil {
		log.Fatalf("MountRemoteServer fatal error: %v", err)
	}
	folder := fmt.Sprintf("%s/%s", cfg.LocalEndpoint, strings.TrimPrefix(cfg.Folder, "/"))

	if err := setMainLogFile(folder); err != nil {
		log.Fatalf("setMainLogFile error: %v", err)
	}

	log.Printf("Running ETP AutoCopy %s", currentVersion)

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
				go app.CopyingMoviesFromFlash(folder, flash, true)
				go app.CopyingLogs(folder, flash)
			}
		case <-checkMounter.C:
			if err := app.MountRemoteServer(cfg.UploadPath, cfg.LocalEndpoint); err != nil {
				log.Fatalf("MountRemoteServer fatal error: %v", err)
			}
		}
	}
}
