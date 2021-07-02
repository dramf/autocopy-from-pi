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
var ch = make(chan struct{}, maxFiles)

// getSubDirs returns data and code
func getSubDirs(filename string) (string, string) {
	d := strings.Split(filename,"_" )
	if len(d) < 2 {
		return fmt.Sprintf("%d", rand.Int63n(100000000)), strconv.Itoa(rand.Intn(100000))
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

func prepareCopy(fullpath, remote string)  {
	defer func() {<- ch}()

	_, filename := path.Split(fullpath)
	date, code := getSubDirs(filename)

	remoteFolder := fmt.Sprintf("%s/%s/%s/", remote, date, code)

	if err := os.MkdirAll(remoteFolder,0777); err != nil {
		log.Printf("prepareCopy error for folder %q: %v", remoteFolder, err)
		return
	}

	remoteFile := remoteFolder+filename
	newfile := remoteFile

	for i:=1; ; i++ {
		if !isFileExist(newfile) {
			break
		}
		log.Printf("File %q is already exist!", newfile)
		s := strings.TrimSuffix(remoteFile, ".mp4")
		newfile = fmt.Sprintf("%s_%d.mp4", s, i)
	}

	log.Printf("start copy %q from %q to %q", filename, fullpath, newfile)
	_, err := copyFile(fullpath, newfile)
	if err != nil {
		log.Printf("Copy file %q to %q err: %v", fullpath, remoteFolder+filename, err)
		return
	}
	log.Printf("File %q was copied succesfully!", filename)
}

//func CopyFolder(remote, flash string) {
func CopyFolder(remote, flash string) {
	files, err := ioutil.ReadDir(flash)
	if err != nil {
		log.Printf("readDir error: %v", err)
		return
	}

	for _, file := range files {
		name := file.Name()
		if file.IsDir() {
			CopyFolder(remote, flash+"/"+name)
			continue
		}
		if !strings.HasSuffix(name, ".mp4") {
			continue
		}

		ch <- struct{}{}
		go prepareCopy(flash + "/" +name, remote)
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
