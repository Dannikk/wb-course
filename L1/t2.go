package main

import (
	"fmt"
	"sync"
)


func task2_1() {
	numbers := []int{2, 4, 6, 8, 10}

	wg := sync.WaitGroup{}

	for _, num := range numbers {
		num := num
		wg.Add(1)
		go func(){
			defer wg.Done()
			fmt.Println(num*num)
		}()
	}
	wg.Wait()
}

func task2_2() {
	numbers := []int{2, 4, 6, 8, 10}
	wg := sync.WaitGroup{}

	for _, num := range numbers {
		wg.Add(1)
		go func(x int){
			defer wg.Done()
			fmt.Println(x*x)
		}(num)
	}
	wg.Wait()
}