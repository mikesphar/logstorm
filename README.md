# logstorm
A very bad go program for generating way too many log messages and blasts them at the local host syslog port (UDP 514)

Usage:

*  -rate int
  *  Number of messages each worker will generate every second (default 1)

*  -workers int
  * Number of workers that will simultaneously generate log messages (default 1)

Will generate (potentially very large amounts of) messages that look like this:

```
test log message from ./logstorm worker 2 at 2019-06-14 15:56:31.502922819 +0000 UTC m=+3.001445236
```
