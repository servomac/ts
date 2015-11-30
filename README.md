# ts
Filter logs lines by specifying a timedelta.

> Show me the last 5 minutes of the syslog! Search me in this bunch of aggregated and not sorted nginx logs the entries from exactly a year ago between 5pm and 7pm! 

Using the reference time layout of Golang, *ts* parses a log searching for a timestamp with the layout specified format. With the parameters `-d`, `-h` and `-m` you specify the timedelta before of the actual instant where you want to start seeing logs.

## Examples

* Last hour and 5 minutes of syslog: `ts -l "Jan 02 15:04:05" -h 1 -m 5 /var/log/syslog`

