package main

import (
	"flag"
	"log"
	"strings"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/saiAuth/config"
	"github.com/webmakom-com/saiAuth/server"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	var mode string
	flag.StringVar(&mode, "mode", "production", "change mode in service")
	flag.Parse()

	var logger *zap.Logger
	var err error

	if strings.ToLower(mode) == "debug" {
		logger, err = zap.NewDevelopment(zap.AddStacktrace(zap.DPanicLevel))
		if err != nil {
			log.Fatal("error creating logger : ", err.Error())
		}
		logger.Info("Logger started", zap.String("mode", "debug"))
	} else {
		logger, err = zap.NewProduction(zap.AddStacktrace(zap.DPanicLevel))
		if err != nil {
			log.Fatal("error creating logger : ", err.Error())
		}
		logger.Info("Logger started", zap.String("mode", "production"))
	}

	srv := server.NewServer(cfg, false, logger)

	if cfg.SocketServer.Host != "" {
		go srv.SocketStart()
	}

	//server for metrics
	mr := gin.Default()
	pprof.Register(mr)
	go mr.Run()
	//

	//srv.StartHttps()
	srv.Start()
}

type Monitor struct {
	Alloc,
	TotalAlloc,
	Sys,
	Mallocs,
	Frees,
	LiveObjects,
	PauseTotalNs uint64

	NumGC        uint32
	NumGoroutine int
}
