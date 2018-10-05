package main

import (
	"fmt"
	"sync"
	"time"
)

func tickAndEnd(stop chan (struct{}), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		fmt.Println("Started")

		ticker := time.NewTicker(time.Second).C
		for {
			select {
			case <-ticker:
				fmt.Print(".")
			case <-stop:
				fmt.Println("Ended")
				wg.Done()
				return
			}
		}
	}()
}

func main() {
	stop := make(chan (struct{}))
	wg := &sync.WaitGroup{}

	tickAndEnd(stop, wg)
	tickAndEnd(stop, wg)
	tickAndEnd(stop, wg)

	time.Sleep(5 * time.Second)
	close(stop)

	wg.Wait()
}
