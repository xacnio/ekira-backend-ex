package utils

import "time"

var TZ *time.Location = nil

func LoadTimezone() {
	loc, errTZ := time.LoadLocation("Europe/Istanbul")
	if errTZ != nil {
		panic(errTZ)
	}
	time.Local = loc
	TZ = loc
}
