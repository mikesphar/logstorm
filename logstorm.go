package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func udp_send_measure(conn net.Conn, worker int) time.Duration {
	start := time.Now()
	fmt.Fprintf(conn, "test log message from %s: worker %d at %s", os.Args[0], worker, start)
	elapsed := time.Since(start)
	return elapsed
}

func udp_loop(target string, worker int, msg_rate int, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("udp", target)
	defer conn.Close()
	if err != nil {
		fmt.Println("Error opening connecton:  %v", err)
		return
	}

	for {
		_ = udp_send_measure(conn, worker)
		time.Sleep(time.Duration(1000000000 / msg_rate))
	}
}

func main() {

	port := "514"
	target := "127.0.0.1"

	var msg_rate int
	var workers int
	flag.IntVar(&msg_rate, "rate", 1, "Number of messages each worker will generate every second")
	flag.IntVar(&workers, "workers", 1, "Number of workers that will simultaneously generate log messages")

	flag.Parse()

	fmt.Printf("Spawning %d workers to each generate %d log messages every second\n", workers, msg_rate)

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go udp_loop(target+":"+port, i, msg_rate, &wg)
	}

	wg.Wait()
}
