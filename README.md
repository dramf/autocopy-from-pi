# Autocopy from USB flash drives via Raspberry Pi
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