package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "I'm up")
	})

	// TODO: repeat these two lines every 5 minutes
	draws := getDraws()
	sendSms(draws)

	var port string
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		port = "8000"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
