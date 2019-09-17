package main

import (
	"fmt"
	"sync"
)

// сюда писать код
func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}

	l := len(jobs) + 1
	chans := make([]chan interface{}, l)
	fmt.Println(l, "channels created")

	for i := 0; i < l; i++ {
		chans[i] = make(chan interface{}, 3)
	}

	for i, job := range jobs {
		wg.Add(1)
		go func(in, out chan interface{}) {
			defer wg.Done()
			job(in, out)
			fmt.Println("done", i)
			close(in)
			close(out)
		}(chans[i], chans[i+1])
	}

	//chan1 := make(chan interface{}, 3)
	//chan2 := make(chan interface{}, 3)
	//
	//job1 := jobs[0]
	//job2 := jobs[1]
	//
	//wg.Add(1)
	//go func(in, out chan interface{}) {
	//	defer wg.Done()
	//	job1(in, out)
	//	fmt.Println("done 1")
	//	close(out)
	//}(chan1, chan2)
	//
	//wg.Add(1)
	//go func(in, out chan interface{}) {
	//	defer wg.Done()
	//	job2(in, out)
	//	fmt.Println("done 2")
	//	close(out)
	//}(chan2, chan1)

	wg.Wait()
}
