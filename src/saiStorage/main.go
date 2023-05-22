package main

import (
	"fmt"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/saiStorage/config"
	"github.com/webmakom-com/saiStorage/mongo"
	"github.com/webmakom-com/saiStorage/server"
)

func main() {
	cfg := config.Load()
	srv := server.NewServer(cfg, false)
	mSrv := mongo.NewMongoServer(cfg)

	go mSrv.Start()

	if cfg.EnableProfiling {
		mr := gin.Default()
		pprof.Register(mr)
		go mr.Run(fmt.Sprintf(":%d", cfg.ProfilingPort))
	}

	srv.Start()
	//srv.StartHttps()
}
