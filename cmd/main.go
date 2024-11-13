package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AsaHero/dastyor-bot/internal/app"
	"github.com/AsaHero/dastyor-bot/pkg/config"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg := config.New()

	app := app.New(cfg)

	// run application
	go func() {
		log.Println("dastyor bot starting...")
		if err := app.Run(); err != nil {
			log.Fatalln("app run", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	// app stops
	log.Println("dastyor bot stopping...")
	app.Stop()
}
