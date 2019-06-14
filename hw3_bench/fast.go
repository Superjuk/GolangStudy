package main

import (
	"bufio"
	json "encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

type Browsers struct {
	Browsers []string
	Email    string
	Name     string
}

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson89aae3efDecodeGolangStudyHw3BenchEasyjson(in *jlexer.Lexer, out *Browsers) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "email":
			out.Email = string(in.String())
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson89aae3efEncodeGolangStudyHw3BenchEasyjson(out *jwriter.Writer, in Browsers) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"browsers\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Browsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Browsers {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Browsers) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson89aae3efEncodeGolangStudyHw3BenchEasyjson(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Browsers) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson89aae3efEncodeGolangStudyHw3BenchEasyjson(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Browsers) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson89aae3efDecodeGolangStudyHw3BenchEasyjson(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Browsers) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson89aae3efDecodeGolangStudyHw3BenchEasyjson(l, v)
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lineContents := bufio.NewReader(file)

	seenBrowsers := []string{}
	uniqueBrowsers := 0

	parseLine := func(index int, line string) {
		user := new(Browsers)
		err := user.UnmarshalJSON([]byte(line))
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		browsers := user.Browsers

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

		emailRaw := user.Email
		name := user.Name

		email := strings.Replace(emailRaw, "@", " [at] ", 1)
		fmt.Fprintf(out, "[%d] %s <%s>\n", index, name, email)
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
