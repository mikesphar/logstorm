package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"math/rand"
	"sync"
	"time"
)

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randChars(n int) string {
	randstr := make([]rune, n)
	for i := range randstr {
		randstr[i] = chars[rand.Intn(len(chars))]
	}
	return string(randstr)
}

type Flags struct {
	Msg_Rate  int
	Workers   int
	Source    string
	Message   string
	Size      int
	Msg_Count int
	Json_Out  bool
}

func pad_string(str string, padsize int) string {
	padcount := padsize - len(str)
	if padcount > 0 {
		return str + randChars(padcount)
	} else {
		return str
	}
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

	// feels like a waste to repeatedly get the size of Message so do it outside the loop
	var do_padding bool = false
	if flags.Size > len(flags.Message) {
		do_padding = true
	}

	for i := 0; i < flags.Msg_Count || infinite; i++ {
		raw_message.Timestamp = time.Now()

		if do_padding {
			raw_message.Message = pad_string(flags.Message, flags.Size)
		}

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

	rand.Seed(time.Now().UnixNano())

	port := "514"
	target := "127.0.0.1"

	flags := new(Flags)

	flag.IntVar(&flags.Msg_Rate, "rate", 1, "Number of messages each worker will generate every second")
	flag.IntVar(&flags.Workers, "workers", 1, "Number of workers that will simultaneously generate log messages")
	flag.StringVar(&flags.Source, "source", "logstorm", "String identifying the source of the log messages")
	flag.StringVar(&flags.Message, "message", "Test Message", "Message payload for every log message")
	flag.IntVar(&flags.Msg_Count, "count", -1, "Number of messages to generate per worker (-1 for unlimited)")
	flag.BoolVar(&flags.Json_Out, "json", false, "Format message as json")
	flag.IntVar(&flags.Size, "size", 0, "Minimum size of message string, pads message with random characters to this size if it's less than")

	flag.Parse()

	fmt.Printf("Spawning %d worker(s) to each generate %d log message(s) every second\n", flags.Workers, flags.Msg_Rate)
	fmt.Printf("Message source: %s\n", flags.Source)
	fmt.Printf("Message text: %s\n", flags.Message)
	if flags.Size > 0 {
		if flags.Size >= len(flags.Message) {
			fmt.Printf("Minimum size for Message field: %d\n", flags.Size)
		} else {
			fmt.Printf("Message field length already greater than minimum size (%d)\n", flags.Size)
			fmt.Printf("Message field size is: %d\n", len(flags.Message))
		}
	}

	fmt.Printf("Total messages generated per worker: ")
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
