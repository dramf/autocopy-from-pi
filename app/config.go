package app

import (
	"log"
)

type Config struct {
	MountPrefix string `yaml:"mountOn"`
	Server      `yaml:"server"`
	TimeOuts    `yaml:"timeouts"`
}

type Server struct {
	UploadPath    string `yaml:"path"`
	Folder        string `yaml:"folder"`
	LocalEndpoint string `yaml:"local_endpoint"`
}

type TimeOuts struct {
	PollInterval int64 `yaml:"poll_interval"`
}

func (cfg *Config) LogConfig() {
	log.Print("Using next configuration:")
	log.Printf("Path for upload: %q", cfg.UploadPath)
	log.Printf("Poll interval: %d ms.", cfg.PollInterval)
	log.Printf("Endpoint for flash drives: %q", cfg.MountPrefix)
}
