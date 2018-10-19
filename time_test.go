package go_pretty_print

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"testing"
	"time"
)

const ns8 = 8 * time.Nanosecond
const us7 = 7 * time.Microsecond
const ms6 = 6 * time.Millisecond
const s5 = 5 * time.Second
const m4 = 4 * time.Minute
const h3 = 3 * time.Hour
const d2 = 2 * 24 * time.Hour
const w1 = 7 * 24 * time.Hour

func TestDuration_MarshalJSON(t *testing.T) {
	assertDuration_MarshalJSON(t, 0, "0")

	assertDuration_MarshalJSON(t, ns8, "8e-09")
	assertDuration_MarshalJSON(t, us7, "7e-06")
	assertDuration_MarshalJSON(t, ms6, "0.006")
	assertDuration_MarshalJSON(t, s5, "5")
	assertDuration_MarshalJSON(t, m4, "240")
	assertDuration_MarshalJSON(t, h3, "10800")
	assertDuration_MarshalJSON(t, d2, "172800")
	assertDuration_MarshalJSON(t, w1, "604800")

	assertDuration_MarshalJSON(t, -ns8, "-8e-09")
	assertDuration_MarshalJSON(t, -us7, "-7e-06")
	assertDuration_MarshalJSON(t, -ms6, "-0.006")
	assertDuration_MarshalJSON(t, -s5, "-5")
	assertDuration_MarshalJSON(t, -m4, "-240")
	assertDuration_MarshalJSON(t, -h3, "-10800")
	assertDuration_MarshalJSON(t, -d2, "-172800")
	assertDuration_MarshalJSON(t, -w1, "-604800")
}

func assertDuration_MarshalJSON(t *testing.T, d time.Duration, expected string) {
	t.Helper()

	jsn, err := json.Marshal(Duration(d))

	assertCallResult(
		t,
		fmt.Sprintf("json.Marshal(Duration(%v))", d),
		[]interface{}{[]byte(expected), nil},
		[]interface{}{jsn, err},
	)
}

func TestDuration_String(t *testing.T) {
	assertDuration_String(t, 0, "0s")

	assertDuration_String(t, ns8, "8ns")
	assertDuration_String(t, us7+ns8, "7us 8ns")
	assertDuration_String(t, ms6+us7+ns8, "6ms 7us")
	assertDuration_String(t, s5+us7+ns8, "5s 7us")
	assertDuration_String(t, m4, "4m")
	assertDuration_String(t, h3, "3h")
	assertDuration_String(t, d2, "2d")
	assertDuration_String(t, w1, "1w")

	assertDuration_String(t, -ns8, "-8ns")
	assertDuration_String(t, -us7-ns8, "-7us 8ns")
	assertDuration_String(t, -ms6-us7-ns8, "-6ms 7us")
	assertDuration_String(t, -s5-us7-ns8, "-5s 7us")
	assertDuration_String(t, -m4, "-4m")
	assertDuration_String(t, -h3, "-3h")
	assertDuration_String(t, -d2, "-2d")
	assertDuration_String(t, -w1, "-1w")
}

func assertDuration_String(t *testing.T, d time.Duration, expected string) {
	t.Helper()

	assertCallResult(
		t,
		fmt.Sprintf("Duration(%v).String()", d),
		[]interface{}{expected},
		[]interface{}{Duration(d).String()},
	)
}

func TestDuration_Format_0(t *testing.T) {
	assertDuration_Format(t, 0, "%b", "0p-1074")
	assertDuration_Format(t, 0, "%.2b", "0p-1074")
	assertDuration_Format(t, 0, "%e", "0e+00")
	assertDuration_Format(t, 0, "%.2e", "0.00e+00")
	assertDuration_Format(t, 0, "%E", "0E+00")
	assertDuration_Format(t, 0, "%.2E", "0.00E+00")
	assertDuration_Format(t, 0, "%f", "0")
	assertDuration_Format(t, 0, "%.2f", "0.00")
	assertDuration_Format(t, 0, "%g", "0")
	assertDuration_Format(t, 0, "%.2g", "0")
	assertDuration_Format(t, 0, "%G", "0")
	assertDuration_Format(t, 0, "%.2G", "0")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, 0, "%"+c, "0s")
		assertDuration_Format(t, 0, "%.0"+c, "0s")
		assertDuration_Format(t, 0, "%.7"+c, "0s")
		assertDuration_Format(t, 0, "%.128"+c, "0s")
	}

	assertDuration_Format(t, 0, "%d", "%!d(go_pretty_print.Duration=0s)")
}

