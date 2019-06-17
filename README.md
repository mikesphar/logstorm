# logstorm
A very bad go program for generating way too many log messages and blasts them at the local host syslog port (UDP 514)

Usage:

*  -rate int
   -  Number of messages each worker will generate every second (default 1)
*  -message string
   -  Message payload for every log message (default "Test Message")
*  -rate int
   -  Number of messages each worker will generate every second (default 1)
*  -source string
   -  String identifying the source of the log messages (default "logstorm")
*  -workers int
   -  Number of workers that will simultaneously generate log messages (default 1)

Will generate (potentially very large amounts of) messages that look like this:

```
Jun 17 14:08:11 127.0.0.1  {"Source":"logstorm","Worker":0,"Message":"Test Message","Timestamp":"2019-06-17T14:08:11.676208505Z"}
```
