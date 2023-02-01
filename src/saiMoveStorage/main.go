package main

import (
	"github.com/saiset-co/saiService"
	"github.com/webmakom-com/saiMoveStorage/internal"
)

func main() {
	svc := saiService.NewService("sai_VM1")
	svc.RegisterConfig("config.yml")

	is := internal.Service(svc)

	svc.RegisterHandlers(
		is.Handlers(),
	)

	svc.Start()
}
