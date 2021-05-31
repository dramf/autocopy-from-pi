package app

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func IsMounted(endpoint string) bool {
	res, _ := exec.Command("sh", "-c", fmt.Sprintf("mount | grep %s", endpoint)).Output()
	if len(res) > 0 {
		return true
	}
	return false
}

func MountRemoteServer(path, localFolder string) error {
	if IsMounted(localFolder) {
		return nil
	}
	if localFolder == "" {
		var err error
		localFolder, err = os.MkdirTemp("", "*")
		if err != nil {
			return err
		}
	}
	cmd := "mount.cifs"
	args := []string{
		path,
		localFolder,
	}
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Printf("Auth to %q error! output: %s", path, out)
		log.Printf("Probably, you should include the mount options for the cifs folder in the configuration file /etc/fstab: "+
			"%s %s cifs credentials=<Path To Credentials>,noauto,user 0 0", path, localFolder)
		return err
	}
	log.Printf("MS Server mounted to %q", localFolder)
	return nil
}
