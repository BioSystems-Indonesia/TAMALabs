package constant

import (
	"time"
)

var (
	LocationJakarta *time.Location = loadTimezone("Asia/Jakarta")
)

func loadTimezone(timezone string) *time.Location {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
	return location
}
