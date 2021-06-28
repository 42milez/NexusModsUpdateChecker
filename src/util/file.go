package util

import (
	"os"
	"regexp"
	"strings"
)
import Err "github.com/42milez/NexusModsWatcher/src/error"

var WorkDir string

func Unpack(srcFile string, dstDir string) ([]string, error) {
	typ := strings.TrimPrefix(regexp.MustCompile(`\.(zip|7z)$`).FindString(srcFile), ".")

	if typ == "zip" {
		return UnpackZip(srcFile, dstDir)
	}

	if typ == "7z" {
		return Unpack7z(srcFile, dstDir)
	}

	return nil, Err.UnsupportedArchiveFormat
}

func init() {
	if wd, err := os.Getwd(); err != nil {
		Exit(Err.GetWorkingDirectoryFailed)
	} else {
		WorkDir = wd
	}
}
