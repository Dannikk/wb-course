package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func task3_1() {
	numbers := []int{2, 4, 6, 8, 10}

	wg := sync.WaitGroup{}
	summa := atomic.Int32{}

	for _, num := range numbers {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			summa.Add(int32(x*x))
			summa.Load()
		}(num)
	}
	wg.Wait()
	fmt.Println(summa.Load())
}

func task3_2() {
	numbers := []int{2, 4, 6, 8, 10}

	var summa int
	mx := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(numbers))

	for _, num := range numbers {
		go func(x int) {
			defer wg.Done()
			mx.Lock()
			summa += x * x
			mx.Unlock()
		}(num)
	}

	wg.Wait()
	fmt.Println(summa)
}

func task3_3() {
	numbers := []int{2, 4, 6, 8, 10}

	var summa int
	ch := make(chan int, len(numbers))

	var wg sync.WaitGroup
	wg.Add(len(numbers))

	for _, num := range numbers {
		go func(x int) {
			defer wg.Done()
			ch <- x * x
		}(num)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for x := range ch {
		summa += x
	}
	fmt.Println(summa)
}
