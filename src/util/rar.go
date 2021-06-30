package util

import (
	Err "github.com/42milez/NexusModsWatcher/src/error"
	"github.com/mholt/archiver/v3"
	"os"
	"strings"
)

func UnpackRar(srcFile string, dstDir string) ([]string, error) {
	var err error

	if _, err = os.Stat(dstDir); err == nil {
		return nil, Err.DirectoryAlreadyExists
	}

	if err = archiver.Unarchive(srcFile, dstDir); err != nil {
		return nil, Err.ExtractFailed
	}

	var filesIncluded []string
	prefix := strings.TrimPrefix(dstDir, WorkDir+"/")

	err = archiver.Walk(srcFile, func(f archiver.File) error {
		if !f.IsDir() {
			filesIncluded = append(filesIncluded, prefix+"/"+f.Name())
		}
		return nil
	})
	if err != nil {
		return nil, Err.ListArchivedContentFailed
	}

	return filesIncluded, nil
}
