package config

import (
	. "ivory/src/model"
	"log"
	"os"
)

type Env struct {
	Version Version
}

func NewEnv() *Env {
	tag := "none"
	if val, ok := os.LookupEnv("IVORY_VERSION_TAG"); ok {
		tag = val
	}
	commit := "none"
	if val, ok := os.LookupEnv("IVORY_VERSION_COMMIT"); ok {
		commit = val
	}

	log.Println("IVORY ENV VARIABLES")
	log.Println("IVORY_VERSION_TAG:", tag)
	log.Println("IVORY_VERSION_COMMIT:", commit)

	return &Env{
		Version: Version{
			Tag:    tag,
			Commit: commit,
		},
	}
}
