package go_pretty_print

import (
	"strconv"
	"strings"
	"time"
)

var durationUnits = [8]struct {
	unit string
	one  time.Duration
}{
	{"w", 7 * 24 * time.Hour},
	{"d", 24 * time.Hour},
	{"h", time.Hour},
	{"m", time.Minute},
	{"s", time.Second},
	{"ms", time.Millisecond},
	{"us", time.Microsecond},
	{"ns", time.Nanosecond},
}

func Duration(dur time.Duration, units uint8) string {
	if dur == 0 {
		return "0s"
	}

	negative := dur < 0
	if negative {
		dur = -dur
	}

	largestUnit := uint8(7)

	for i, unit := range durationUnits {
		if dur >= unit.one {
			largestUnit = uint8(i)
			break
		}
	}

	segments := []string{}

	for i := largestUnit; i < 8 && units > 0; i++ {
		unit := durationUnits[i]
		amount := dur / unit.one
		dur %= unit.one

		if amount > 0 {
			segments = append(segments, strconv.FormatInt(int64(amount), 10)+unit.unit)
		}

		units--
	}

	result := strings.Join(segments, " ")

	if negative {
		result = "-" + result;
	}

	return result
}
