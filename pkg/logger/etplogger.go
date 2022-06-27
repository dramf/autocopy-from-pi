package logger

import (
	"fmt"
	"os"

	"github.com/dramf/autocopy/pkg/utils"
)

type etpLogger struct {
	folder           string
	fileNameTemplate string
}

func getETPLogger(folder, fileTemplate string) *etpLogger {
	return &etpLogger{
		folder:           folder,
		fileNameTemplate: fileTemplate,
	}
}

func (logger *etpLogger) Write(data []byte) (int, error) {
	if err := os.MkdirAll(logger.folder, 0x666); err != nil {
		return 0, err
	}
	now := utils.GetCurrentDateName()
	fn := fmt.Sprintf(logger.fileNameTemplate, now)
	f, err := os.OpenFile(logger.folder+fn, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}
