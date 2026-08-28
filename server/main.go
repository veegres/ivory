package main

import (
	"ivory/clients/console/shell"
	"ivory/clients/console/ssh"
	"ivory/clients/http"
	"ivory/core"
	"ivory/core/config"
	"ivory/core/engine"
	"ivory/features"
	"ivory/plugins"
	"ivory/tools"
)

func main() {
	appEnv := config.NewEnvironment()

	httpClient := http.NewClient()
	sshClient := ssh.NewClient()
	shellClient := shell.NewClient()

	cc := core.NewContext(sshClient)
	pc := plugins.NewContext(httpClient, sshClient)
	tc := tools.NewContext(shellClient, cc.Service)
	fc := features.NewContext(
		appEnv,
		pc.DatabaseRegistry,
		pc.PlatformRegistry,
		pc.KeeperRegistry,
		tc.Registry,
		cc.Service,
	)

	engine.NewHttpServer(appEnv, cc.Router, fc.Router, tc.Router)
}
