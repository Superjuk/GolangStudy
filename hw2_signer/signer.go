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
	wgGl         = &sync.WaitGroup{}
	dataArr      []string
	dataArrMutex = &sync.Mutex{}
	md5Mutex     = &sync.Mutex{}
	testMutex    = &sync.Mutex{}
	numsStrArr   = [...]string{"0", "1", "2", "3", "4", "5"}
)

func ExecutePipeline(jobs ...job) {
	var chans []chan interface{}
	dummy := make(chan interface{})

	chans = append(chans, dummy)

	for range jobs {
		wgGl.Add(1)
		chans = append(chans, make(chan interface{}, MaxInputDataLen))
	}

	for i, job := range jobs {
		go jobWrapper(job, chans[i], chans[i+1])
	}

	close(chans[0])

	wgGl.Wait()
}

/**/
func jobWrapper(jb job, in, out chan interface{}) {
	defer wgGl.Done()
	defer close(out)

	jb(in, out)

	runtime.Gosched()
}

/*SingleHash count MD5 and CRC32*/
func SingleHash(in, out chan interface{}) {
	md5 := func(str string) string {
		md5Mutex.Lock()
		md5Out := DataSignerMd5(str)
		md5Mutex.Unlock()
		return md5Out
	}

	hash := func(str string, out chan interface{}, wgSH *sync.WaitGroup) {
		wg := &sync.WaitGroup{}
		wg.Add(2)

		dataSl := make([]string, 2)

		go Crc32Worker(wg, str, dataSl[0:1])
		go Crc32Worker(wg, md5(str), dataSl[1:2])

		wg.Wait()

		fmt.Println("SingleHash out =", strings.Join(dataSl, "~"))
		out <- strings.Join(dataSl, "~")

		wgSH.Done()

		runtime.Gosched()
	}

	wgSH := &sync.WaitGroup{}

	for data := range in {
		wgSH.Add(1)
		dataStr := strconv.Itoa(data.(int))

		go hash(dataStr, out, wgSH)
	}

	wgSH.Wait()

	fmt.Println("Single hash done!!!")
}

func MultiHash(in, out chan interface{}) {
	hash := func(str string, out chan interface{}, wgMh *sync.WaitGroup) {
		wg := &sync.WaitGroup{}
		wg.Add(6)

		dataSl := make([]string, 6)

		for i, num := range numsStrArr {
			line := num + str
			go Crc32Worker(wg, line, dataSl[i:i+1])
		}

		wg.Wait()

		fmt.Println("MultiHash out =", strings.Join(dataSl, ""))
		out <- strings.Join(dataSl, "")

		wgMh.Done()

		runtime.Gosched()
	}

	wgMh := &sync.WaitGroup{}

	for data := range in {
		wgMh.Add(1)

		go hash(data.(string), out, wgMh)
	}

	wgMh.Wait()

	fmt.Println("MultiHash done!!!")
}

func CombineResults(in, out chan interface{}) {
	for data := range in {
		dataArrMutex.Lock()
		dataArr = append(dataArr, data.(string))
		fmt.Println("dataArr =", dataArr)
		dataArrMutex.Unlock()

		runtime.Gosched()
	}

	fmt.Println("Combine results fires!")
	fmt.Println("DataArr len=", len(dataArr))

	sort.Strings(dataArr)
	result := strings.Join(dataArr, "_")
	fmt.Println("result =", result)

	out <- result
}

func Crc32Worker(wg *sync.WaitGroup, data string, slice []string) {
	crc32 := DataSignerCrc32(data)
	slice[0] = crc32
	wg.Done()
	runtime.Gosched()
}
