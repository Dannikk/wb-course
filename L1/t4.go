package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func worker(id int, dataCh <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range dataCh {
		time.Sleep(1 * time.Second)
		fmt.Printf("Worker %d received data: %d\n", id, data)

	}
	fmt.Printf("Worker %d stopped\n", id)
}

func task_4_1() {

	numWorkers := 5
	if args := flag.Args(); len(args) > 1 {
		numWorkers, _ = strconv.Atoi(args[1])
	}

	dataCh := make(chan int)
	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, dataCh, &wg)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

loop1:
	for i := 1; ; i++ {
		select {
		case sig := <-signals:
			fmt.Println("\nCtrl+C pressed: ", sig)
			close(dataCh)
			break loop1
		default:
			dataCh <- i
			fmt.Printf("Sent data: %d to the channel\n", i)
			time.Sleep(80 * time.Millisecond)
		}
	}

	wg.Wait()
}
