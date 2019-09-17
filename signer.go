package main

import (
	"sync"
)

// сюда писать код
func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}

	l := len(jobs)

	chans := make([]chan interface{}, l)

	for i := 0; i < l; i++ {
		chans[i] = make(chan interface{}, 1)
	}

	for i, j := range jobs {
		wg.Add(1)
		go func(in, out chan interface{}, job job) {
			defer wg.Done()
			job(in, out)
			close(out)
		}(chans[i], chans[(i+1)%l], j)
	}

	wg.Wait()
}
