package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/saiAuth/config"
	"github.com/webmakom-com/saiAuth/server"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger, err := zap.NewDevelopment(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		log.Fatal("error creating logger : ", err.Error())
	}
	logger.Debug("Logger started", zap.String("mode", "debug"))
	srv := server.NewServer(cfg, false, logger)

	if cfg.SocketServer.Host != "" {
		go srv.SocketStart()
	}

	if cfg.EnableProfiling {
		mr := gin.Default()
		pprof.Register(mr)
		go mr.Run(fmt.Sprintf(":%d", cfg.ProfilingPort))
	}

	//srv.StartHttps()
	srv.Start()
}
