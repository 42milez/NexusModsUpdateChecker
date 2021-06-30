package util

import (
	"archive/zip"
	"fmt"
	Err "github.com/42milez/NexusModsWatcher/src/error"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnpackZip(srcFile string, dstDir string) ([]string, error) {
	extract := func(f *zip.File) (string, error) {
		var origFile io.ReadCloser
		var err error

		if origFile, err = f.Open(); err != nil {
			return "", Err.OpenFileFailed
		}
		defer closeIoReadCloser(origFile)

		var path string

		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(fmt.Sprintf("%s/%s", dstDir, f.Name), 0755); err != nil {
				return "", Err.CreateDirectoryFailed
			}
		} else {
			var dstFile *os.File
			flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
			dstPath := fmt.Sprintf("%s/%s", dstDir, f.Name)

			if _, err = os.Stat(dstPath); err != nil {
				if err = os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
					return "", Err.CreateDirectoryFailed
				}
			}

			if dstFile, err = os.OpenFile(dstPath, flag, f.Mode()); err != nil {
				return "", Err.OpenFileFailed
			}
			defer closeFile(dstFile)

			if _, err = io.Copy(dstFile, origFile); err != nil {
				return "", Err.CopyFileFailed
			}

			path = dstFile.Name()
		}

		return path, nil
	}

	var err error

	if _, err = os.Stat(dstDir); err == nil {
		return nil, Err.DirectoryAlreadyExists
	}

	if err = os.MkdirAll(dstDir, 0755); err != nil {
		return nil, Err.CreateDirectoryFailed
	}

	reader, err := zip.OpenReader(srcFile)
	if err != nil {
		return nil, Err.OpenFileFailed
	}
	defer closeZipReader(reader)

	var filesExtracted []string

	for _, f := range reader.Reader.File {
		filePath, errE := extract(f)
		if errE != nil {
			return nil, Err.ExtractFailed
		}
		if len(filePath) > 0 {
			filesExtracted = append(filesExtracted, strings.TrimPrefix(filePath, WorkDir+"/"))
		}
	}

	return filesExtracted, nil
}

func closeFile(f *os.File) {
	if f == nil {
		return
	}
	if err := f.Close(); err != nil {
		Exit(err)
	}
}

func closeIoReadCloser(r io.ReadCloser) {
	if r == nil {
		return
	}
	if err := r.Close(); err != nil {
		Exit(err)
	}
}

func closeZipReader(r *zip.ReadCloser) {
	if r == nil {
		return
	}
	if err := r.Close(); err != nil {
		Exit(err)
	}
}
