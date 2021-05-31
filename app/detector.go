package app

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

var flashes = make(map[string]bool)

func FlashDetector(pref, endpoint *string) []string {
	cmd := "df"
	args := []string{"--output=target"}

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Fatal(err)
	}

	newFlashes := []string{}

	for _, mountOn := range bytes.Split(out, []byte("\n")) {
		flash := string(mountOn)
		if !strings.HasPrefix(flash, *pref) {
			continue
		}
		if flash == *endpoint {
			continue
		}
		_, ok := flashes[flash]
		if ok {
			continue
		}
		log.Printf("Found a new flash! Endpoint: %q", flash)
		flashes[flash] = true
		newFlashes = append(newFlashes, flash)
	}
	return newFlashes
}
