package timeregex

import (
    "strings"
)


// This section of code, (constants, isDigit, nextStdChunk, startsWithLowerCase)
// are extracted from https://github.com/golang/go/blob/master/src/time/format.go
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
const (
    _                        = iota
    stdLongMonth             = iota + stdNeedDate  // "January"
    stdMonth                                       // "Jan"
    stdNumMonth                                    // "1"
    stdZeroMonth                                   // "01"
    stdLongWeekDay                                 // "Monday"
    stdWeekDay                                     // "Mon"
    stdDay                                         // "2"
    stdUnderDay                                    // "_2"
    stdZeroDay                                     // "02"
    stdHour                  = iota + stdNeedClock // "15"
    stdHour12                                      // "3"
    stdZeroHour12                                  // "03"
    stdMinute                                      // "4"
    stdZeroMinute                                  // "04"
    stdSecond                                      // "5"
    stdZeroSecond                                  // "05"
    stdLongYear              = iota + stdNeedDate  // "2006"
    stdYear                                        // "06"
    stdPM                    = iota + stdNeedClock // "PM"
    stdpm                                          // "pm"
    stdTZ                    = iota                // "MST"
    stdISO8601TZ                                   // "Z0700"  // prints Z for UTC
    stdISO8601SecondsTZ                            // "Z070000"
    stdISO8601ColonTZ                              // "Z07:00" // prints Z for UTC
    stdISO8601ColonSecondsTZ                       // "Z07:00:00"
    stdNumTZ                                       // "-0700"  // always numeric
    stdNumSecondsTz                                // "-070000"
    stdNumShortTZ                                  // "-07"    // always numeric
    stdNumColonTZ                                  // "-07:00" // always numeric
    stdNumColonSecondsTZ                           // "-07:00:00"
    stdFracSecond0                                 // ".0", ".00", ... , trailing zeros included
    stdFracSecond9                                 // ".9", ".99", ..., trailing zeros omitted

    stdNeedDate  = 1 << 8             // need month, day, year
    stdNeedClock = 2 << 8             // need hour, minute, second
    stdArgShift  = 16                 // extra argument in high bits, above low stdArgShift
    stdMask      = 1<<stdArgShift - 1 // mask out argument
)

// std0x records the std values for "01", "02", ..., "06".
var std0x = [...]int{stdZeroMonth, stdZeroDay, stdZeroHour12, stdZeroMinute, stdZeroSecond, stdYear}

// startsWithLowerCase reports whether the string has a lower-case letter at the beginning.
// Its purpose is to prevent matching strings like "Month" when looking for "Mon".
func startsWithLowerCase(str string) bool {
    if len(str) == 0 {
        return false
    }
    c := str[0]
    return 'a' <= c && c <= 'z'
}

