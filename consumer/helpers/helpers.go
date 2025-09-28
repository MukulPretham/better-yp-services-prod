package helpers

import (
	// "fmt"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	
)

// Connect to database and get all websites
func ConnectDB() gorm.DB {
	const dsn = "host=better-up-postgres user=postgres password=9059015626 dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal("unable to connect ot the database")
	}
	return *db
}

// func getRegions(db *gorm.DB) []string {
// 	var regions []Region
// 	db.Find(&regions)
// 	var regionList []string
// 	for _, region := range regions {
// 		regionList = append(regionList, region.Id)
// 	}
// 	return regionList
// }

func getRegionId(db *gorm.DB, regionName string) (string, error) {
	var currRegion Region
	err := db.Where("name = ?", regionName).First(&currRegion)
	if err.Error == nil {
		return currRegion.Id, nil
	} else {
		return "", errors.New("region wich was passed via env variable is not a valid region")
	}
}

func getSiteId(db *gorm.DB, url string) string {
	var website Website
	db.Where("url = ?", url).First(&website)
	return website.Id
}

func setStatus(db *gorm.DB, siteId string, regionId string, status bool) bool {
	err := db.Model(&Status{}).Where(`"siteId" = ? AND "regionId" = ?`, siteId, regionId).Update("status", status)
	if err.Error != nil {
		return false
	}
	return true
}

func GetStatus(db *gorm.DB, siteId string, regionId string) bool {
	var currSite Status
	err := db.Model(&Status{}).Where(`"siteId" = ? AND "regionId" = ?`, siteId, regionId).First(&currSite)
	if err.Error != nil {
		return false
	}
	return currSite.Status
}

func setLatency(db *gorm.DB, siteId string, regionId string, latency float64) {
	latencyRepot := Latency{
		Id:       uuid.NewString(),
		SiteId:   siteId,
		RegionId: regionId,
		Latency:  latency,
		Time:     time.Now(),
	}
	db.Create(&latencyRepot)
}

func GetEmails(db *gorm.DB,siteId string)[]string{
	var currUserIds []UserToWebsite 
	db.Find(&currUserIds,`"siteId"= ?`,siteId)
	
	var mails []string
	for _,userId := range currUserIds{
		
		var currUser User
		db.First(&currUser, "id = ?", userId.UserId)
		mails = append(mails, currUser.Email)
	}
	return mails
}

func fetch(url string) int {
	client := &http.Client{
		Timeout: 15 * time.Second, // set a timeout
	}

	res, err := client.Get(fmt.Sprintf("https://%s", url))
	if err != nil {
		fmt.Println("Request error:", err)
		return 0
	}
	defer res.Body.Close() // always close body

	if res.StatusCode == 200 {
		return 200
	}
	return 0
}

func WriteToDB(url string, client *redis.Client, ID string) {
	var currLatency float64

	start := time.Now()

	db := ConnectDB()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()
	res := fetch(url)

	currLatency = float64(time.Since(start).Milliseconds())

	env := os.Getenv("REGION")
	currRegionId, err := getRegionId(&db, env)
	if err != nil {
		log.Fatal(err)
	}

	currSiteId := getSiteId(&db, url)

	if res == 200 {
		setStatus(&db, currSiteId, currRegionId, true)
		setLatency(&db, currSiteId, currRegionId, currLatency)
	} else {
		prevState := GetStatus(&db, currSiteId, currRegionId)
		setLatency(&db, currSiteId, currRegionId, 0)
		setStatus(&db, currSiteId, currRegionId, false)
		// Adding this failed siteId to notifications queue
		
		if prevState == true {
			data, err := json.Marshal(map[string]string{"siteId": currSiteId,"regionId":currRegionId})
			if err != nil {
				fmt.Println("error while writing to redis stream : &%v ", err)
			}
			client.XAdd(context.Background(), &redis.XAddArgs{
				Stream: "notifications",
				Values: map[string]any{
					"site": string(data),
				},
			})
		}
	}
	fmt.Println("updated", url)
	ctx := context.Background()
	client.XAck(ctx, "websites", "consumerGroup", ID)
}


