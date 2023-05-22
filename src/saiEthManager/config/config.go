package config

import (
	"fmt"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	HttpServer struct {
		Host            string
		Port            int
		Token           string
		EnableProfiling bool
		ProfilingPort   int64
	}
	Contract struct {
		Server   string
		ABI      string
		Address  string
		Private  string
		GasLimit uint64
	}
}

var config Configuration

func Load() {
	configErr := gonfig.GetConf("config/config.json", &config)

	if configErr != nil {
		fmt.Println("Config load error: ", configErr)
		panic(configErr)
	}
}

func Get() Configuration {
	return config
}
