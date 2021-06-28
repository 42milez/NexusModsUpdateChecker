package main

import (
	"fmt"
	"github.com/42milez/NexusModsWatcher/src/api"
	Err "github.com/42milez/NexusModsWatcher/src/error"
	"github.com/42milez/NexusModsWatcher/src/log"
	"github.com/42milez/NexusModsWatcher/src/nexus"
	"github.com/42milez/NexusModsWatcher/src/util"
	"time"
)

const ModYaml = "mod.yml"

var Release string

func DownloadNewRelease(releases []*nexus.Release, mods *nexus.ModInfoSet) ([]string, error) {
	var lcFiles []*nexus.LocalFile
	var err error

	if lcFiles, err = api.NexusMods.Download(releases); err != nil {
		return nil, Err.DownloadFailed
	}

	var filesDownloaded []string

	for _, f := range lcFiles {
		domain := f.Release.Domain
		modId := f.Release.Mod.ID
		fileVer := f.Release.File.Version
		dstDir := fmt.Sprintf("%s/%s/%s/%s", util.WorkDir, domain, f.Release.Mod.Path, fileVer)
		var filesExtracted []string

		if filesExtracted, err = util.Unpack(f.Path, dstDir); err != nil {
			return nil, Err.ExtractFailed
		}

		if len(filesExtracted) > 0 {
			filesDownloaded = append(filesDownloaded, filesExtracted...)
		}

		log.D(fmt.Sprintf("file extracted: %s", dstDir))

		if !mods.Update(domain, modId, fileVer) {
			log.E(fmt.Sprintf("can't find mod info: id=%d, version=%s", modId, fileVer))
		}
	}

	return filesDownloaded, nil
}

func CreatePullRequest(filesUpload []string, releases []*nexus.Release) (string, error) {
	var desc string
	var err error

	if desc, err = nexus.CreateReleaseNote(releases); err != nil {
		return "", Err.CreateReleaseNoteFailed
	}

	sub := fmt.Sprintf("New Release %s", time.Now().Format("2006-01-02"))
	branch := fmt.Sprintf("new%d", time.Now().Unix())
	files := filesUpload
	msg := "new release"
	var url string

	if url, err = api.GitHub.CreatePullRequest(sub, desc, branch, files, msg); err != nil {
		return "", Err.CreatePullRequestFailed
	}

	return url, nil
}

func init() {
	if Release == "true" {
		log.DisableDebug()
	}
}

func main() {
	var err error

	mods := &nexus.ModInfoSet{}

	if err = mods.Setup(ModYaml); err != nil {
		util.Exit(err)
	}

	var releases []*nexus.Release

	if releases, err = api.NexusMods.GetRelease(mods); err != nil {
		util.Exit(err)
	}

	if len(releases) == 0 {
		log.I("no new release found")
		return
	}

	var filesUpload []string

	if filesUpload, err = DownloadNewRelease(releases, mods); err != nil {
		util.Exit(err)
	}

	if err = mods.Export(ModYaml); err != nil {
		util.Exit(fmt.Errorf("can't export %s", ModYaml))
	}
	filesUpload = append(filesUpload, ModYaml)

	log.D("files upload", filesUpload...)

	var prUrl string

	if prUrl, err = CreatePullRequest(filesUpload, releases); err != nil {
		util.Exit(Err.CreatePullRequestFailed)
	}

	log.I(fmt.Sprintf("pull request created: %s", prUrl))
}
