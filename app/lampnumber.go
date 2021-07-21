package app

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

func GetLampNumber(flash string) string {
	cmcFiles := []string{}
	lampNumber := "00000"
	err := filepath.Walk(flash, func(path string, info fs.FileInfo, err error) error {
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
		s := strings.TrimSuffix(cmcFiles[0], ".txt")
		l := len(s)
		lampNumber = s[l-5:]
	}
	log.Printf("Using next lamp number: %q for %q", lampNumber, flash)
	return lampNumber
}
