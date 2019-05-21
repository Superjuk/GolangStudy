package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := regexp.MustCompile("@")
	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := ""

	lines := strings.Split(string(fileContents), "\n")

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
			if ok, err := regexp.MatchString("Android", browser); ok && err == nil {
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

			if ok, err := regexp.MatchString("MSIE", browser); ok && err == nil {
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
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", index, nameRaw, email)
	}

	for i, line := range lines {
		parseLine(i, line)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

//go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1
