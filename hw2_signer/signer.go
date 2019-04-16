package main

import (
	//	"context"
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	md5Lock uint32 = 0
)

func ExecutePipeline(jobs ...job) {
	var chans []chan interface{}
	dummy := make(chan interface{})

	chans = append(chans, dummy)
	for i := 0; i < len(jobs); i++ {
		chans = append(chans, make(chan interface{}, MaxInputDataLen))
	}
	chans = append(chans, dummy)

	wg := &sync.WaitGroup{}

	fmt.Println("run ExecutePipline")

	for i, job := range jobs {
		wg.Add(1)
		go jobWrapper(job, chans[i], chans[i+1], wg)
	}

LOOP:
	for {
		select {
		case num := <-chans[0]:
			fmt.Println(num)

			if num == nil {
				close(chans[1])
				break LOOP
			}

			if numU32, isOk := num.(uint32); isOk == true {
				chans[1] <- numU32
			} else if num32, isOk := num.(int32); isOk == true {
				chans[1] <- num32
			} else if numU64, isOk := num.(uint64); isOk == true {
				chans[1] <- numU64
			} else if num64, isOk := num.(int64); isOk == true {
				chans[1] <- num64
			}

		default:
			break LOOP
		}
	}

	wg.Wait()

	fmt.Println("End ExecutePipeline")

}

func jobWrapper(jb job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)
	defer fmt.Println("job over")

	jb(in, out)
}

func SingleHash(in, out chan interface{}) {
	runtime.Gosched()
LOOP:
	for {
		select {
		case data := <-in:

			dataStr := strconv.Itoa(int(data.(uint32)))

			if dataStr == "the_end" {
				out <- "the_end"
				//fmt.Println("the_end")
				break LOOP
			}

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
		if data == "the_end" {
			out <- data
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
		fmt.Println(longData)
		out <- longData
	}
}

func CombineResults(in, out chan interface{}) {
	var dataArr []string
LOOP:
	for data := range in {
		if data == "the_end" {
			sort.Strings(dataArr)
			var result string
			for i, res := range dataArr {
				result += res
				if i != len(dataArr)-1 {
					result += "_"
				}
			}
			out <- result
			break LOOP
		}

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

func main() {
	//inputData := []int{0, 1, 1, 2, 3, 5, 8}
	//inputData := []int{0, 1}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			//			for _, fibNum := range inputData {
			//				out <- fibNum
			//			}
			for i := 0; i < 100; i++ {
				out <- i
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data := dataRaw.(string)
			fmt.Println(data)
		}),
	}

	start := time.Now()

	ExecutePipeline(hashSignJobs...)

	end := time.Since(start)

	fmt.Println("Time =", end)

	fmt.Scanln()
	//	start := time.Now()

	//	wg := &sync.WaitGroup{}
	//	wg.Add(1)

	//	chan_out := make(chan interface{}, 1)
	//	chan_in := make(chan interface{}, 1)
	//	chan_out2 := make(chan interface{}, 1)
	//	chan_in2 := make(chan interface{}, 1)

	//	go EndProgram(wg, chan_in2)

	//	go SingleHash(chan_out, chan_in)
	//	go MultiHash(chan_in, chan_out2)
	//	go CombineResults(chan_out2, chan_in2)

	//	chan_out <- 0
	//	chan_out <- 1
	//	chan_out <- 2
	//	chan_out <- 3
	//	chan_out <- 4
	//	chan_out <- 5
	//	chan_out <- 6

	//	wg.Wait()
	//	stop := time.Now()
	//	fmt.Println("Time =", stop.Sub(start))
}
