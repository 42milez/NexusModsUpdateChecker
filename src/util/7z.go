package util

import (
	Err "github.com/42milez/NexusModsUpdateChecker/src/error"
	"github.com/gen2brain/go-unarr"
	"os"
	"strings"
)

func Unpack7z(srcFile string, dstDir string) ([]string, error) {
	var filesExtracted []string
	var reader *unarr.Archive
	var err error

	if _, err = os.Stat(dstDir); err == nil {
		return nil, Err.DirectoryAlreadyExists
	}

	if reader, err = unarr.NewArchive(srcFile); err != nil {
		return nil, Err.OpenFileFailed
	}
	defer close7zReader(reader)

	if err = os.MkdirAll(dstDir, 0755); err != nil {
		return nil, Err.CreateDirectoryFailed
	}

	var files []string

	if files, err = reader.Extract(dstDir); err != nil {
		return nil, Err.ExtractFailed
	}

	for _, f := range files {
		filesExtracted = append(filesExtracted, strings.TrimPrefix(dstDir, WorkDir+"/")+"/"+f)
	}

	return filesExtracted, nil
}

func close7zReader(r *unarr.Archive) {
	if r == nil {
		return
	}
	if err := r.Close(); err != nil {
		Exit(err)
	}
}
