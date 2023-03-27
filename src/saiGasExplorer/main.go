package main

import (
	"fmt"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/hv/src/saiGasExplorer/config"
	"github.com/webmakom-com/hv/src/saiGasExplorer/server"
)

func main() {
	cfg := config.Load()
	srv := server.NewServer(cfg, true)

	if cfg.EnableProfiling {
		mr := gin.Default()
		pprof.Register(mr)
		go mr.Run(fmt.Sprintf(":%d", cfg.ProfilingPort))
	}

	srv.Start()
}
