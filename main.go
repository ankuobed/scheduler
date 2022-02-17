package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// load .env file

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "I'm up")
	})

	go func() {
		// run every 5 minutes
		for range time.Tick(time.Minute * 5) {
			draws := getDraws()
			sendSms(draws)
		}
	}()

	var port string
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		port = "8000"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
