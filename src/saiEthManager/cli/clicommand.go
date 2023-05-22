package cli

import (
	"fmt"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/thatisuday/commando"
	"github.com/webmakom-com/saiEthManager/config"
	"github.com/webmakom-com/saiEthManager/server"
)

func InitCli() {
	commando.
		SetExecutableName("sai-eth-manager").
		SetVersion("1.0.0")

	commando.
		Register("start").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			config.Load()

			if config.Get().HttpServer.EnableProfiling {
				mr := gin.Default()
				pprof.Register(mr)
				go mr.Run(fmt.Sprintf(":%d", config.Get().HttpServer.ProfilingPort))
			}

			server.NewServer().Start()
		})

	commando.Parse(nil)
}