func TestDuration_Format_NS(t *testing.T) {
	assertDuration_Format(t, ns8, "%b", "4835703278458517p-79")
	assertDuration_Format(t, ns8, "%.2b", "4835703278458517p-79")
	assertDuration_Format(t, ns8, "%e", "8e-09")
	assertDuration_Format(t, ns8, "%.2e", "8.00e-09")
	assertDuration_Format(t, ns8, "%E", "8E-09")
	assertDuration_Format(t, ns8, "%.2E", "8.00E-09")
	assertDuration_Format(t, ns8, "%f", "0.000000008")
	assertDuration_Format(t, ns8, "%.2f", "0.00")
	assertDuration_Format(t, ns8, "%g", "8e-09")
	assertDuration_Format(t, ns8, "%.2g", "8e-09")
	assertDuration_Format(t, ns8, "%G", "8E-09")
	assertDuration_Format(t, ns8, "%.2G", "8E-09")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, ns8, "%"+c, "8ns")
		assertDuration_Format(t, ns8, "%.0"+c, "8ns")
		assertDuration_Format(t, ns8, "%.7"+c, "8ns")
		assertDuration_Format(t, ns8, "%.128"+c, "8ns")
	}

	assertDuration_Format(t, ns8, "%d", "%!d(go_pretty_print.Duration=8ns)")

	assertDuration_Format(t, -ns8, "%b", "-4835703278458517p-79")
	assertDuration_Format(t, -ns8, "%.2b", "-4835703278458517p-79")
	assertDuration_Format(t, -ns8, "%e", "-8e-09")
	assertDuration_Format(t, -ns8, "%.2e", "-8.00e-09")
	assertDuration_Format(t, -ns8, "%E", "-8E-09")
	assertDuration_Format(t, -ns8, "%.2E", "-8.00E-09")
	assertDuration_Format(t, -ns8, "%f", "-0.000000008")
	assertDuration_Format(t, -ns8, "%.2f", "-0.00")
	assertDuration_Format(t, -ns8, "%g", "-8e-09")
	assertDuration_Format(t, -ns8, "%.2g", "-8e-09")
	assertDuration_Format(t, -ns8, "%G", "-8E-09")
	assertDuration_Format(t, -ns8, "%.2G", "-8E-09")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, -ns8, "%"+c, "-8ns")
		assertDuration_Format(t, -ns8, "%.0"+c, "-8ns")
		assertDuration_Format(t, -ns8, "%.7"+c, "-8ns")
		assertDuration_Format(t, -ns8, "%.128"+c, "-8ns")
	}

	assertDuration_Format(t, -ns8, "%d", "%!d(go_pretty_print.Duration=-8ns)")
}

func TestDuration_Format_US(t *testing.T) {
	assertDuration_Format(t, us7, "%b", "8264141345021879p-70")
	assertDuration_Format(t, us7, "%.2b", "8264141345021879p-70")
	assertDuration_Format(t, us7, "%e", "7e-06")
	assertDuration_Format(t, us7, "%.2e", "7.00e-06")
	assertDuration_Format(t, us7, "%E", "7E-06")
	assertDuration_Format(t, us7, "%.2E", "7.00E-06")
	assertDuration_Format(t, us7, "%f", "0.000007")
	assertDuration_Format(t, us7, "%.2f", "0.00")
	assertDuration_Format(t, us7, "%g", "7e-06")
	assertDuration_Format(t, us7, "%.2g", "7e-06")
	assertDuration_Format(t, us7, "%G", "7E-06")
	assertDuration_Format(t, us7, "%.2G", "7E-06")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, us7, "%"+c, "7us")
		assertDuration_Format(t, us7, "%.0"+c, "7us")
		assertDuration_Format(t, us7, "%.7"+c, "7us")
		assertDuration_Format(t, us7, "%.128"+c, "7us")

		assertDuration_Format(t, us7+ns8, "%"+c, "7us 8ns")
		assertDuration_Format(t, us7+ns8, "%.0"+c, "7us")
		assertDuration_Format(t, us7+ns8, "%.1"+c, "7us 8ns")
		assertDuration_Format(t, us7+ns8, "%.7"+c, "7us 8ns")
		assertDuration_Format(t, us7+ns8, "%.128"+c, "7us 8ns")
	}

	assertDuration_Format(t, us7, "%d", "%!d(go_pretty_print.Duration=7µs)")

	assertDuration_Format(t, -us7, "%b", "-8264141345021879p-70")
	assertDuration_Format(t, -us7, "%.2b", "-8264141345021879p-70")
	assertDuration_Format(t, -us7, "%e", "-7e-06")
	assertDuration_Format(t, -us7, "%.2e", "-7.00e-06")
	assertDuration_Format(t, -us7, "%E", "-7E-06")
	assertDuration_Format(t, -us7, "%.2E", "-7.00E-06")
	assertDuration_Format(t, -us7, "%f", "-0.000007")
	assertDuration_Format(t, -us7, "%.2f", "-0.00")
	assertDuration_Format(t, -us7, "%g", "-7e-06")
	assertDuration_Format(t, -us7, "%.2g", "-7e-06")
	assertDuration_Format(t, -us7, "%G", "-7E-06")
	assertDuration_Format(t, -us7, "%.2G", "-7E-06")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, -us7, "%"+c, "-7us")
		assertDuration_Format(t, -us7, "%.0"+c, "-7us")
		assertDuration_Format(t, -us7, "%.7"+c, "-7us")
		assertDuration_Format(t, -us7, "%.128"+c, "-7us")

		assertDuration_Format(t, -us7-ns8, "%"+c, "-7us 8ns")
		assertDuration_Format(t, -us7-ns8, "%.0"+c, "-7us")
		assertDuration_Format(t, -us7-ns8, "%.1"+c, "-7us 8ns")
		assertDuration_Format(t, -us7-ns8, "%.7"+c, "-7us 8ns")
		assertDuration_Format(t, -us7-ns8, "%.128"+c, "-7us 8ns")
	}

	assertDuration_Format(t, -us7, "%d", "%!d(go_pretty_print.Duration=-7µs)")
}

