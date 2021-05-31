package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func CopyFolder(remote, flash string) {
	files, err := ioutil.ReadDir(flash)
	if err != nil {
		log.Printf("readDir error: %v", err)
		return
	}

	for _, file := range files {
		log.Printf("start copy %q from %q to %q", file.Name(), flash, remote)
		_, err := copyFile(flash+"/"+file.Name(), remote+"/"+file.Name())
		if err != nil {
			log.Printf("Copy file %q err: %v", file.Name(), err)
			continue
		}
		log.Printf("File %q was copied succesfully!", file.Name())
	}
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
