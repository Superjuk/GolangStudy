package main

import (
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
	numsStrArr   = [...]string{"0", "1", "2", "3", "4", "5"}
)

func ExecutePipeline(jobs ...job) {
	var chans []chan interface{}
	dummy := make(chan interface{})

	chans = append(chans, dummy)

	for range jobs {
		chans = append(chans, make(chan interface{}))
	}

	wgGl.Add(len(jobs))

	for i, job := range jobs {
		go jobWrapper(job, chans[i], chans[i+1])
	}

	close(chans[0])

	wgGl.Wait()
}

func jobWrapper(jb job, in, out chan interface{}) {
	defer wgGl.Done()
	defer close(out)

	jb(in, out)

	runtime.Gosched()
}

func SingleHash(in, out chan interface{}) {
	hash := func(data interface{}, wgSH *sync.WaitGroup) {
		wg := &sync.WaitGroup{}
		wg.Add(1)

		str := strconv.Itoa(data.(int))

		dataSl := make([]string, 2)

		go Md5Worker(wg, str, dataSl)

		wg.Wait()

		out <- strings.Join(dataSl, "~")

		wgSH.Done()

		runtime.Gosched()
	}

	wgSH := &sync.WaitGroup{}

	for data := range in {
		wgSH.Add(1)
		go hash(data, wgSH)
	}

	wgSH.Wait()
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
}

func CombineResults(in, out chan interface{}) {
	for data := range in {
		dataArrMutex.Lock()
		dataArr = append(dataArr, data.(string))
		dataArrMutex.Unlock()

		runtime.Gosched()
	}

	sort.Strings(dataArr)
	result := strings.Join(dataArr, "_")

	out <- result
}

func Crc32Worker(wg *sync.WaitGroup, data string, slice []string) {
	crc32 := DataSignerCrc32(data)
	slice[0] = crc32
	wg.Done()
	runtime.Gosched()
}

func Md5Worker(wg *sync.WaitGroup, str string, slice []string) {
	wgMd5 := &sync.WaitGroup{}
	md5Out := make(chan string)
	out := make(chan string)

	md5 := func(str string) {
		out <- str
		md5Mutex.Lock()
		md5Out <- DataSignerMd5(str)
		md5Mutex.Unlock()
	}

	crc32 := func(in chan string, index int) {
		data := <-in
		crc32 := DataSignerCrc32(data)
		slice[index] = crc32
		wgMd5.Done()
	}

	wgMd5.Add(2)

	go md5(str)
	go crc32(out, 0)
	go crc32(md5Out, 1)

	go func() {
		wgMd5.Wait()
		wg.Done()
		close(md5Out)
	}()

	runtime.Gosched()
}
