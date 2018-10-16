package go_pretty_print

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var durationUnits = [8]struct {
	unit string
	one  Duration
}{
	{"w", Duration(7 * 24 * time.Hour)},
	{"d", Duration(24 * time.Hour)},
	{"h", Duration(time.Hour)},
	{"m", Duration(time.Minute)},
	{"s", Duration(time.Second)},
	{"ms", Duration(time.Millisecond)},
	{"us", Duration(time.Microsecond)},
	{"ns", Duration(time.Nanosecond)},
}

type Duration time.Duration

func (dur Duration) MarshalJSON() ([]byte, error) {
	return []byte(dur.floatString('g', -1)), nil
}

func (dur Duration) String() string {
	return dur.string(2)
}

func (dur Duration) Format(f fmt.State, c rune) {
	switch c {
	case 'b', 'e', 'E', 'f', 'g', 'G':
		prec, hasPrec := f.Precision()
		if !hasPrec {
			prec = -1
		}

		fmt.Fprint(f, dur.floatString(byte(c), prec))
	case 's':
		prec, hasPrec := f.Precision()
		if !hasPrec {
			prec = 1
		}

		fmt.Fprint(f, dur.string(uint8(prec+1)))
	case 'v':
		fmt.Fprintf(f, "go_pretty_print.Duration(%v)", time.Duration(dur))
	default:
		fmt.Fprintf(f, "%%!%c(go_pretty_print.Duration=%s)", c, time.Duration(dur))
	}
}

func (dur Duration) floatString(fmt byte, prec int) string {
	return strconv.FormatFloat(float64(dur)/float64(time.Second), fmt, prec, 64)
}

func (dur Duration) string(units uint8) string {
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
