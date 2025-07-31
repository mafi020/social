package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mafi020/social/internal/env"
)

func init(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func main() {
	cfg := config{
		addr: env.GetEnvOrPanic("ADDR"),
	}

	app := &application{
		config: cfg,
	}
	
	log.Fatal(app.start(app.mount()))
}