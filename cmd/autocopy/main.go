package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dramf/autocopy/pkg/commands"
)

var (
	version   = "dev"
	buildtime = "1655367163"
)

func getVersion() string {
	bti, err := strconv.ParseInt(buildtime, 10, 64)
	if err != nil {
		log.Printf("parsing build time error: %v, will use a default value", err)
		bti = 1655367163
	}
	return fmt.Sprintf("v%s build: %s", version, time.Unix(bti, 0).Format("2006.01.02 15:04:05"))
}

func main() {
	app := commands.NewApp(getVersion())
	if err := app.Execute(); err != nil {
		log.Fatalf("execute error: %v", err)
	}
}
