package util

import "github.com/42milez/NexusModsWatcher/src/log"

func Exit(err error) {
	log.F(err.Error())
}