func TestDuration_Format_MS(t *testing.T) {
	assertDuration_Format(t, ms6, "%b", "6917529027641082p-60")
	assertDuration_Format(t, ms6, "%.2b", "6917529027641082p-60")
	assertDuration_Format(t, ms6, "%e", "6e-03")
	assertDuration_Format(t, ms6, "%.2e", "6.00e-03")
	assertDuration_Format(t, ms6, "%E", "6E-03")
	assertDuration_Format(t, ms6, "%.2E", "6.00E-03")
	assertDuration_Format(t, ms6, "%f", "0.006")
	assertDuration_Format(t, ms6, "%.2f", "0.01")
	assertDuration_Format(t, ms6, "%g", "0.006")
	assertDuration_Format(t, ms6, "%.2g", "0.006")
	assertDuration_Format(t, ms6, "%G", "0.006")
	assertDuration_Format(t, ms6, "%.2G", "0.006")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, ms6, "%"+c, "6ms")
		assertDuration_Format(t, ms6, "%.0"+c, "6ms")
		assertDuration_Format(t, ms6, "%.7"+c, "6ms")
		assertDuration_Format(t, ms6, "%.128"+c, "6ms")

		assertDuration_Format(t, ms6+us7, "%"+c, "6ms 7us")
		assertDuration_Format(t, ms6+us7, "%.0"+c, "6ms")
		assertDuration_Format(t, ms6+us7, "%.1"+c, "6ms 7us")
		assertDuration_Format(t, ms6+us7, "%.7"+c, "6ms 7us")
		assertDuration_Format(t, ms6+us7, "%.128"+c, "6ms 7us")

		assertDuration_Format(t, ms6+ns8, "%"+c, "6ms 8ns")
		assertDuration_Format(t, ms6+ns8, "%.0"+c, "6ms")
		assertDuration_Format(t, ms6+ns8, "%.1"+c, "6ms 8ns")
		assertDuration_Format(t, ms6+ns8, "%.7"+c, "6ms 8ns")
		assertDuration_Format(t, ms6+ns8, "%.128"+c, "6ms 8ns")

		assertDuration_Format(t, ms6+us7+ns8, "%"+c, "6ms 7us")
		assertDuration_Format(t, ms6+us7+ns8, "%.0"+c, "6ms")
		assertDuration_Format(t, ms6+us7+ns8, "%.1"+c, "6ms 7us")
		assertDuration_Format(t, ms6+us7+ns8, "%.2"+c, "6ms 7us 8ns")
		assertDuration_Format(t, ms6+us7+ns8, "%.7"+c, "6ms 7us 8ns")
		assertDuration_Format(t, ms6+us7+ns8, "%.128"+c, "6ms 7us 8ns")
	}

	assertDuration_Format(t, ms6, "%d", "%!d(go_pretty_print.Duration=6ms)")

	assertDuration_Format(t, -ms6, "%b", "-6917529027641082p-60")
	assertDuration_Format(t, -ms6, "%.2b", "-6917529027641082p-60")
	assertDuration_Format(t, -ms6, "%e", "-6e-03")
	assertDuration_Format(t, -ms6, "%.2e", "-6.00e-03")
	assertDuration_Format(t, -ms6, "%E", "-6E-03")
	assertDuration_Format(t, -ms6, "%.2E", "-6.00E-03")
	assertDuration_Format(t, -ms6, "%f", "-0.006")
	assertDuration_Format(t, -ms6, "%.2f", "-0.01")
	assertDuration_Format(t, -ms6, "%g", "-0.006")
	assertDuration_Format(t, -ms6, "%.2g", "-0.006")
	assertDuration_Format(t, -ms6, "%G", "-0.006")
	assertDuration_Format(t, -ms6, "%.2G", "-0.006")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, -ms6, "%"+c, "-6ms")
		assertDuration_Format(t, -ms6, "%.0"+c, "-6ms")
		assertDuration_Format(t, -ms6, "%.7"+c, "-6ms")
		assertDuration_Format(t, -ms6, "%.128"+c, "-6ms")

		assertDuration_Format(t, -ms6-us7, "%"+c, "-6ms 7us")
		assertDuration_Format(t, -ms6-us7, "%.0"+c, "-6ms")
		assertDuration_Format(t, -ms6-us7, "%.1"+c, "-6ms 7us")
		assertDuration_Format(t, -ms6-us7, "%.7"+c, "-6ms 7us")
		assertDuration_Format(t, -ms6-us7, "%.128"+c, "-6ms 7us")

		assertDuration_Format(t, -ms6-ns8, "%"+c, "-6ms 8ns")
		assertDuration_Format(t, -ms6-ns8, "%.0"+c, "-6ms")
		assertDuration_Format(t, -ms6-ns8, "%.1"+c, "-6ms 8ns")
		assertDuration_Format(t, -ms6-ns8, "%.7"+c, "-6ms 8ns")
		assertDuration_Format(t, -ms6-ns8, "%.128"+c, "-6ms 8ns")

		assertDuration_Format(t, -ms6-us7-ns8, "%"+c, "-6ms 7us")
		assertDuration_Format(t, -ms6-us7-ns8, "%.0"+c, "-6ms")
		assertDuration_Format(t, -ms6-us7-ns8, "%.1"+c, "-6ms 7us")
		assertDuration_Format(t, -ms6-us7-ns8, "%.2"+c, "-6ms 7us 8ns")
		assertDuration_Format(t, -ms6-us7-ns8, "%.7"+c, "-6ms 7us 8ns")
		assertDuration_Format(t, -ms6-us7-ns8, "%.128"+c, "-6ms 7us 8ns")
	}

	assertDuration_Format(t, -ms6, "%d", "%!d(go_pretty_print.Duration=-6ms)")
}

