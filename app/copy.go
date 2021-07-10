package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
)

const maxFiles = 3

// getSubDirs returns data and code
func getSubDirs(filename string) (string, string) {
	d := strings.Split(filename, "_")
	if len(d) < 2 {
		return fmt.Sprintf("%d", rand.Int63n(100000000)), strconv.Itoa(rand.Intn(100000))
	}
	for i := 0; i < len(d); i++ {
		if d[i] == "" {
			d = append(d[:i], d[i+1:]...)
		}
	}

	if len(d[1]) < 8 {
		return fmt.Sprintf("%d", rand.Int63n(100000000)), d[0]
	}
	date := d[1][:8]
	return date, d[0]
}

func isFileExist(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

func prepareCopy(fullpath, remote string) {
	_, filename := path.Split(fullpath)
	date, code := getSubDirs(filename)

	remoteFolder := fmt.Sprintf("%s/%s/%s/", remote, date, code)

	if err := os.MkdirAll(remoteFolder, 0777); err != nil {
		log.Printf("prepareCopy error for folder %q: %v", remoteFolder, err)
		return
	}

	remoteFile := remoteFolder + filename
	newfile := remoteFile

	for i := 1; ; i++ {
		if !isFileExist(newfile) {
			break
		}
		log.Printf("File %q is already exist!", newfile)
		filenameLength := len(remoteFile)
		newfile = fmt.Sprintf("%s_%d%s", remoteFile[:filenameLength-4], i, remoteFile[filenameLength-4:])
	}

	log.Printf("start copy %q from %q to %q", filename, fullpath, newfile)
	_, err := copyFile(fullpath, newfile)
	if err != nil {
		log.Printf("Copy file %q to %q err: %v", fullpath, remoteFolder+filename, err)
		return
	}
	if err := os.Remove(fullpath); err != nil {
		log.Printf("Removing file %q err: %v", fullpath, err)
		return
	}
	log.Printf("File %q was copied and removed succesfully!", filename)
}

//func CopyFolder(remote, flash string) {
func CopyFolder(remote, flash string, isBaseLevel bool) {
	files, err := ioutil.ReadDir(flash)
	if err != nil {
		log.Printf("readDir error: %v", err)
		return
	}

	for _, file := range files {
		name := file.Name()
		if file.IsDir() {
			CopyFolder(remote, flash+"/"+name, false)
			continue
		}
		if !strings.HasSuffix(strings.ToLower(name), ".mp4") {
			continue
		}
		prepareCopy(flash+"/"+name, remote)
	}
	if isBaseLevel {
		fileUmount, err := os.Create(flash + "/umount")
		if err != nil {
			log.Printf("Creating 'umount'-file error: %v", err)
			return
		}
		fileUmount.Close()
		log.Print("'umount' was created successfully")
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
