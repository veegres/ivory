package config

import (
	"log"
	"os"
)

type Env struct {
	Company  string
	Auth     string
	Tag      string
	Commit   string
	Username *string
	Password *string
}

func NewEnv() *Env {
	// TODO think how to move this to UI set up page it can be next after secret page
	company := "Veegres"
	if val, ok := os.LookupEnv("IVORY_COMPANY_LABEL"); ok {
		company = val
	}
	auth := "none"
	if val, ok := os.LookupEnv("IVORY_AUTHENTICATION"); ok {
		auth = val
	}
	tag := "none"
	if val, ok := os.LookupEnv("IVORY_VERSION_TAG"); ok {
		tag = val
	}
	commit := "none"
	if val, ok := os.LookupEnv("IVORY_VERSION_COMMIT"); ok {
		commit = val
	}
	var user *string
	if val, ok := os.LookupEnv("IVORY_BASIC_USER"); ok && auth == "basic" {
		user = &val
	}
	var password *string
	if val, ok := os.LookupEnv("IVORY_BASIC_PASSWORD"); ok && auth == "basic" {
		password = &val
	}

	if auth == "basic" {
		if user == nil || password == nil {
			panic("Provide IVORY_BASIC_USER and IVORY_BASIC_PASSWORD when you use IVORY_AUTHENTICATION=basic")
		}
	}

	log.Println("IVORY ENV VARIABLES")
	log.Println("IVORY_VERSION_TAG:", tag)
	log.Println("IVORY_VERSION_COMMIT:", commit)
	log.Println("IVORY_COMPANY_LABEL:", company)
	log.Println("IVORY_AUTHENTICATION:", auth)
	log.Println("IVORY_BASIC_USER:", user != nil)
	log.Println("IVORY_BASIC_PASSWORD:", password != nil)

	return &Env{
		Company:  company,
		Auth:     auth,
		Tag:      tag,
		Commit:   commit,
		Username: user,
		Password: password,
	}
}
