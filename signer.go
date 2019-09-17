package main

import (
	"fmt"
	"sync"
)

// сюда писать код
func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}

	l := len(jobs)

	chans := make([]chan interface{}, l)
	fmt.Println(l, "channels created")

	for i := 0; i < l; i++ {
		chans[i] = make(chan interface{}, 1)
	}

	for i, j := range jobs {
		wg.Add(1)
		go func(in, out chan interface{}, job job) {
			defer wg.Done()
			job(in, out)
			fmt.Println("done", i)
			close(out)
		}(chans[i], chans[(i+1)%l], j)
	}

	wg.Wait()
}
