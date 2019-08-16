# logstorm
A very bad go program for generating way too many log messages and blasting them at the local syslog service.

Usage:


*  -count int
   -  Number of messages to generate per worker (-1 for unlimited) (default -1)
*  -json
   -  Format message as json
*  -message string
   -  Message payload for every log message (default "Test Message")
*  -rate int
   -  Number of messages each worker will generate every second (default 1)
*  -size int
   -  Minimum size of message string, pads message to this size if it's less than
*  -source string
   -  String identifying the source of the log messages (default "logstorm")
*  -workers int
   -  Number of workers that will simultaneously generate log messages (default 1)

Will generate (potentially very large amounts of) messages that look like this:

JSON:
```
Jun 21 22:14:39 somehost.somewhere.foo logstorm[11790]: {"Source":"logstorm","Worker":1,"Message":"Test Message","Timestamp":"2019-06-21T22:14:39.880695653Z"}
```

Plain text:
```
Jun 21 22:11:41 somehost.somewhere.foo logstorm[21343]: From logstorm: worker 0 at 2019-06-21 22:11:41.24345764 +0000 UTC m=+3.001405705 - Test Message
```
