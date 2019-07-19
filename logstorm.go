package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"sync"
	"time"
)

func send_logs(target string, worker int, msg_rate int, msg_count int, source string, message string, json_out bool, wg *sync.WaitGroup) {
	defer wg.Done()

	sysLog, err := syslog.Dial("udp", target,
		syslog.LOG_WARNING|syslog.LOG_DAEMON, "logstorm")
	if err != nil {
		log.Fatal(err)
	}

	type LogMessage struct {
		Source    string
		Worker    int
		Message   string
		Timestamp time.Time
	}

	var infinite bool = false

	if msg_count < 0 {
		infinite = true
	}

	for i := 0; i < msg_count || infinite; i++ {
		raw_message := LogMessage{
			Source:    source,
			Worker:    worker,
			Message:   message,
			Timestamp: time.Now(),
		}

		if json_out {
			json_message, _ := json.Marshal(raw_message)
			fmt.Fprintf(sysLog, string(json_message))
		} else {
			fmt.Fprintf(sysLog, "From %s: worker %d at %s - %s", raw_message.Source, raw_message.Worker, raw_message.Timestamp, raw_message.Message)
		}

		time.Sleep(time.Duration(1000000000 / msg_rate))
	}
	return
}

func main() {

	port := "514"
	target := "127.0.0.1"

	var msg_rate, workers, msg_count int
	var source, message string
	var json_out bool

	flag.IntVar(&msg_rate, "rate", 1, "Number of messages each worker will generate every second")
	flag.IntVar(&workers, "workers", 1, "Number of workers that will simultaneously generate log messages")
	flag.StringVar(&source, "source", "logstorm", "String identifying the source of the log messages")
	flag.StringVar(&message, "message", "Test Message", "Message payload for every log message")
	flag.IntVar(&msg_count, "count", -1, "Number of messages to generate per worker (-1 for unlimited)")
	flag.BoolVar(&json_out, "json", false, "Format message as json")

	flag.Parse()

	fmt.Printf("Spawning %d workers to each generate %d log messages every second\n", workers, msg_rate)
	fmt.Printf("Message source: %s\n", source)
	fmt.Printf("Message text: %s\n", message)
	fmt.Printf("Total per worker: ")

	if msg_count == 0 {
	  fmt.Printf("Unlimited\n")
	} else {
	  fmt.Printf("%d\n", msg_count)
	}

	if json_out {
		fmt.Printf("Format: json \n")
	} else {
		fmt.Printf("Format: text \n")
	}
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go send_logs(target+":"+port, i, msg_rate, msg_count, source, message, json_out, &wg)
	}

	wg.Wait()
}
