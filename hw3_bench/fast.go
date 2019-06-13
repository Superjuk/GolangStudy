package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type Browsers struct {
	Browsers []string
	Email    string
	Name     string
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lineContents := bufio.NewReader(file)

	r := regexp.MustCompile("@")
	seenBrowsers := []string{}
	uniqueBrowsers := 0

	parseLine := func(index int, line string) {
		user := new(Browsers)
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		browsersArr := reflect.ValueOf(user).Elem()

		browsers := browsersArr.Field(0).Interface().([]string)

		for _, browser := range browsers {
			tempAndr := strings.Split(browser, "Android")
			if len(tempAndr) > 1 {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}

			tempMsie := strings.Split(browser, "MSIE")
			if len(tempMsie) > 1 {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			return
		}

		emailRaw := browsersArr.Field(1).Interface().(string)
		nameRaw := browsersArr.Field(2).Interface().(string)

		email := r.ReplaceAllString(emailRaw, " [at] ")
		fmt.Fprintf(out, "[%d] %s <%s>\n", index, nameRaw, email)
	}

	i := 0
	fmt.Fprintln(out, "found users:")
	for {
		line, err := lineContents.ReadString('\n')
		if err == io.EOF {
			break
		}

		parseLine(i, line)
		i++
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}

//go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1
//go tool pprof hw3_bench.test.exe cpu.out
//go tool pprof hw3_bench.test.exe mem.out
