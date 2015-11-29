package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "regexp"
  "time"
  "github.com/servomac/ts/timeregex"
)

func main() {
    layoutPtr := flag.String("l", "", "timestamp layout")
    minutesPtr := flag.Int("m", 0, "minutes")
    hoursPtr := flag.Int("h", 0, "hours")
    daysPtr := flag.Int("d", 0, "days")

    flag.Parse()
    if flag.NArg() != 1 || flag.NFlag() < 2 || len(*layoutPtr) == 0 {
        flag.Usage()
        os.Exit(1)
    }

    layout := *layoutPtr
    re_str := timeregex.GenerateRegex(layout)
    rePtr := regexp.MustCompile(re_str)
    re := *rePtr

    initial_formated := time.Now().AddDate(0, 0, -*daysPtr).
            Add(-time.Duration(*hoursPtr)*time.Hour).
            Add(-time.Duration(*minutesPtr)*time.Minute)//.
            //Format(layout)
    //initial := time.Parse(layout, initial_formated)
    fmt.Println("Filter starts at ", initial)

    filename := flag.Arg(0)
    f, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)

    for scanner.Scan() {
        line := scanner.Text()
        timestamp := re.FindString(line)
        if len(timestamp) > 0 {

            t, err := time.Parse(layout, timestamp)
            if err != nil {
                panic(err)
            }

            if t.After(initial) || t.Equal(initial) {
                fmt.Println("MATCH ", line)
            }
        }
    }

}
