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

type Flags struct {
	Msg_Rate  int
	Workers   int
	Source    string
	Message   string
	Msg_Count int
	Json_Out  bool
}

func send_logs(target string, worker int, flags *Flags, wg *sync.WaitGroup) {
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

	if flags.Msg_Count < 0 {
		infinite = true
	}

	raw_message := LogMessage{
		Source:  flags.Source,
		Message: flags.Message,
		Worker:  worker,
	}

	for i := 0; i < flags.Msg_Count || infinite; i++ {
		raw_message.Timestamp = time.Now()

		if flags.Json_Out {
			json_message, _ := json.Marshal(raw_message)
			fmt.Fprintf(sysLog, string(json_message))
		} else {
			fmt.Fprintf(sysLog, "From %s: worker %d at %s - %s", raw_message.Source, raw_message.Worker, raw_message.Timestamp, raw_message.Message)
		}

		time.Sleep(time.Duration(1000000000 / flags.Msg_Rate))
	}
	return
}

func main() {

	port := "514"
	target := "127.0.0.1"

	flags := new(Flags)

	flag.IntVar(&flags.Msg_Rate, "rate", 1, "Number of messages each worker will generate every second")
	flag.IntVar(&flags.Workers, "workers", 1, "Number of workers that will simultaneously generate log messages")
	flag.StringVar(&flags.Source, "source", "logstorm", "String identifying the source of the log messages")
	flag.StringVar(&flags.Message, "message", "Test Message", "Message payload for every log message")
	flag.IntVar(&flags.Msg_Count, "count", -1, "Number of messages to generate per worker (-1 for unlimited)")
	flag.BoolVar(&flags.Json_Out, "json", false, "Format message as json")

	flag.Parse()

	fmt.Printf("Spawning %d workers to each generate %d log messages every second\n", flags.Workers, flags.Msg_Rate)
	fmt.Printf("Message source: %s\n", flags.Source)
	fmt.Printf("Message text: %s\n", flags.Message)
	fmt.Printf("Total per worker: ")

	if flags.Msg_Count == 0 {
		fmt.Printf("Unlimited\n")
	} else {
		fmt.Printf("%d\n", flags.Msg_Count)
	}

	if flags.Json_Out {
		fmt.Printf("Format: json \n")
	} else {
		fmt.Printf("Format: text \n")
	}
	var wg sync.WaitGroup
	wg.Add(flags.Workers)

	for i := 0; i < flags.Workers; i++ {
		go send_logs(target+":"+port, i, flags, &wg)
	}

	wg.Wait()
}
