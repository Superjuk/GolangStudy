package main

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	dummy        = make(chan interface{})
	dataArr      []string
	dataArrMutex = &sync.Mutex{}
	wgGl         = &sync.WaitGroup{}
	md5Mutex     = &sync.Mutex{}
	numsStrArr   = [...]string{"0", "1", "2", "3", "4", "5"}
)

func ExecutePipeline(jobs ...job) {
	rawDataCh := make(chan interface{}, MaxInputDataLen)

	wgGl.Add(1)
	go jobWrapper(jobs[0], dummy, rawDataCh)

	jobsSlice := jobs[1:]

	for num := range rawDataCh {
		var chans []chan interface{}
		for range jobsSlice {
			wgGl.Add(1)
			chans = append(chans, make(chan interface{}, MaxInputDataLen))
		}
		chans = append(chans, dummy)

		chans[0] <- num

		for i, job := range jobsSlice {
			go jobWrapper(job, chans[i], chans[i+1])
		}

		close(chans[0])
	}

	wgGl.Wait()
}

/**/
func jobWrapper(jb job, in, out chan interface{}) {
	defer wgGl.Done()
	if out != dummy {
		defer close(out)
	}

	jb(in, out)

	runtime.Gosched()
}

/*SingleHash count MD5 and CRC32*/
func SingleHash(in, out chan interface{}) {
LOOP:
	for {
		select {
		case data := <-in:
			if data == nil {
				out <- data
				break LOOP
			}

			dataStr := strconv.Itoa(data.(int))

			wg := &sync.WaitGroup{}
			wg.Add(2)

			dataSl := make([]string, 2)

			md5 := func(str string) *string {
				md5Mutex.Lock()
				md5Out := DataSignerMd5(str)
				md5Mutex.Unlock()
				return &md5Out
			}

			go Crc32Worker(wg, md5(dataStr), &dataSl[1])
			go Crc32Worker(wg, &dataStr, &dataSl[0])

			wg.Wait()

			out <- strings.Join(dataSl, "~")

			runtime.Gosched()
		}
	}
}

func MultiHash(in, out chan interface{}) {
LOOP:
	for {
		select {
		case data := <-in:
			if data == nil {
				out <- data
				break LOOP
			}

			wg := &sync.WaitGroup{}
			wg.Add(6)

			dataSl := make([]string, 6)

			for i, num := range numsStrArr {
				line := num + data.(string)
				go Crc32Worker(wg, &line, &dataSl[i])
			}

			wg.Wait()

			out <- strings.Join(dataSl, "")

			runtime.Gosched()
		}
	}
}

func CombineResults(in, out chan interface{}) {
LOOP:
	for data := range in {
		fmt.Println("Combine results fires!")
		fmt.Println("DataArr len=", len(dataArr))
		if data == nil {
			fmt.Println("data nil!")
			dataArrMutex.Lock()
			sort.Strings(dataArr)
			result := strings.Join(dataArr, "_")
			dataArrMutex.Unlock()
			out <- result
			break LOOP
		}

		dataArrMutex.Lock()
		dataArr = append(dataArr, data.(string))
		dataArrMutex.Unlock()

		runtime.Gosched()
	}
}

func Crc32Worker(wg *sync.WaitGroup, data *string, out *string) {
	crc32 := DataSignerCrc32(*data)
	*out = crc32
	wg.Done()
	runtime.Gosched()
}
