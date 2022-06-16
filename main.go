package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"etpribor.ru/autocopy/app"
	"gopkg.in/yaml.v2"
)

var (
	version   = "dev"
	buildtime = fmt.Sprintf("%d", time.Now().Unix())

	configFile = ""
	ready      = make(chan struct{})
)

func init() {
	flag.StringVar(&configFile, "config", "settings.yml", "a path to the config yaml file")
}

func getLoggerForLamp(folder, lampNumber string) *log.Logger {
	now := app.GetCurrentDateName()
	dir := fmt.Sprintf("%s/%s/%s/logs", folder, now, lampNumber)
	logger := app.GetETPLogger(dir, "/lamp_%s.txt")
	w := io.MultiWriter(os.Stdout, logger)
	return log.New(w, "["+lampNumber+" log] ", log.LstdFlags)
}

func getMainWriter(mainFolder string) io.Writer {
	hostname, _ := os.Hostname()
	return io.MultiWriter(os.Stdout, app.GetETPLogger(mainFolder+"/+logs/"+hostname, "/sys_%s.txt"))
}

func main() {
	bti, err := strconv.ParseInt(buildtime, 10, 64)
	if err != nil {
		log.Fatalf("parsing build time error: %v", err)
	}
	log.Printf("Running ETP AutoCopy %s %s", version, time.Unix(bti, 0).Format("2006.01.02 15:04:05"))

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
	folder := fmt.Sprintf("%s/%s", cfg.LocalEndpoint, strings.TrimPrefix(cfg.Folder, "/"))
	log.SetOutput(getMainWriter(folder))

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
				logger := getLoggerForLamp(folder, lampNumber)
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
