package main

import (
	"runtime"
	"sort"
	"strconv"
	"sync"
)

var (
	md5Lock      uint32 = 0
	dummy               = make(chan interface{})
	dataArr      []string
	dataArrMutex = &sync.Mutex{}
	wgGl         = &sync.WaitGroup{}
	md5Mutex     = &sync.Mutex{}
	jobMutex     = &sync.Mutex{}
	numsStrArr   = [...]string{"0", "1", "2", "3", "4", "5"}
)

func ExecutePipeline(jobs ...job) {
	rawDataCh := make(chan interface{}, MaxInputDataLen)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go jobWrapper(jobs[0], dummy, rawDataCh, wg)

	jobsSlice := jobs[1:]
	jobsSlCount := len(jobsSlice)

	for num := range rawDataCh {
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
}

func jobWrapper(jb job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	if out != dummy {
		defer close(out)
	}

	jobMutex.Lock()
	jb(in, out)
	jobMutex.Unlock()

	runtime.Gosched()
}

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

			md5 := func(str string) string {
				md5Mutex.Lock()
				md5Out := DataSignerMd5(str)
				md5Mutex.Unlock()
				return md5Out
			}

			go Crc32Worker(wg, md5(dataStr), dataSl[1:2])
			go Crc32Worker(wg, dataStr, dataSl[0:1])

			wg.Wait()

			out <- dataSl[0] + "~" + dataSl[1]

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
				go Crc32Worker(wg, num+data.(string), dataSl[i:i+1])
			}

			wg.Wait()

			var longData string
			for _, str := range dataSl {
				longData += str
			}
			out <- longData

			runtime.Gosched()
		}
	}
}

func CombineResults(in, out chan interface{}) {
LOOP:
	for data := range in {
		if data == nil {
			dataArrMutex.Lock()
			sort.Strings(dataArr)
			var result string
			for i, res := range dataArr {
				result += res
				if i != len(dataArr)-1 {
					result += "_"
				}
			}
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

func Crc32Worker(wg *sync.WaitGroup, data string, slice []string) {
	crc32 := DataSignerCrc32(data)
	slice[0] = crc32
	wg.Done()
	runtime.Gosched()
}

func main() {

}