func TestDuration_Format_S(t *testing.T) {
	assertDuration_Format(t, s5, "%b", "5629499534213120p-50")
	assertDuration_Format(t, s5, "%.2b", "5629499534213120p-50")
	assertDuration_Format(t, s5, "%e", "5e+00")
	assertDuration_Format(t, s5, "%.2e", "5.00e+00")
	assertDuration_Format(t, s5, "%E", "5E+00")
	assertDuration_Format(t, s5, "%.2E", "5.00E+00")
	assertDuration_Format(t, s5, "%f", "5")
	assertDuration_Format(t, s5, "%.2f", "5.00")
	assertDuration_Format(t, s5, "%g", "5")
	assertDuration_Format(t, s5, "%.2g", "5")
	assertDuration_Format(t, s5, "%G", "5")
	assertDuration_Format(t, s5, "%.2G", "5")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, s5, "%"+c, "5s")
		assertDuration_Format(t, s5, "%.0"+c, "5s")
		assertDuration_Format(t, s5, "%.7"+c, "5s")
		assertDuration_Format(t, s5, "%.128"+c, "5s")

		assertDuration_Format(t, s5+ms6, "%"+c, "5s 6ms")
		assertDuration_Format(t, s5+ms6, "%.0"+c, "5s")
		assertDuration_Format(t, s5+ms6, "%.1"+c, "5s 6ms")
		assertDuration_Format(t, s5+ms6, "%.7"+c, "5s 6ms")
		assertDuration_Format(t, s5+ms6, "%.128"+c, "5s 6ms")

		assertDuration_Format(t, s5+ns8, "%"+c, "5s 8ns")
		assertDuration_Format(t, s5+ns8, "%.0"+c, "5s")
		assertDuration_Format(t, s5+ns8, "%.1"+c, "5s 8ns")
		assertDuration_Format(t, s5+ns8, "%.7"+c, "5s 8ns")
		assertDuration_Format(t, s5+ns8, "%.128"+c, "5s 8ns")
	}

	assertDuration_Format(t, s5, "%d", "%!d(go_pretty_print.Duration=5s)")

	assertDuration_Format(t, -s5, "%b", "-5629499534213120p-50")
	assertDuration_Format(t, -s5, "%.2b", "-5629499534213120p-50")
	assertDuration_Format(t, -s5, "%e", "-5e+00")
	assertDuration_Format(t, -s5, "%.2e", "-5.00e+00")
	assertDuration_Format(t, -s5, "%E", "-5E+00")
	assertDuration_Format(t, -s5, "%.2E", "-5.00E+00")
	assertDuration_Format(t, -s5, "%f", "-5")
	assertDuration_Format(t, -s5, "%.2f", "-5.00")
	assertDuration_Format(t, -s5, "%g", "-5")
	assertDuration_Format(t, -s5, "%.2g", "-5")
	assertDuration_Format(t, -s5, "%G", "-5")
	assertDuration_Format(t, -s5, "%.2G", "-5")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, -s5, "%"+c, "-5s")
		assertDuration_Format(t, -s5, "%.0"+c, "-5s")
		assertDuration_Format(t, -s5, "%.7"+c, "-5s")
		assertDuration_Format(t, -s5, "%.128"+c, "-5s")

		assertDuration_Format(t, -s5-ms6, "%"+c, "-5s 6ms")
		assertDuration_Format(t, -s5-ms6, "%.0"+c, "-5s")
		assertDuration_Format(t, -s5-ms6, "%.1"+c, "-5s 6ms")
		assertDuration_Format(t, -s5-ms6, "%.7"+c, "-5s 6ms")
		assertDuration_Format(t, -s5-ms6, "%.128"+c, "-5s 6ms")

		assertDuration_Format(t, -s5-ns8, "%"+c, "-5s 8ns")
		assertDuration_Format(t, -s5-ns8, "%.0"+c, "-5s")
		assertDuration_Format(t, -s5-ns8, "%.1"+c, "-5s 8ns")
		assertDuration_Format(t, -s5-ns8, "%.7"+c, "-5s 8ns")
		assertDuration_Format(t, -s5-ns8, "%.128"+c, "-5s 8ns")
	}

	assertDuration_Format(t, -s5, "%d", "%!d(go_pretty_print.Duration=-5s)")
}

