package main

import (
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	quotaCh := make(chan struct{}, 1)

	for data := range in {
		wg.Add(1)
		data1 := (data).(int)
		dataS := strconv.Itoa(data1)

		crc32_ := make(chan string)
		crc32md5_ := make(chan string)

		go func(dataS string) {
			crc32_ <- DataSignerCrc32(dataS)
		}(dataS)

		go func(dataS string, quotaCh chan struct{}) {
			quotaCh <- struct{}{}
			md5_ := DataSignerMd5(dataS)
			<-quotaCh
			runtime.Gosched()
			crc32md5_ <- DataSignerCrc32(md5_)
		}(dataS, quotaCh)

		go func() {
			defer wg.Done()
			res := <-crc32_ + "~" + <-crc32md5_
			out <- res
		}()
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg_ := &sync.WaitGroup{}
	for data := range in {
		wg_.Add(1)
		mu := &sync.Mutex{}
		wg := &sync.WaitGroup{}

		data1 := data.(string)

		var res string
		var resStrings [6]string

		for i := 0; i < 6; i++ {
			wg.Add(1)
			go func(data string, i int, mu *sync.Mutex, resStrings *[6]string) {
				defer wg.Done()
				res := DataSignerCrc32(strconv.Itoa(i) + data1)
				mu.Lock()
				resStrings[i] = res
				mu.Unlock()
			}(data1, i, mu, &resStrings)
		}

		go func() {
			defer wg_.Done()
			wg.Wait()

			mu.Lock()
			for _, s := range resStrings {
				res += s
			}
			mu.Unlock()

			out <- res
		}()
	}
	wg_.Wait()
}

func CombineResults(in, out chan interface{}) {
	var data []string

	for d := range in {
		data = append(data, d.(string))
	}

	sort.Strings(data)
	out <- strings.Join(data, "_")
}

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
