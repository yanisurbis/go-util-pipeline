package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func SingleHash(in, out chan interface{}) {
	for data := range in {
		data1 := (data).(int)
		dataS := strconv.Itoa(data1)

		crc32_ := make(chan string)
		crc32md5_ := make(chan string)

		go func() {
			crc32_ <- DataSignerCrc32(dataS)
		}()

		go func() {
			crc32md5_ <- DataSignerCrc32(DataSignerMd5(dataS))
		}()

		res := <-crc32_ + "~" + <-crc32md5_
		out <- res
	}
}

func MultiHash(in, out chan interface{}) {
	for data := range in {
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

		wg.Wait()

		mu.Lock()
		for _, s := range resStrings {
			res += s
		}
		mu.Unlock()

		out <- res
	}
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