func TestDuration_Format_M(t *testing.T) {
	assertDuration_Format(t, m4, "%b", "8444249301319680p-45")
	assertDuration_Format(t, m4, "%.2b", "8444249301319680p-45")
	assertDuration_Format(t, m4, "%e", "2.4e+02")
	assertDuration_Format(t, m4, "%.2e", "2.40e+02")
	assertDuration_Format(t, m4, "%E", "2.4E+02")
	assertDuration_Format(t, m4, "%.2E", "2.40E+02")
	assertDuration_Format(t, m4, "%f", "240")
	assertDuration_Format(t, m4, "%.2f", "240.00")
	assertDuration_Format(t, m4, "%g", "240")
	assertDuration_Format(t, m4, "%.2g", "2.4e+02")
	assertDuration_Format(t, m4, "%G", "240")
	assertDuration_Format(t, m4, "%.2G", "2.4E+02")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, m4, "%"+c, "4m")
		assertDuration_Format(t, m4, "%.0"+c, "4m")
		assertDuration_Format(t, m4, "%.7"+c, "4m")
		assertDuration_Format(t, m4, "%.128"+c, "4m")

		assertDuration_Format(t, m4+s5, "%"+c, "4m 5s")
		assertDuration_Format(t, m4+s5, "%.0"+c, "4m")
		assertDuration_Format(t, m4+s5, "%.1"+c, "4m 5s")
		assertDuration_Format(t, m4+s5, "%.7"+c, "4m 5s")
		assertDuration_Format(t, m4+s5, "%.128"+c, "4m 5s")

		assertDuration_Format(t, m4+ns8, "%"+c, "4m 8ns")
		assertDuration_Format(t, m4+ns8, "%.0"+c, "4m")
		assertDuration_Format(t, m4+ns8, "%.1"+c, "4m 8ns")
		assertDuration_Format(t, m4+ns8, "%.7"+c, "4m 8ns")
		assertDuration_Format(t, m4+ns8, "%.128"+c, "4m 8ns")
	}

	assertDuration_Format(t, m4, "%d", "%!d(go_pretty_print.Duration=4m0s)")

	assertDuration_Format(t, -m4, "%b", "-8444249301319680p-45")
	assertDuration_Format(t, -m4, "%.2b", "-8444249301319680p-45")
	assertDuration_Format(t, -m4, "%e", "-2.4e+02")
	assertDuration_Format(t, -m4, "%.2e", "-2.40e+02")
	assertDuration_Format(t, -m4, "%E", "-2.4E+02")
	assertDuration_Format(t, -m4, "%.2E", "-2.40E+02")
	assertDuration_Format(t, -m4, "%f", "-240")
	assertDuration_Format(t, -m4, "%.2f", "-240.00")
	assertDuration_Format(t, -m4, "%g", "-240")
	assertDuration_Format(t, -m4, "%.2g", "-2.4e+02")
	assertDuration_Format(t, -m4, "%G", "-240")
	assertDuration_Format(t, -m4, "%.2G", "-2.4E+02")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, -m4, "%"+c, "-4m")
		assertDuration_Format(t, -m4, "%.0"+c, "-4m")
		assertDuration_Format(t, -m4, "%.7"+c, "-4m")
		assertDuration_Format(t, -m4, "%.128"+c, "-4m")

		assertDuration_Format(t, -m4-s5, "%"+c, "-4m 5s")
		assertDuration_Format(t, -m4-s5, "%.0"+c, "-4m")
		assertDuration_Format(t, -m4-s5, "%.1"+c, "-4m 5s")
		assertDuration_Format(t, -m4-s5, "%.7"+c, "-4m 5s")
		assertDuration_Format(t, -m4-s5, "%.128"+c, "-4m 5s")

		assertDuration_Format(t, -m4-ns8, "%"+c, "-4m 8ns")
		assertDuration_Format(t, -m4-ns8, "%.0"+c, "-4m")
		assertDuration_Format(t, -m4-ns8, "%.1"+c, "-4m 8ns")
		assertDuration_Format(t, -m4-ns8, "%.7"+c, "-4m 8ns")
		assertDuration_Format(t, -m4-ns8, "%.128"+c, "-4m 8ns")
	}

	assertDuration_Format(t, -m4, "%d", "%!d(go_pretty_print.Duration=-4m0s)")
}

