package main

import (
	//	"context"
	"fmt"
	//	"runtime"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	//	"time"
)

var (
	md5Lock uint32 = 0
	dummy          = make(chan interface{})
)

func ExecutePipeline(jobs ...job) {
	rawDataCh := make(chan interface{}, MaxInputDataLen)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go jobWrapper(jobs[0], dummy, rawDataCh, wg)

	jobsSlice := jobs[1:]
	jobsSlCount := len(jobsSlice)

	for num := range rawDataCh {
		fmt.Println("ExecPipe: num=", num)

		var chans []chan interface{}
		for i := 0; i < jobsSlCount; i++ {
			chans = append(chans, make(chan interface{}, MaxInputDataLen))
		}
		chans = append(chans, dummy)

		chans[0] <- num

		for i, job := range jobsSlice {
			wg.Add(1)

			go jobWrapper(job, chans[i], chans[i+1], wg)
		}

		close(chans[0])
	}

	wg.Wait()

	fmt.Println("End ExecutePipeline")
}

func jobWrapper(jb job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	if out != dummy {
		defer close(out)
	}
	defer fmt.Println("job over")

	jb(in, out)
}

func SingleHash(in, out chan interface{}) {
LOOP:
	for {
		select {
		case data := <-in:

			if data == nil {
				out <- data
				fmt.Println("Single hash closed")
				break LOOP
			}

			dataStr := strconv.Itoa(data.(int))

			wg := &sync.WaitGroup{}
			wg.Add(2)

			dataSl := make([]string, 2)

			md5 := func(str string) string {
				for {
					if unlock := atomic.CompareAndSwapUint32(&md5Lock, 0, 1); unlock {
						defer atomic.StoreUint32(&md5Lock, 0)
						return DataSignerMd5(str)
					}
				}
			}

			go Crc32Worker(wg, md5(dataStr), dataSl[1:2])
			go Crc32Worker(wg, dataStr, dataSl[0:1])

			wg.Wait()

			out <- dataSl[0] + "~" + dataSl[1]

		}
	}
}

func MultiHash(in, out chan interface{}) {
LOOP:
	for data := range in {

		if data == nil {
			out <- data
			fmt.Println("Multi hash closed")
			break LOOP
		}

		wg := &sync.WaitGroup{}
		wg.Add(6)

		dataSl := make([]string, 6)

		for i := 0; i < 6; i++ {
			go Crc32Worker(wg, strconv.Itoa(i)+data.(string), dataSl[i:i+1])
		}

		wg.Wait()

		var longData string
		for _, str := range dataSl {
			longData += str
		}
		fmt.Println("LongData = ", longData)
		out <- longData
	}
}

func CombineResults(in, out chan interface{}) {
	var dataArr []string
LOOP:
	for data := range in {
		if data == nil {
			sort.Strings(dataArr)
			var result string
			for i, res := range dataArr {
				result += res
				if i != len(dataArr)-1 {
					result += "_"
				}
			}
			out <- result
			fmt.Println("Combine results closed")
			break LOOP
		}

		fmt.Println("DataArr = ", dataArr)
		dataArr = append(dataArr, data.(string))
		//sort.Strings(dataArr)
	}
}

func Crc32Worker(wg *sync.WaitGroup, data string, slice []string) {
	crc32 := DataSignerCrc32(data)
	slice[0] = crc32
	//fmt.Println(crc32)
	wg.Done()
}
