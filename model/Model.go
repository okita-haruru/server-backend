package model

import (
	"time"
)

type PlayerKillStats struct {
	PlayerId             string `gorm:"primaryKey"`
	AncientGuardianKills int    `gorm:"default:0"`
	PhantomKills         int    `gorm:"default:0"`
	PiglinBruteKills     int    `gorm:"default:0"`
	EnderDragonKills     int    `gorm:"default:0"`
	WitherKills          int    `gorm:"default:0"`
	WardenKills          int    `gorm:"default:0"`
}
type PlayerDeathStats struct {
	PlayerId   string `gorm:"primaryKey"`
	DeathCount int    `gorm:"default:0"`
}
type Xconomy struct { // Xconomy is a table in the database
	UID     string `gorm:"primaryKey"`
	Player  string
	Balance float32 `gorm:"default:0"`
	Hidden  int     `gorm:"default:0"`
}
type Player struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
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
	Uuid     string `json:"uuid"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	States   States `json:"states"`
	Join     int64  `json:"join"`
	LastSeen int64  `json:"lastSeen"`
	PlayTime int64  `json:"playTime"`
	IsOnline bool   `json:"isOnline"`
}
type RankingInt struct {
	Rank  int `json:"rank"`
	Value int `json:"value"`
}
type RankingFloat struct {
	Rank  int     `json:"rank"`
	Value float32 `json:"value"`
}
type States struct {
	TotalFishing RankingInt       `json:"totalFishing"`
	Fishing      []FishingRanking `json:"fishing"`
	Balance      RankingFloat     `json:"balance"`
	Death        RankingInt       `json:"death"`
	Kills        KillRanking      `json:"kills"`
}
type KillRanking struct {
	Warden          RankingInt `json:"warden"`
	AncientGuardian RankingInt `json:"ancientGuardian"`
	EnderDragon     RankingInt `json:"enderDragon"`
	Wither          RankingInt `json:"wither"`
	PiglinBrute     RankingInt `json:"piglinBrute"`
	Phantom         RankingInt `json:"phantom"`
}
type FishingRanking struct {
	Name    string       `json:"name"`
	Amount  RankingInt   `json:"amount"`
	MaxSize RankingFloat `json:"maxSize"`
}
type PlayTime struct {
	UUID       string    `gorm:"primaryKey" json:"uuid"`
	Name       string    `json:"name"`
	FirstLogin time.Time `json:"first_login"`
	PlayTime   int       `json:"play_time"`
	LastLogin  time.Time `json:"last_login"`
}

type LoginRecord struct {
	UUID      string `gorm:"primaryKey"`
	Name      string
	LoginTime time.Time `gorm:"type:date"`
}
type UserInfo struct {
	UUID string `gorm:"primaryKey"`
}

func (LoginRecord) TableName() string {
	return "login_record"
}

func (PlayTime) TableName() string {
	return "play_time"
}

type PlayerDetail struct {
	Ping   int    `json:"ping"`
	Name   string `json:"name"`
	UUID   string `json:"uuid"`
	Avatar string `json:"avatar"`
}
type PlayerListResponse struct {
	Lobby    RoomJson `json:"lobby"`
	Survival RoomJson `json:"survival"`
}
type RoomJson struct {
	Players []PlayerDetail `json:"players"`
	Count   int            `json:"count"`
}