func TestDuration_Format_H(t *testing.T) {
	assertDuration_Format(t, h3, "%b", "5937362789990400p-39")
	assertDuration_Format(t, h3, "%.2b", "5937362789990400p-39")
	assertDuration_Format(t, h3, "%e", "1.08e+04")
	assertDuration_Format(t, h3, "%.2e", "1.08e+04")
	assertDuration_Format(t, h3, "%E", "1.08E+04")
	assertDuration_Format(t, h3, "%.2E", "1.08E+04")
	assertDuration_Format(t, h3, "%f", "10800")
	assertDuration_Format(t, h3, "%.2f", "10800.00")
	assertDuration_Format(t, h3, "%g", "10800")
	assertDuration_Format(t, h3, "%.2g", "1.1e+04")
	assertDuration_Format(t, h3, "%G", "10800")
	assertDuration_Format(t, h3, "%.2G", "1.1E+04")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, h3, "%"+c, "3h")
		assertDuration_Format(t, h3, "%.0"+c, "3h")
		assertDuration_Format(t, h3, "%.7"+c, "3h")
		assertDuration_Format(t, h3, "%.128"+c, "3h")

		assertDuration_Format(t, h3+m4, "%"+c, "3h 4m")
		assertDuration_Format(t, h3+m4, "%.0"+c, "3h")
		assertDuration_Format(t, h3+m4, "%.1"+c, "3h 4m")
		assertDuration_Format(t, h3+m4, "%.7"+c, "3h 4m")
		assertDuration_Format(t, h3+m4, "%.128"+c, "3h 4m")

		assertDuration_Format(t, h3+ns8, "%"+c, "3h 8ns")
		assertDuration_Format(t, h3+ns8, "%.0"+c, "3h")
		assertDuration_Format(t, h3+ns8, "%.1"+c, "3h 8ns")
		assertDuration_Format(t, h3+ns8, "%.7"+c, "3h 8ns")
		assertDuration_Format(t, h3+ns8, "%.128"+c, "3h 8ns")
	}

	assertDuration_Format(t, h3, "%d", "%!d(go_pretty_print.Duration=3h0m0s)")

	assertDuration_Format(t, -h3, "%b", "-5937362789990400p-39")
	assertDuration_Format(t, -h3, "%.2b", "-5937362789990400p-39")
	assertDuration_Format(t, -h3, "%e", "-1.08e+04")
	assertDuration_Format(t, -h3, "%.2e", "-1.08e+04")
	assertDuration_Format(t, -h3, "%E", "-1.08E+04")
	assertDuration_Format(t, -h3, "%.2E", "-1.08E+04")
	assertDuration_Format(t, -h3, "%f", "-10800")
	assertDuration_Format(t, -h3, "%.2f", "-10800.00")
	assertDuration_Format(t, -h3, "%g", "-10800")
	assertDuration_Format(t, -h3, "%.2g", "-1.1e+04")
	assertDuration_Format(t, -h3, "%G", "-10800")
	assertDuration_Format(t, -h3, "%.2G", "-1.1E+04")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, -h3, "%"+c, "-3h")
		assertDuration_Format(t, -h3, "%.0"+c, "-3h")
		assertDuration_Format(t, -h3, "%.7"+c, "-3h")
		assertDuration_Format(t, -h3, "%.128"+c, "-3h")

		assertDuration_Format(t, -h3-m4, "%"+c, "-3h 4m")
		assertDuration_Format(t, -h3-m4, "%.0"+c, "-3h")
		assertDuration_Format(t, -h3-m4, "%.1"+c, "-3h 4m")
		assertDuration_Format(t, -h3-m4, "%.7"+c, "-3h 4m")
		assertDuration_Format(t, -h3-m4, "%.128"+c, "-3h 4m")

		assertDuration_Format(t, -h3-ns8, "%"+c, "-3h 8ns")
		assertDuration_Format(t, -h3-ns8, "%.0"+c, "-3h")
		assertDuration_Format(t, -h3-ns8, "%.1"+c, "-3h 8ns")
		assertDuration_Format(t, -h3-ns8, "%.7"+c, "-3h 8ns")
		assertDuration_Format(t, -h3-ns8, "%.128"+c, "-3h 8ns")
	}

	assertDuration_Format(t, -h3, "%d", "%!d(go_pretty_print.Duration=-3h0m0s)")
}

