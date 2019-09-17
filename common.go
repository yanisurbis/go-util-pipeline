package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"strconv"
	"sync/atomic"
	"time"
	"sort"
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
	data := <-in
	data1 := data.(string)
	out <- DataSignerCrc32(data1) + "~" + DataSignerCrc32(data1)
}

func MultiHash(in, out chan interface{}) {
	data := <-in
	data1 := data.(string)

	var res string

	for i := 0; i < 6; i++ {
		res += DataSignerCrc32(strconv.Itoa(i) + data1)
	}

	out <- res
}

func CombineResults(in, out chan interface{}) {
	data := make([]string, 0, 6)

	for i := 0; i < 6; i++ {
		d := <-in
		d1 := d.(string)
		data = append(data, d1)
	}

	sort.Strings(data)

	res := data[0]

	for _, s := range data {
		res += "_" + s
	}

	out <- res
}
