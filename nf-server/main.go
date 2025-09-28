package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mukulpretham/betterUpConsumer/helpers"
	"mukulpretham/betterUpPublisher/utils"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	client := utils.CreateRedisClient("better-up-redis:6379", 0, "", 2)
	
	err := utils.CreateRedisGroup(client, "notifications", "notificationGroup")
	if err != nil {
		log.Fatal(err)
	}else{
		fmt.Println("group already exist or it has been created")
	}

	for {
		readRes, readErr := utils.ReadXGroup(client, []string{"notifications", ">"}, "notificationGroup")
		if readErr != nil {
			log.Fatal(readErr)
		}
		currMessage := readRes[0].Values["site"].(string)
		m := make(map[string]string)
		if err := json.Unmarshal([]byte(currMessage), &m); err == nil {
			db := helpers.ConnectDB()
			sqlDB, Serr := db.DB()
			if Serr != nil {
				log.Fatal(err)
			}
			defer sqlDB.Close()
			mails := helpers.GetEmails(&db, m["siteId"])
			fmt.Println(mails)

			msg := fmt.Sprintf("From: "+os.Getenv("FromEmail")+"\r\n"+
				"To: "+"\r\n"+
				"Subject: Website Down Alert\r\n"+
				"MIME-Version: 1.0\r\n"+
				"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n you webiste with sidtId: %s was down in the regionId %s", m["siteId"],m["regionId"]) +
				" "
			err := SendMain(mails, msg)
			if err != nil {
				fmt.Print("failed to send eamil")
				fmt.Print(err)
			}
			fmt.Println("sent")
			client.XAck(context.Background(), "notifications", readRes[0].ID)
		}
	}
}
