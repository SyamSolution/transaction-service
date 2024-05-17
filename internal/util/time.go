package util

import (
	"os"
	"time"
)

var TimeNow = time.Now

func LoadLocation() *time.Location {
	zone := os.Getenv("TIMEZONE")
	if zone == "" {
		zone = "Asia/Jakarta"
	}

	loc, err := time.LoadLocation(zone)
	if err != nil {
		return time.Local
	}

	return loc
}
