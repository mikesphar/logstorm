package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

func send_logs(target string, worker int, msg_rate int, source string, message string, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("udp", target)
	defer conn.Close()
	if err != nil {
		fmt.Println("Error opening connecton:  %v", err)
		return
	}

	type LogMessage struct {
		Source    string
		Worker    int
		Message   string
		Timestamp time.Time
	}

	for {
		raw_message := LogMessage{
			Source:    source,
			Worker:    worker,
			Message:   message,
			Timestamp: time.Now(),
		}
		json_message, _ := json.Marshal(raw_message)
		fmt.Fprintf(conn, string(json_message))
		time.Sleep(time.Duration(1000000000 / msg_rate))
	}
}

func main() {

	port := "514"
	target := "127.0.0.1"

	var msg_rate, workers int
	var source, message string

	flag.IntVar(&msg_rate, "rate", 1, "Number of messages each worker will generate every second")
	flag.IntVar(&workers, "workers", 1, "Number of workers that will simultaneously generate log messages")
	flag.StringVar(&source, "source", "logstorm", "String identifying the source of the log messages")
	flag.StringVar(&message, "message", "Test Message", "Message payload for every log message")

	flag.Parse()

	fmt.Printf("Spawning %d workers to each generate %d log messages every second\n", workers, msg_rate)
	fmt.Printf("Message source: %s\n", source)
	fmt.Printf("Message text: %s\n", message)

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go send_logs(target+":"+port, i, msg_rate, source, message, &wg)
	}

	wg.Wait()
}
