package model

import "time"

type PlayerKillStats struct {
	PlayerId             string `gorm:"primaryKey"`
	AncientGuardianKills int    `gorm:"default:0"`
	PhantomKills         int    `gorm:"default:0"`
	PiglinKills          int    `gorm:"default:0"`
	EnderDragonKills     int    `gorm:"default:0"`
	WitherKills          int    `gorm:"default:0"`
	WardenKills          int    `gorm:"default:0"`
}
type Xconomy struct { // Xconomy is a table in the database
	UID     string `gorm:"primaryKey"`
	Player  string
	Balance float32 `gorm:"default:0"`
	Hidden  int     `gorm:"default:0"`
}

func (Xconomy) TableName() string {
	return "xconomy"
}

type CustomfishingData struct {
	Uuid string `gorm:"primaryKey"`
	Lock int
	Data string
}
type CustomfishingDataDecoded struct {
	Uuid    string `gorm:"primaryKey"`
	Amount  map[string]int
	MaxSize map[string]float32
}
type FishAmountData struct {
	Uuid   string
	Amount map[string]int
}
type FishSizeData struct {
	Uuid    string
	MaxSize map[string]float32
}
type PlayerProfile struct {
	UUID      string `gorm:"primaryKey"`
	Name      string
	FishData  CustomfishingDataDecoded
	KillStats PlayerKillStats
	Xconomy   Xconomy
}
type PlayerInfo struct {
	Uuid   string `json:"uuid"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}
type PlayTime struct {
	UUID       string    `gorm:"primaryKey" json:"uuid"`
	Name       string    `json:"name"`
	FirstLogin time.Time `json:"first_login"`
	PlayTime   int       `json:"play_time"`
	LastLogin  time.Time `json:"last_login"`
}

func (PlayTime) TableName() string {
	return "play_time"
}
