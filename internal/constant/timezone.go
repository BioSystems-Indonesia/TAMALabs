package constant

import (
	"time"
)

// these specifies commonly used timezone
var (
	IndonesiaWesternStandardTimezone *time.Location
	IndonesiaCentralStandardTimezone *time.Location
	IndonesiaEasternStandardTimezone *time.Location
)

func init() {
	utcToWIB := int((7 * time.Hour).Seconds())
	IndonesiaWesternStandardTimezone = time.FixedZone("WIB", utcToWIB)

	utcToWITA := int((8 * time.Hour).Seconds())
	IndonesiaCentralStandardTimezone = time.FixedZone("WITA", utcToWITA)

	utcToWIT := int((9 * time.Hour).Seconds())
	IndonesiaEasternStandardTimezone = time.FixedZone("WIT", utcToWIT)
}

func loadTimezone(timezone string) *time.Location {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
	return location
}
