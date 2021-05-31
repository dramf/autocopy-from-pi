# Autocopy from USB flash drives via Raspberry Pi
## Usage
1. Install
```shell
https://github.com/dramf/autocopy-from-pi/releases/download/v0.0.1/autocopy_linux_armv7.tar.gz
tar -xf autocopy_linux_armv7.tar.gz
```
2. Include the mount options for the cifs folder in the configuration file `/etc/fstab`:
```
//192.168.0.2/shared /home/pi/winserver cifs credentials=<Path To Credentials>,noauto,user 0 0
```
3. Run with config
```shell
./autocopy --config settings.yml
```
## Set up credentials
To be able to mount that folder as a normal user (without `sudo`), include the mount options for the cifs folder in the configuration file `/etc/fstab`, and add the options `noauto,user`
```
//192.168.0.2/shared /path/to/local/endpoint cifs credentials=/home/user/path/to/credentials/.passwd,noauto,user 0 0
```
`.passwd` contains a `username` and `password` from Windows Storage:
```
username=user
password=pwd
```