func TestDuration_Format_D(t *testing.T) {
	assertDuration_Format(t, d2, "%b", "5937362789990400p-35")
	assertDuration_Format(t, d2, "%.2b", "5937362789990400p-35")
	assertDuration_Format(t, d2, "%e", "1.728e+05")
	assertDuration_Format(t, d2, "%.2e", "1.73e+05")
	assertDuration_Format(t, d2, "%E", "1.728E+05")
	assertDuration_Format(t, d2, "%.2E", "1.73E+05")
	assertDuration_Format(t, d2, "%f", "172800")
	assertDuration_Format(t, d2, "%.2f", "172800.00")
	assertDuration_Format(t, d2, "%g", "172800")
	assertDuration_Format(t, d2, "%.2g", "1.7e+05")
	assertDuration_Format(t, d2, "%G", "172800")
	assertDuration_Format(t, d2, "%.2G", "1.7E+05")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, d2, "%"+c, "2d")
		assertDuration_Format(t, d2, "%.0"+c, "2d")
		assertDuration_Format(t, d2, "%.7"+c, "2d")
		assertDuration_Format(t, d2, "%.128"+c, "2d")

		assertDuration_Format(t, d2+h3, "%"+c, "2d 3h")
		assertDuration_Format(t, d2+h3, "%.0"+c, "2d")
		assertDuration_Format(t, d2+h3, "%.1"+c, "2d 3h")
		assertDuration_Format(t, d2+h3, "%.7"+c, "2d 3h")
		assertDuration_Format(t, d2+h3, "%.128"+c, "2d 3h")

		assertDuration_Format(t, d2+ns8, "%"+c, "2d 8ns")
		assertDuration_Format(t, d2+ns8, "%.0"+c, "2d")
		assertDuration_Format(t, d2+ns8, "%.1"+c, "2d 8ns")
		assertDuration_Format(t, d2+ns8, "%.7"+c, "2d 8ns")
		assertDuration_Format(t, d2+ns8, "%.128"+c, "2d 8ns")
	}

	assertDuration_Format(t, d2, "%d", "%!d(go_pretty_print.Duration=48h0m0s)")

	assertDuration_Format(t, -d2, "%b", "-5937362789990400p-35")
	assertDuration_Format(t, -d2, "%.2b", "-5937362789990400p-35")
	assertDuration_Format(t, -d2, "%e", "-1.728e+05")
	assertDuration_Format(t, -d2, "%.2e", "-1.73e+05")
	assertDuration_Format(t, -d2, "%E", "-1.728E+05")
	assertDuration_Format(t, -d2, "%.2E", "-1.73E+05")
	assertDuration_Format(t, -d2, "%f", "-172800")
	assertDuration_Format(t, -d2, "%.2f", "-172800.00")
	assertDuration_Format(t, -d2, "%g", "-172800")
	assertDuration_Format(t, -d2, "%.2g", "-1.7e+05")
	assertDuration_Format(t, -d2, "%G", "-172800")
	assertDuration_Format(t, -d2, "%.2G", "-1.7E+05")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, -d2, "%"+c, "-2d")
		assertDuration_Format(t, -d2, "%.0"+c, "-2d")
		assertDuration_Format(t, -d2, "%.7"+c, "-2d")
		assertDuration_Format(t, -d2, "%.128"+c, "-2d")

		assertDuration_Format(t, -d2-h3, "%"+c, "-2d 3h")
		assertDuration_Format(t, -d2-h3, "%.0"+c, "-2d")
		assertDuration_Format(t, -d2-h3, "%.1"+c, "-2d 3h")
		assertDuration_Format(t, -d2-h3, "%.7"+c, "-2d 3h")
		assertDuration_Format(t, -d2-h3, "%.128"+c, "-2d 3h")

		assertDuration_Format(t, -d2-ns8, "%"+c, "-2d 8ns")
		assertDuration_Format(t, -d2-ns8, "%.0"+c, "-2d")
		assertDuration_Format(t, -d2-ns8, "%.1"+c, "-2d 8ns")
		assertDuration_Format(t, -d2-ns8, "%.7"+c, "-2d 8ns")
		assertDuration_Format(t, -d2-ns8, "%.128"+c, "-2d 8ns")
	}

	assertDuration_Format(t, -d2, "%d", "%!d(go_pretty_print.Duration=-48h0m0s)")
}