// nextStdChunk finds the first occurrence of a std string in
// layout and returns the text before, the std string, and the text after.
func nextStdChunk(layout string) (prefix string, std int, suffix string) {
	for i := 0; i < len(layout); i++ {
		switch c := int(layout[i]); c {
		case 'J': // January, Jan
			if len(layout) >= i+3 && layout[i:i+3] == "Jan" {
				if len(layout) >= i+7 && layout[i:i+7] == "January" {
					return layout[0:i], stdLongMonth, layout[i+7:]
				}
				if !startsWithLowerCase(layout[i+3:]) {
					return layout[0:i], stdMonth, layout[i+3:]
				}
			}

		case 'M': // Monday, Mon, MST
			if len(layout) >= i+3 {
				if layout[i:i+3] == "Mon" {
					if len(layout) >= i+6 && layout[i:i+6] == "Monday" {
						return layout[0:i], stdLongWeekDay, layout[i+6:]
					}
					if !startsWithLowerCase(layout[i+3:]) {
						return layout[0:i], stdWeekDay, layout[i+3:]
					}
				}
				if layout[i:i+3] == "MST" {
					return layout[0:i], stdTZ, layout[i+3:]
				}
			}

		case '0': // 01, 02, 03, 04, 05, 06
			if len(layout) >= i+2 && '1' <= layout[i+1] && layout[i+1] <= '6' {
				return layout[0:i], std0x[layout[i+1]-'1'], layout[i+2:]
			}

		case '1': // 15, 1
			if len(layout) >= i+2 && layout[i+1] == '5' {
				return layout[0:i], stdHour, layout[i+2:]
			}
			return layout[0:i], stdNumMonth, layout[i+1:]

		case '2': // 2006, 2
			if len(layout) >= i+4 && layout[i:i+4] == "2006" {
				return layout[0:i], stdLongYear, layout[i+4:]
			}
			return layout[0:i], stdDay, layout[i+1:]

		case '_': // _2, _2006
			if len(layout) >= i+2 && layout[i+1] == '2' {
				//_2006 is really a literal _, followed by stdLongYear
				if len(layout) >= i+5 && layout[i+1:i+5] == "2006" {
					return layout[0 : i+1], stdLongYear, layout[i+5:]
				}
				return layout[0:i], stdUnderDay, layout[i+2:]
			}

		case '3':
			return layout[0:i], stdHour12, layout[i+1:]

		case '4':
			return layout[0:i], stdMinute, layout[i+1:]

		case '5':
			return layout[0:i], stdSecond, layout[i+1:]

		case 'P': // PM
			if len(layout) >= i+2 && layout[i+1] == 'M' {
				return layout[0:i], stdPM, layout[i+2:]
			}

		case 'p': // pm
			if len(layout) >= i+2 && layout[i+1] == 'm' {
				return layout[0:i], stdpm, layout[i+2:]
			}

		case '-': // -070000, -07:00:00, -0700, -07:00, -07
			if len(layout) >= i+7 && layout[i:i+7] == "-070000" {
				return layout[0:i], stdNumSecondsTz, layout[i+7:]
			}
			if len(layout) >= i+9 && layout[i:i+9] == "-07:00:00" {
				return layout[0:i], stdNumColonSecondsTZ, layout[i+9:]
			}
			if len(layout) >= i+5 && layout[i:i+5] == "-0700" {
				return layout[0:i], stdNumTZ, layout[i+5:]
			}
			if len(layout) >= i+6 && layout[i:i+6] == "-07:00" {
				return layout[0:i], stdNumColonTZ, layout[i+6:]
			}
			if len(layout) >= i+3 && layout[i:i+3] == "-07" {
				return layout[0:i], stdNumShortTZ, layout[i+3:]
			}

		case 'Z': // Z070000, Z07:00:00, Z0700, Z07:00,
			if len(layout) >= i+7 && layout[i:i+7] == "Z070000" {
				return layout[0:i], stdISO8601SecondsTZ, layout[i+7:]
			}
			if len(layout) >= i+9 && layout[i:i+9] == "Z07:00:00" {
				return layout[0:i], stdISO8601ColonSecondsTZ, layout[i+9:]
			}
			if len(layout) >= i+5 && layout[i:i+5] == "Z0700" {
				return layout[0:i], stdISO8601TZ, layout[i+5:]
			}
			if len(layout) >= i+6 && layout[i:i+6] == "Z07:00" {
				return layout[0:i], stdISO8601ColonTZ, layout[i+6:]
			}

		case '.': // .000 or .999 - repeated digits for fractional seconds.
			if i+1 < len(layout) && (layout[i+1] == '0' || layout[i+1] == '9') {
				ch := layout[i+1]
				j := i + 1
				for j < len(layout) && layout[j] == ch {
					j++
				}
				// String of digits must end here - only fractional second is all digits.
				if !isDigit(layout, j) {
					std := stdFracSecond0
					if layout[i+1] == '9' {
						std = stdFracSecond9
					}
					std |= (j - (i + 1)) << stdArgShift
					return layout[0:i], std, layout[j:]
				}
			}
		}
	}
	return layout, 0, ""
}

// isDigit reports whether s[i] is in range and is a decimal digit.
func isDigit(s string, i int) bool {
    if len(s) <= i {
        return false
    }
    c := s[i]
    return '0' <= c && c <= '9'
}

// End of Go Authors code (Copyright 2010 The Go Authors).


func GenerateRegex(layout string) string {
    re_equiv := map[int]string {
        stdYear:        "\\d{2}",
        stdLongYear:    "\\d{4}",
        stdMonth:       "(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)",
        stdLongMonth:   "(January|February|March|April|May|June|July|August|September|October|November|December)",
        stdNumMonth:    "([1-9]|1[0-2])",
        stdZeroMonth:   "(0[1-9]|1[0-2])",
        stdWeekDay:     "(Sun|Mon|Tue|Wed|Thu|Fri|Sat)",
        stdLongWeekDay: "(Sunday|Monday|Tuesday|Wednesday|Thursday|Friday|Saturday)",
        stdDay:         "([1-9]|[12]\\d|3[0-1])",
        stdUnderDay:    "(_[1-9]|[12]\\d|3[0-1])",
        stdZeroDay:     "(0[1-9]|[12]\\d|3[0-1])",
        stdHour:        "([01][0-9]|2[0-3])",
        stdHour12:      "([1-9]|1[0-2])",
        stdZeroHour12:  "(0[1-9]|1[0-2])",
        stdMinute:      "[1-5]?\\d",
        stdZeroMinute:  "[0-5]\\d",
        stdSecond:      "[1-5]?\\d",
        stdZeroSecond:  "[0-5]\\d",
        stdPM:          "[AP]M",
        stdpm:          "[ap]m",
        stdISO8601TZ:   "",
        stdISO8601ColonTZ:          "",
        stdISO8601SecondsTZ:        "",
        stdISO8601ColonSecondsTZ:   "",
        stdNumTZ:                   "(\\+|-)\\d{4}",
        stdNumShortTZ:              "(\\+|-)\\d{2}",
        stdNumColonTZ:              "(\\+|-)\\d{2}:\\d{2}",
        stdNumSecondsTz:            "(\\+|-)\\d{6}",
        stdNumColonSecondsTZ:       "(\\+|-)\\d{2}:\\d{2}:\\d{2}",
        stdTZ:                      "(MST|GSM|UTC)",
        stdFracSecond0:             "",
        stdFracSecond9:             "",
    }

    var regex []string
    for {
        prefix, std, suffix := nextStdChunk(layout)
        if std == 0 {
            break
        }
        regex = append(regex, prefix, re_equiv[std & stdMask])
        layout = suffix
    }

    return strings.Join(regex, "")
}
