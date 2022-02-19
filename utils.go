package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getDraws() string {
	var result string

	c := colly.NewCollector()

	c.OnHTML("ul.lottery-number-list.lottery-number-list2", func(e *colly.HTMLElement) {
		e.ForEach("li.active", func(i int, h *colly.HTMLElement) {
			result += h.Text + " "
		})
	})

	c.Visit("https://www.theb2blotto.com/home")

	return strings.TrimSpace(result)
}

func getLastDraw() string {
	lastDraw := ""

	uri := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("scheduler").Collection("draws")

	cursor, err := coll.Find(context.TODO(), bson.D{{}})

	var results []map[string]string

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	if len(results) > 0 {
		lastDraw = results[len(results)-1]["draw"]
	}

	if err != nil {
		panic(err)
	}

	return lastDraw
}

func saveDraw(draw string) {
	uri := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("scheduler").Collection("draws")

	doc := bson.D{{"draw", draw}}

	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}

func sendSms(message string) {
	lastDraw := getLastDraw()

	if message != lastDraw {
		data, _ := json.Marshal(map[string]string{
			"phone":   "+2330540810791",
			"message": message,
			"key":     "textbelt",
		})

		requestBody := bytes.NewBuffer(data)

		resp, err := http.Post("https://textbelt.com/text", "application/json", requestBody)

		if err != nil {
			panic(err)
		} else {
			saveDraw(message)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		log.Println(string(body))
	}
}

// func getCollection() *mongo.Collection {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	uri := os.Getenv("MONGO_URI")

// 	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer func() {
// 		if err := client.Disconnect(context.TODO()); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	coll := client.Database("scheduler").Collection("draws")

// 	return coll
// }
