package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "regexp"
  "time"
)



// given a string and a time layout, search for a timestamp in the
// line and return true if it's greather than the initial time
func between(line string, layout string, re regexp.Regexp, initial time.Time) bool {
    timestamp := re.FindString(line)
    if len(timestamp) > 0 {
        t, err := time.Parse(layout, timestamp)
        if err != nil {
            panic(err)
        }

        return initial.Before(t) || initial == t
    }

    return false
}

func main() {
    re := regexp.MustCompile(`\d{2}/\w{3}/\d{4}:\d{2}:\d{2}:\d{2} (\+|-)\d{4}`)
    layout := "02/Jan/2006:15:04:05 -0700"

    minutesPtr := flag.Int("m", 0, "minutes")
    hoursPtr := flag.Int("h", 0, "hours")
    daysPtr := flag.Int("d", 0, "days")
    flag.Parse()

    initial := time.Now().
        AddDate(0, 0, -*daysPtr).
        Add(-time.Duration(*hoursPtr)*time.Hour).
        Add(-time.Duration(*minutesPtr)*time.Minute)

    f, err := os.Open("prova.log")
    if err != nil {
        panic(err)
    }

    scanner := bufio.NewScanner(f)

    for scanner.Scan() {
        line := scanner.Text()
        if between(line, layout, *re, initial) {
            fmt.Println(line)
        }
    }
}
