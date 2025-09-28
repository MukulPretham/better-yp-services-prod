package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"mukulpretham/betterUpPublisher/utils"
	"mukulpretham/betterUpPublisher/redis_utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/redis/go-redis/v9"
)

func main() {

	//Connect to database and get all websites
	dsn := "host=better-up-postgres user=postgres password=9059015626 dbname=postgres port=5432"
	db, err := gorm.Open(postgres.Open(dsn))
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()
	if err != nil {
		log.Fatal("unable to connect ot the database")
	}
	//Creating redis client
	client := utils.CreateRedisClient("better-up-redis:6379", 0, "", 2)

	for {
		func(db *gorm.DB, client *redis.Client) {
			var cueeWebsites []utils.Website
			db.Find(&cueeWebsites)

			for _, rec := range cueeWebsites {
				data, err := json.Marshal(rec)
				if err != nil {
					log.Println("Failed to marshal:", err)
					continue
				}
				// Making it run in a go routine for concurrency
				go redis_utils.Xadd(client,data)
			}

		}(db, client)
		fmt.Println("iteration completed")
		time.Sleep(180 * time.Second)
	}
}
