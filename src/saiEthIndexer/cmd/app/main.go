package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/saiset-co/saiEthIndexer/internal/app"
	"github.com/saiset-co/saiEthIndexer/tasks"
)

func main() {
	args := os.Args

	app, err := app.New(args)
	if err != nil {
		log.Fatal(err)
	}

	//register config with specific options
	err = app.RegisterConfig("./config/config.json", "./config/contracts.json")
	if err != nil {
		log.Fatal(err)
	}

	t, err := tasks.NewManager(app.Cfg, app.Logger)
	if err != nil {
		log.Fatal(err)
	}

	defer t.Logger.Sync()

	app.RegisterTask(t)

	app.RegisterHandlers()

	if app.Cfg.Common.EnableProfiling {
		mr := gin.Default()
		pprof.Register(mr)
		go mr.Run(fmt.Sprintf(":%d", app.Cfg.Common.ProfilingPort))
	}

	app.Run()

}
