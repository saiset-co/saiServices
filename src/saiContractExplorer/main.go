package main

import (
	"fmt"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/saiContractExplorer/config"
	"github.com/webmakom-com/saiContractExplorer/explorer"
	"github.com/webmakom-com/saiContractExplorer/server"
)

func main() {
	cfg := config.Load()
	srv := server.NewServer(cfg, true)
	exp := explorer.NewExplorer(cfg)

	go srv.WSProcess()

	if cfg.Geth.Socket.Enabled {
		go exp.Process()
	}

	if cfg.Geth.Web.Enabled {
		go exp.WProcess()
	}

	if cfg.EnableProfiling {
		mr := gin.Default()
		pprof.Register(mr)
		go mr.Run(fmt.Sprintf(":%d", cfg.ProfilingPort))
	}

	srv.Start()
}
