package config

import (
	"fmt"
	"golang.org/x/exp/slog"
	. "ivory/src/model"
	"os"
	"regexp"
)

type Env struct {
	Config  Config
	Version Version
}

func NewEnv() *Env {
	ginMode := "debug"
	if val, ok := os.LookupEnv("GIN_MODE"); ok {
		ginMode = val
	}

	certFilePath := ""
	if val, ok := os.LookupEnv("IVORY_CERT_FILE_PATH"); ok {
		certFilePath = val
	}
	certKeyFilePath := ""
	if val, ok := os.LookupEnv("IVORY_CERT_KEY_FILE_PATH"); ok {
		certKeyFilePath = val
	}

	urlAddress := ":8080"
	if val, ok := os.LookupEnv("IVORY_URL_ADDRESS"); ok {
		urlAddress = val
	} else if ginMode == "release" {
		if certFilePath != "" && certKeyFilePath != "" {
			urlAddress = ":443"
		} else {
			urlAddress = ":80"
		}
	}
	urlPath := "/"
	if val, ok := os.LookupEnv("IVORY_URL_PATH"); ok {
		urlPath = parseUrlPath(val)
	}

	staticFilesPath := ""
	if val, ok := os.LookupEnv("IVORY_STATIC_FILES_PATH"); ok {
		staticFilesPath = val
		updateBaseUrlTag(staticFilesPath, urlPath)
	}

	tag := "v0.0.0"
	if val, ok := os.LookupEnv("IVORY_VERSION_TAG"); ok {
		tag = val
	}
	commit := "000000000000"
	if val, ok := os.LookupEnv("IVORY_VERSION_COMMIT"); ok {
		commit = val
	}

	println()
	slog.Info("IVORY ENV VARIABLES:")
	slog.Info("ENV", "GIN_MODE", ginMode)
	slog.Info("ENV", "IVORY_URL_ADDRESS", urlAddress)
	slog.Info("ENV", "IVORY_URL_PATH", urlPath)
	slog.Info("ENV", "IVORY_STATIC_FILES_PATH", staticFilesPath)
	slog.Info("ENV", "IVORY_CERT_FILE_PATH", certFilePath)
	slog.Info("ENV", "IVORY_CERT_KEY_FILE_PATH", certKeyFilePath)
	slog.Info("ENV", "IVORY_VERSION_TAG", tag)
	slog.Info("ENV", "IVORY_VERSION_COMMIT", commit)
	println()

	return &Env{
		Config: Config{
			UrlAddress:      urlAddress,
			UrlPath:         urlPath,
			StaticFilesPath: staticFilesPath,
			CertFilePath:    certFilePath,
			CertKeyFilePath: certKeyFilePath,
		},
		Version: Version{
			Tag:    tag,
			Commit: commit,
			Label:  "Ivory " + tag + " (" + fmt.Sprintf("%.7s", commit) + ")",
		},
	}
}

func updateBaseUrlTag(path string, href string) {
	// NOTE: we need to remove `/` because we set it in replace string
	//  correct url for base href is `/`, `/ivory/`, `/example/` (it should always be covered by two `/`
	if href == "/" {
		href = ""
	}
	indexPath := path + "/index.html"
	index, errRead := os.ReadFile(indexPath)
	if errRead != nil {
		panic(errRead)
	}
	file := string(index)
	regex := regexp.MustCompile("<base href=\".*\" />")
	newFile := regex.ReplaceAllString(file, "<base href=\""+href+"/\" />")
	errWrite := os.WriteFile(indexPath, []byte(newFile), 0666)
	if errWrite != nil {
		panic(errWrite)
	}
}

func parseUrlPath(path string) string {
	if path == "" {
		return "/"
	}
	newPath := path

	// NOTE: path should either be empty or with slash as first symbol
	if newPath[0] != '/' {
		panic("path should contain `/` at the beginning")
	}

	// NOTE we need to always remove last slash to handle routes correctly
	newPathLen := len(newPath) - 1
	if newPathLen > 0 && newPath[newPathLen] == '/' {
		newPath = newPath[0 : len(newPath)-1]
	}

	return newPath
}
