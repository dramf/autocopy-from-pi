package main

import (
	"etpribor.ru/autocopy/app"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/rand"
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

func main() {
	rand.Seed(12212112)
	flag.Parse()
	log.Print("Running ETP AutoCopy")
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
	tick := time.NewTicker(time.Millisecond * time.Duration(cfg.PollInterval))

	const checkMoreOften = 100
	checkMounter := time.NewTicker(time.Millisecond * time.Duration(cfg.PollInterval) * checkMoreOften)
	log.Print("Start listening of new USB flashes")

	if err := app.MountRemoteServer(cfg.UploadPath, cfg.LocalEndpoint); err != nil {
		log.Fatalf("MountRemoteServer fatal error: %v", err)
	}

	f := strings.TrimPrefix(cfg.Folder, "/")
	folder := fmt.Sprintf("%s/%s", cfg.LocalEndpoint , f)

	for {
		select {
		case <-tick.C:
			flashes := app.FlashDetector(&cfg.MountPrefix, &cfg.LocalEndpoint)
			for _, flash := range flashes {
				log.Printf("Start coping from %q to %q", flash, cfg.UploadPath)
				go app.CopyFolder(folder, flash)
			}
		case <-checkMounter.C:
			if err := app.MountRemoteServer(cfg.UploadPath, cfg.LocalEndpoint); err != nil {
				log.Fatalf("MountRemoteServer fatal error: %v", err)
			}
		}
	}
}
