package log

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

const blue = "1;34"
const red = "1;31"
const yellow = "1;33"
const dtFormat = "15:04:05.000"

var mtx sync.Mutex
var stdout io.Writer = os.Stdout
var stderr io.Writer = os.Stderr

var debug = true

func doPrint(w io.Writer, prefix string, s string, args ...string) {
	dt := time.Now().Format(dtFormat)
	if s != "" {
		_, _ = fmt.Fprintf(w, "%s %s %s\n", dt, prefix, s)
	}
	for _, v := range args {
		_, _ = fmt.Fprintf(w, "                 %s\n", v)
	}
}

func doColorPrint(w io.Writer, color string, prefix string, s string, args ...string) {
	dt := time.Now().Format(dtFormat)
	if s != "" {
		_, _ = fmt.Fprintf(w, "\u001B[%sm%s %s %s\u001B[0m\n", color, dt, prefix, s)
	}
	for _, v := range args {
		_, _ = fmt.Fprintf(w, "\u001B[%sm                 %s\u001B[0m\n", color, v)
	}
}

func D(s string, args ...string) {
	if !debug {
		return
	}
	defer mtx.Unlock()
	mtx.Lock()
	doPrint(stdout, "[D]", s, args...)
}

func I(s string, args ...string) {
	defer mtx.Unlock()
	mtx.Lock()
	doColorPrint(stderr, blue, "[I]", s, args...)
}

func W(s string, args ...string) {
	defer mtx.Unlock()
	mtx.Lock()
	doColorPrint(stderr, yellow, "[W]", s, args...)
}

func E(s string, args ...string) {
	defer mtx.Unlock()
	mtx.Lock()
	doColorPrint(stderr, red, "[E]", s, args...)
}

func F(s string, args ...string) {
	defer mtx.Unlock()
	mtx.Lock()
	doColorPrint(stderr, red, "[F]", s, args...)
	os.Exit(1)
}

func CaptureLogOutput(f func()) string {
	defer resetOutput()

	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	setOutput(writer)

	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, reader)
		out <- buf.String()
	}()

	f()
	_ = writer.Close()
	ret := <-out
	_ = reader.Close()

	return ret
}

func EnableDebug() {
	defer mtx.Unlock()
	mtx.Lock()
	debug = true
}

func DisableDebug() {
	defer mtx.Unlock()
	mtx.Lock()
	debug = false
}

func EnableOutput() {
	resetOutput()
}

func DisableOutput() {
	setOutput(ioutil.Discard)
}

func resetOutput() {
	defer mtx.Unlock()
	mtx.Lock()
	stdout = os.Stdout
	stderr = os.Stderr
}

func setOutput(writer io.Writer) {
	defer mtx.Unlock()
	mtx.Lock()
	stdout = writer
	stderr = writer
}