func TestDuration_Format_W(t *testing.T) {
	assertDuration_Format(t, w1, "%b", "5195192441241600p-33")
	assertDuration_Format(t, w1, "%.2b", "5195192441241600p-33")
	assertDuration_Format(t, w1, "%e", "6.048e+05")
	assertDuration_Format(t, w1, "%.2e", "6.05e+05")
	assertDuration_Format(t, w1, "%E", "6.048E+05")
	assertDuration_Format(t, w1, "%.2E", "6.05E+05")
	assertDuration_Format(t, w1, "%f", "604800")
	assertDuration_Format(t, w1, "%.2f", "604800.00")
	assertDuration_Format(t, w1, "%g", "604800")
	assertDuration_Format(t, w1, "%.2g", "6e+05")
	assertDuration_Format(t, w1, "%G", "604800")
	assertDuration_Format(t, w1, "%.2G", "6E+05")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, w1, "%"+c, "1w")
		assertDuration_Format(t, w1, "%.0"+c, "1w")
		assertDuration_Format(t, w1, "%.7"+c, "1w")
		assertDuration_Format(t, w1, "%.128"+c, "1w")

		assertDuration_Format(t, w1+d2, "%"+c, "1w 2d")
		assertDuration_Format(t, w1+d2, "%.0"+c, "1w")
		assertDuration_Format(t, w1+d2, "%.1"+c, "1w 2d")
		assertDuration_Format(t, w1+d2, "%.7"+c, "1w 2d")
		assertDuration_Format(t, w1+d2, "%.128"+c, "1w 2d")

		assertDuration_Format(t, w1+ns8, "%"+c, "1w 8ns")
		assertDuration_Format(t, w1+ns8, "%.0"+c, "1w")
		assertDuration_Format(t, w1+ns8, "%.1"+c, "1w 8ns")
		assertDuration_Format(t, w1+ns8, "%.7"+c, "1w 8ns")
		assertDuration_Format(t, w1+ns8, "%.128"+c, "1w 8ns")
	}

	assertDuration_Format(t, w1, "%d", "%!d(go_pretty_print.Duration=168h0m0s)")

	assertDuration_Format(t, -w1, "%b", "-5195192441241600p-33")
	assertDuration_Format(t, -w1, "%.2b", "-5195192441241600p-33")
	assertDuration_Format(t, -w1, "%e", "-6.048e+05")
	assertDuration_Format(t, -w1, "%.2e", "-6.05e+05")
	assertDuration_Format(t, -w1, "%E", "-6.048E+05")
	assertDuration_Format(t, -w1, "%.2E", "-6.05E+05")
	assertDuration_Format(t, -w1, "%f", "-604800")
	assertDuration_Format(t, -w1, "%.2f", "-604800.00")
	assertDuration_Format(t, -w1, "%g", "-604800")
	assertDuration_Format(t, -w1, "%.2g", "-6e+05")
	assertDuration_Format(t, -w1, "%G", "-604800")
	assertDuration_Format(t, -w1, "%.2G", "-6E+05")

	for _, c := range [2]string{"s", "v"} {
		assertDuration_Format(t, -w1, "%"+c, "-1w")
		assertDuration_Format(t, -w1, "%.0"+c, "-1w")
		assertDuration_Format(t, -w1, "%.7"+c, "-1w")
		assertDuration_Format(t, -w1, "%.128"+c, "-1w")

		assertDuration_Format(t, -w1-d2, "%"+c, "-1w 2d")
		assertDuration_Format(t, -w1-d2, "%.0"+c, "-1w")
		assertDuration_Format(t, -w1-d2, "%.1"+c, "-1w 2d")
		assertDuration_Format(t, -w1-d2, "%.7"+c, "-1w 2d")
		assertDuration_Format(t, -w1-d2, "%.128"+c, "-1w 2d")

		assertDuration_Format(t, -w1-ns8, "%"+c, "-1w 8ns")
		assertDuration_Format(t, -w1-ns8, "%.0"+c, "-1w")
		assertDuration_Format(t, -w1-ns8, "%.1"+c, "-1w 8ns")
		assertDuration_Format(t, -w1-ns8, "%.7"+c, "-1w 8ns")
		assertDuration_Format(t, -w1-ns8, "%.128"+c, "-1w 8ns")
	}

	assertDuration_Format(t, -w1, "%d", "%!d(go_pretty_print.Duration=-168h0m0s)")
}

func assertDuration_Format(t *testing.T, d time.Duration, format, expected string) {
	t.Helper()

	assertCallResult(
		t,
		fmt.Sprintf("fmt.Sprintf(%#v, Duration(%v))", format, d),
		[]interface{}{expected},
		[]interface{}{fmt.Sprintf(format, Duration(d))},
	)
}

func assertCallResult(t *testing.T, call string, expected, actual []interface{}) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		buf := &bytes.Buffer{}

		fmt.Fprintf(buf, "Got unexpected result from %s:", call)

		for _, e := range expected {
			fmt.Fprint(buf, "\n- ")
			fprintHumanReadable(buf, e)
		}

		for _, a := range actual {
			fmt.Fprint(buf, "\n+ ")
			fprintHumanReadable(buf, a)
		}

		t.Error(buf.String())
	}
}

func fprintHumanReadable(w io.Writer, v interface{}) {
	switch x := v.(type) {
	case nil:
		fmt.Fprint(w, "nil")
	case bool:
		fmt.Fprintf(w, "bool(%v)", x)
	case int:
		fmt.Fprintf(w, "int(%d)", x)
	case int8:
		fmt.Fprintf(w, "int8(%d)", x)
	case int16:
		fmt.Fprintf(w, "int16(%d)", x)
	case int32:
		fmt.Fprintf(w, "int32(%d)", x)
	case int64:
		fmt.Fprintf(w, "int64(%d)", x)
	case uint:
		fmt.Fprintf(w, "uint(%d)", x)
	case uint8:
		fmt.Fprintf(w, "uint8(%d)", x)
	case uint16:
		fmt.Fprintf(w, "uint16(%d)", x)
	case uint32:
		fmt.Fprintf(w, "uint32(%d)", x)
	case uint64:
		fmt.Fprintf(w, "uint64(%d)", x)
	case float32:
		fmt.Fprintf(w, "float32(%s)", strconv.FormatFloat(float64(x), 'g', -1, 64))
	case float64:
		fmt.Fprintf(w, "float64(%s)", strconv.FormatFloat(x, 'g', -1, 64))
	case string:
		fmt.Fprintf(w, "string(%#v)", x)
	case []byte:
		fmt.Fprintf(w, "[]byte(%#v)", string(x))
	default:
		fmt.Fprintf(w, "%#v", v)
	}
}
