package util

import (
	"os"
	"regexp"
	"strings"
)
import Err "github.com/42milez/NexusModsWatcher/src/error"

var WorkDir string

func Unpack(srcFile string, dstDir string) ([]string, error) {
	typ := strings.TrimPrefix(regexp.MustCompile(`\.(7z|rar|zip)$`).FindString(srcFile), ".")

	if typ == "7z" {
		return Unpack7z(srcFile, dstDir)
	}
	if typ == "rar" {
		return UnpackRar(srcFile, dstDir)
	}
	if typ == "zip" {
		return UnpackZip(srcFile, dstDir)
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
