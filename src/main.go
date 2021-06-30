package main

import (
	"flag"
	"fmt"
	"github.com/42milez/NexusModsUpdateChecker/src/api"
	Err "github.com/42milez/NexusModsUpdateChecker/src/error"
	"github.com/42milez/NexusModsUpdateChecker/src/log"
	"github.com/42milez/NexusModsUpdateChecker/src/nexus"
	"github.com/42milez/NexusModsUpdateChecker/src/util"
	"time"
)

const ModYaml = "mod.yml"

var Release string
var createPrFlag = flag.Bool("c", true, "create pull request")

func CreatePullRequest(filesUpload []string, releases []*nexus.Release) (string, error) {
	var desc string
	var err error

	if desc, err = nexus.CreateReleaseNote(releases); err != nil {
		return "", Err.CreateReleaseNoteFailed
	}

	sub := fmt.Sprintf("New Release %s", time.Now().Format("2006-01-02"))
	branch := fmt.Sprintf("new%d", time.Now().Unix())
	files := filesUpload
	msg := sub
	var url string

	if url, err = api.GitHub.CreatePullRequest(sub, desc, branch, files, msg); err != nil {
		return "", Err.CreatePullRequestFailed
	}

	return url, nil
}

func Download(mods nexus.ModInfoSet) error {
	if _, err := api.NexusMods.Download(mods); err != nil {
		return Err.DownloadFailed
	}
	return nil
}

func init() {
	if Release == "true" {
		log.DisableDebug()
	}
}

func main() {
	flag.Parse()

	var downloadCmd bool
	var updateCmd bool

	switch flag.Args()[0] {
	case "download":
		downloadCmd = true
	case "update":
		updateCmd = true
	default:
		util.Exit(Err.UnsupportedCommand)
	}

	mods := nexus.ModInfoSet{}
	var releases []*nexus.Release
	var err error

	if err = mods.Setup(ModYaml); err != nil {
		util.Exit(err)
	}

	if downloadCmd {
		if err = Download(mods); err != nil {
			util.Exit(err)
		}
	}

	if updateCmd {
		if releases, err = api.NexusMods.GetRelease(mods); err != nil {
			util.Exit(err)
		}

		for _, r := range releases {
			if !mods.Update(r) {
				log.E(fmt.Sprintf("can't find mod info: id=%d, version=%s", r.Mod.ID, r.File.Version))
			}
		}

		if err = mods.Export(ModYaml); err != nil {
			util.Exit(fmt.Errorf("can't export %s", ModYaml))
		}

		if len(releases) == 0 {
			log.I("no new release found")
			return
		}

		if *createPrFlag {
			var prUrl string
			if prUrl, err = CreatePullRequest([]string{ModYaml}, releases); err != nil {
				util.Exit(Err.CreatePullRequestFailed)
			}
			log.I(fmt.Sprintf("pull request created: %s", prUrl))
		}
	}
}
