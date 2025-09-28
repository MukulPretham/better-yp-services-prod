package helpers

import "time"

type User struct{
	Id string `gorm:"column:id"`
	Name string `gorm:"column:name"`
	Password string `gorm:"column:password"`
	Email string `gorm:"column:email"`
}

type Status struct {
	Id       string `gorm:"column:id"`
	SiteId   string `gorm:"column:siteId"`
	RegionId string `gorm:"column:regionId"`
	Status   bool   `gorm:"column:status"`
}

type Latency struct {
	Id       string   `gorm:"column:id"`
	SiteId   string `gorm:"column:siteId"`
	RegionId string `gorm:"column:regionId"`
	Latency  float64 `gorm:"column:latency`
	Time     time.Time  `gorm:"column:time`
}

type Region struct {
	Id   string
	Name string
}

type Website struct {
	Id   string
	Name string
	Url  string
}

type StreamMsg struct {
	Id   string
	Name string
	Url  string
}

type UserToWebsite struct {
    Id     string `gorm:"column:id"`
    UserId string `gorm:"column:userId"`
    SiteId string `gorm:"column:siteId"`
}

func (Region) TableName() string {
	return "Region"
}

func (Status) TableName() string {
	return "Status"
}

func (Website) TableName() string {
	return "Website"
}

func (Latency) TableName() string {
	return "Latency"
}

func (UserToWebsite) TableName() string {
	return "UserToWebsite"
}

func (User) TableName() string {
	return "User"
}

