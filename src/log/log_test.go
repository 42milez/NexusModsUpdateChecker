package log

import (
	"regexp"
	"strings"
	"testing"
)

func TestD_1(t *testing.T) {
	want, _ := regexp.Compile(`^[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{3} \[D] Debug$`)
	got := CaptureLogOutput(func() {
		D("Debug")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("D() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile(`^[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{3} \[D] DebugHelloWorld$`)
	got = CaptureLogOutput(func() {
		D("Debug", "Hello", "World")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("D() = %v; want %v", got, want.String())
	}
}

// return early when the debug flag is false
func TestD_2(t *testing.T) {
	defer EnableDebug()
	DisableDebug()

	want := ""
	got := CaptureLogOutput(func() {
		D("Debug")
	})
	got = Trim(got)
	if got != want {
		t.Errorf("D() = %s; want %s", got, want)
	}
}

func TestI(t *testing.T) {
	want, _ := regexp.Compile(`^\[1;34m[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{3} \[I] Info\[0m$`)
	got := CaptureLogOutput(func() {
		I("Info")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("I() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile(`^\[1;34m[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{3} \[I] Info\[0m\[1;34mHello\[0m\[1;34mWorld\[0m$`)
	got = CaptureLogOutput(func() {
		I("Info", "Hello", "World")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("I() = %v; want %v", got, want.String())
	}
}

func TestW(t *testing.T) {
	want, _ := regexp.Compile(`^\[1;33m[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{3} \[W] Warning\[0m$`)
	got := CaptureLogOutput(func() {
		W("Warning")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("W() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile(`^\[1;33m[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{3} \[W] Warning\[0m\[1;33mHello\[0m\[1;33mWorld\[0m$`)
	got = CaptureLogOutput(func() {
		W("Warning", "Hello", "World")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("W() = %v; want %v", got, want.String())
	}
}

func TestE(t *testing.T) {
	want, _ := regexp.Compile(`^\[1;31m[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{3} \[E] Error\[0m$`)
	got := CaptureLogOutput(func() {
		E("Error")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("E() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile(`^\[1;31m[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{3} \[E] Error\[0m\[1;31mHello\[0m\[1;31mWorld\[0m$`)
	got = CaptureLogOutput(func() {
		E("Error", "Hello", "World")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("E() = %v; want %v", got, want.String())
	}
}

func TestF(t *testing.T) {
	// No tests exists because log.Fatal*() calls os.Exit(1).
}

func Trim(s string) (ret string) {
	ret = strings.Replace(s, "\n", "", -1)
	ret = strings.Replace(ret, "                 ", "", -1)
	return
}
