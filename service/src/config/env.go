package config

import (
	"fmt"
	"golang.org/x/exp/slog"
	. "ivory/src/model"
	"os"
)

type Env struct {
	Version Version
}

func NewEnv() *Env {
	tag := "v0.0.0"
	if val, ok := os.LookupEnv("IVORY_VERSION_TAG"); ok {
		tag = val
	}
	commit := "000000000000"
	if val, ok := os.LookupEnv("IVORY_VERSION_COMMIT"); ok {
		commit = val
	}

	slog.Info("IVORY ENV VARIABLES")
	slog.Info("ENV", "IVORY_VERSION_TAG", tag)
	slog.Info("ENV", "IVORY_VERSION_COMMIT", commit)

	return &Env{
		Version: Version{
			Tag:    tag,
			Commit: commit,
			Label:  "Ivory " + tag + " (" + fmt.Sprintf("%.7s", commit) + ")",
		},
	}
}
