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
    if flag.NFlag() > 3 || len(*layoutPtr) == 0 {
        flag.Usage()
        os.Exit(1)
    }

    layout := *layoutPtr
    reStr := timeregex.GenerateRegex(layout)
    fmt.Println(layout, " -> ", reStr)
    re := *regexp.MustCompile(reStr)

    starting_time := time.Now().AddDate(0, 0, -*daysPtr).
        Add(-time.Duration(*hoursPtr)*time.Hour).
        Add(-time.Duration(*minutesPtr)*time.Minute)
    // format and parse again the desired starting log time,
    // to easily compare with formatted timestamps
    starting_time, err := time.Parse(layout, starting_time.Format(layout))
    if err != nil {
        panic(err)
    }
    fmt.Println("Filter from: ", starting_time)

    var scanner *bufio.Scanner;
    if flag.NArg() == 1 {
        filename := flag.Arg(0)
        f, err := os.Open(filename)
        if err != nil {
            panic(err)
        }
        defer f.Close()

        scanner = bufio.NewScanner(f)
    } else if flag.NArg() == 0 {
        scanner = bufio.NewScanner(os.Stdin)
    }

    for scanner.Scan() {
        line := scanner.Text()
        timestamp := re.FindString(line)
        if len(timestamp) > 0 {
            t, err := time.Parse(layout, timestamp)
            if err != nil {
                panic(err)
            }

            if t.After(starting_time) || t.Equal(starting_time) {
                fmt.Println(line)
            }
        }
    }

}
