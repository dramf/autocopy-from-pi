package app

import "time"

func GetCurrentDateName() string {
	return time.Now().Format("20060102")
}
