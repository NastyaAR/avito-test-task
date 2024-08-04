package main

import (
	"avito-test-task/config"
	"avito-test-task/internal/app"
	"log"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal("can't read config file")
	}

	app.Run(cfg)
}
