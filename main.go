package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"sync"

	"perkawis/config"
	"perkawis/src"
)

func main() {
	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Infof(".env is not loaded properly")
	}

	cfg := config.NewConfig()
	srv := src.InitServer(cfg)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		srv.Run()
	}()

	wg.Wait()
}
