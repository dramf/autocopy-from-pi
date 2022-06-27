package logger

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dramf/autocopy/pkg/utils"
)

func InitLoggerForLamp(folder, lampNumber string) *log.Logger {
	now := utils.GetCurrentDateName()
	dir := fmt.Sprintf("%s/%s/%s/logs", folder, now, lampNumber)
	logger := getETPLogger(dir, "/lamp_%s.txt")
	w := io.MultiWriter(os.Stdout, logger)
	return log.New(w, "["+lampNumber+" log] ", log.LstdFlags)
}
