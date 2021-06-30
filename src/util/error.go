package util

import "github.com/42milez/NexusModsUpdateChecker/src/log"

func Exit(err error) {
	log.F(err.Error())
}
