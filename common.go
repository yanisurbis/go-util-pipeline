package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"strconv"
	"sync/atomic"
	"time"
	"sort"
	"sync"
	"strings"
)

type job func(in, out chan interface{})

const (
	MaxInputDataLen = 100
)

var (
	dataSignerOverheat uint32 = 0
	DataSignerSalt            = ""
)

var OverheatLock = func() {
	for {
		if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 0, 1); !swapped {
			fmt.Println("OverheatLock happend")
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

var OverheatUnlock = func() {
	for {
		if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 1, 0); !swapped {
			fmt.Println("OverheatUnlock happend")
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

var DataSignerMd5 = func(data string) string {
	OverheatLock()
	defer OverheatUnlock()
	data += DataSignerSalt
	dataHash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	time.Sleep(10 * time.Millisecond)
	return dataHash
}

var DataSignerCrc32 = func(data string) string {
	data += DataSignerSalt
	crcH := crc32.ChecksumIEEE([]byte(data))
	dataHash := strconv.FormatUint(uint64(crcH), 10)
	time.Sleep(time.Second)
	return dataHash
}

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
		//fmt.Println(res)
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
