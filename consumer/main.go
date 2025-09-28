package main

// TODOs
// handel consumer dead situation
// Notification feature.

// import (
// 	"fmt"
// 	"mukulpretham/betterUpConsumer/helpers"
// 	"os"

// 	"github.com/joho/godotenv"
// )

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"mukulpretham/betterUpPublisher/utils"

	"mukulpretham/betterUpConsumer/helpers"

	"github.com/joho/godotenv"
)

func main() {
	// Load the env file
	godotenv.Load(".env")

	//Cureent Consumer Group
	fmt.Print(os.Getenv("REGION"))
	currConsumerGroup := fmt.Sprintf("%sConsumerGroup",os.Getenv("REGION"))
	
	// Redis client created
	client := utils.CreateRedisClient("better-up-redis:6379",0,"",2)

	// Create redis consumerGroup of the paticular region and a stream, if not exist
	err := utils.CreateRedisGroup(client, "websites",currConsumerGroup )
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Read messages form redis stremas via a consumer group
		res,err := utils.ReadXGroup(client,[]string{"websites",">"},currConsumerGroup)
		if err != nil {
			log.Fatal(err)
		}

		msg := res[0]

		if currMesssage, ok := msg.Values["site"].(string); ok {
			var m map[string]string
			// Parsing to JSON.
			if err := json.Unmarshal([]byte(currMesssage), &m); err != nil {
				panic("error parsing string")
			}
			go helpers.WriteToDB(m["Url"],client,msg.ID)
		}
	}
}

