package service

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sushi/model"
	"sushi/utils/DB"
	"sushi/utils/config"
	"time"
)

type Service struct {
	db   *DB.DB
	log  *logrus.Logger
	conf *config.Config
	Ctx  *context.Context
}

func NewService(db *DB.DB, log *logrus.Logger, conf *config.Config, ctx context.Context) *Service {
	return &Service{
		db:   db,
		log:  log,
		conf: conf,
		Ctx:  &ctx,
	}
}
func (service *Service) GetBalance(uuid string) float32 {
	var record model.Xconomy
	service.db.DB.Model(&record).Where("uid = ?", uuid).First(&record)
	return record.Balance
}

func (service *Service) GetPlayerKillStatsByUUID(uuid string) model.PlayerKillStats {
	var record model.PlayerKillStats
	service.db.DB.Model(&record).Where("player_id = ?", uuid).First(&record)
	return record
}
func (service *Service) GetPlayerKillStats() []model.PlayerKillStats {
	var record []model.PlayerKillStats
	service.db.DB.Model(&record).Find(&record)
	return record
}
func (service *Service) GetPlayTime(page int) []model.PlayTime {
	var record []model.PlayTime
	service.db.DB.Model(&record).Order("play_time desc").Offset((page - 1) * 20).Limit(20).Find(&record)
	return record
}
func (service *Service) GetPlayTimeByUUID(uuid string) model.PlayTime {
	var record model.PlayTime
	service.db.DB.Model(&record).Where("uuid = ?", uuid).First(&record)
	return record
}

func (service *Service) GetPlayerInfo(username string) (*model.PlayerInfo, error) {
	uuid := service.GetUUIDByName(username)
	playTime := service.GetPlayTimeByUUID(uuid)
	fishingRanking, err := service.GetFishRanking(uuid)
	amountRanking, err := service.GetAmountRanking(uuid)
	if err != nil {
		return nil, err
	}
	playerInfo := model.PlayerInfo{
		Uuid:   uuid,
		Name:   username,
		Avatar: service.GetAvatar(username),
		States: model.States{
			TotalFishing: *amountRanking,
			Fishing:      fishingRanking,
			Balance:      service.GetPlayerBalanceRanking(uuid),
			Death:        service.GetPlayerDeathRanking(uuid),
			Kills: model.KillRanking{
				Warden:          service.GetPlayerWardenKillRanking(uuid),
				AncientGuardian: service.GetPlayerAncientGuardianKillRanking(uuid),
				EnderDragon:     service.GetPlayerEnderDragonKillRanking(uuid),
				Wither:          service.GetPlayerWitherKillRanking(uuid),
				PiglinBrute:     service.GetPlayerPiglinBruteKillRanking(uuid),
				Phantom:         service.GetPlayerPhantomKillRanking(uuid),
			},
		},
		Join:     playTime.FirstLogin.Unix(),
		LastSeen: playTime.LastLogin.Unix(),
		PlayTime: int64(playTime.PlayTime),
		IsOnline: service.IsPlayerOnline(uuid),
	}
	return &playerInfo, nil
}

type PlayerFishData struct {
	Name  string `json:"name"`
	Stats Stats  `json:"stats"`
	Trade Trade  `json:"trade"`
}
type Stats struct {
	Amount map[string]int     `json:"amount"`
	Size   map[string]float32 `json:"size"`
	Bag    Bag                `json:"bag"`
}
type Bag struct {
	Inventory string `json:"inventory"`
	Size      int    `json:"size"`
}
type Trade struct {
	Earnings float32 `json:"earnings"`
	Data     int     `json:"data"`
}

func (service *Service) GetNameByUUID(uuid string) string {
	var record model.Xconomy
	service.db.DB.Model(&record).Where("uid = ?", uuid).First(&record)
	return record.Player
}
func (service *Service) GetUUIDByName(name string) string {
	var record model.Xconomy
	service.db.DB.Model(&record).Where("player = ?", name).First(&record)
	return record.UID
}

// GetPlayers
func (service *Service) GetPlayers() []model.Player {
	var records []model.Xconomy
	var result []model.Player
	service.db.DB.Model(&records).Find(&records)
	for _, record := range records {
		result = append(result, model.Player{
			UUID:   record.UID,
			Name:   record.Player,
			Avatar: service.GetAvatar(record.Player),
		})
	}
	return result
}

func (service *Service) GetPlayerProfileByName(name string) model.PlayerProfile {
	var record model.PlayerProfile
	uuid := service.GetUUIDByName(name)
	service.db.DB.Model(&model.Xconomy{}).Where("uid = ?", uuid).First(&record.Xconomy)
	service.db.DB.Model(&model.PlayerKillStats{}).Where("player_id = ?", uuid).First(&record.KillStats)
	var FishData model.CustomfishingData
	service.db.DB.Model(&model.CustomfishingData{}).Where("uuid = ?", uuid).First(&FishData)
	decodedFishData, err := service.decodeFishData(FishData)
	if err != nil {
		return model.PlayerProfile{}
	}
	record.FishData = *decodedFishData
	record.UUID = uuid
	return record
}

func (service *Service) GetAvatar(name string) string {
	return "https://minotar.net/helm/" + name + "/100.png"
}

type Fish struct {
	Name         string `json:"name"`
	FormalName   string `json:"formal_name"`
	Key          string `json:"key"`
	LatinName    string `json:"latin_name"`
	MinSize      int    `json:"min_size"`
	MaxSize      int    `json:"max_size"`
	Price        int    `json:"price"`
	Rarity       string `json:"rarity"`
	Description  string `json:"description"`
	Distribution string `json:"distribution"`
}

func (service *Service) GetPlayerBalanceRanking(uuid string) model.RankingFloat {
	//get ranking of one player by balance
	var count int64
	service.db.DB.Model(&model.Xconomy{}).Where("balance > (select balance from xconomy where uid = ?)", uuid).Count(&count)
	return model.RankingFloat{
		Rank:  int(count) + 1,
		Value: service.GetBalance(uuid),
	}
}
func (service *Service) GetPlayerDeathRanking(uuid string) model.RankingInt {
	//get ranking of one player by death
	var count int64
	death := service.GetDeathCount(uuid)
	service.db.DB.Model(&model.PlayerDeathStats{}).Where("death_count > ?", death).Count(&count)
	return model.RankingInt{
		Rank:  int(count) + 1,
		Value: death,
	}
}

type LoginRecordCount struct {
	LoginTime time.Time `json:"login_time"`
	Count     int       `json:"count"`
}
type LoginRecordRes struct {
	LoginTime string `json:"login_time"`
	Count     int    `json:"count"`
}

func (service *Service) GetLoginRecordCountByDate() []LoginRecordRes {
	var record []model.LoginRecord
	var recordCount []LoginRecordCount
	var recordRes []LoginRecordRes
	service.db.DB.Model(&record).Select("login_time, count(*) as count").Group("login_time").Order("login_time").Find(&recordCount)

	recordCountMap := make(map[string]int)
	for _, v := range recordCount {
		recordCountMap[v.LoginTime.Format("2006-01-02")] = v.Count
	}

	startDate := recordCount[0].LoginTime
	endDate := recordCount[len(recordCount)-1].LoginTime

	for d := startDate; d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		count, exists := recordCountMap[dateStr]
		if !exists {
			count = 0
		}
		recordRes = append(recordRes, LoginRecordRes{
			LoginTime: dateStr,
			Count:     count,
		})
	}

	return recordRes
}
func (service *Service) GetDeathRanking(pageInt int) []model.PlayerDeathStats {
	var record []model.PlayerDeathStats
	service.db.DB.Model(&record).Order("death_count desc").Offset((pageInt - 1) * 20).Limit(20).Find(&record)
	return record
}
func (service *Service) GetDeathCount(uuid string) int {
	var record model.PlayerDeathStats
	service.db.DB.Model(&record).Where("player_id = ?", uuid).First(&record)
	return record.DeathCount
}

func (service *Service) getOnlinePlayerList() ([]model.PlayerDetail, error) {
	resp, err := http.Get("http://localhost:25577/api/players")
	if err != nil {
		return nil, err
	}
	var response model.PlayerListResponse
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	var result []model.PlayerDetail
	for _, player := range response.Lobby.Players {
		result = append(result, player)
	}
	for _, player := range response.Survival.Players {
		result = append(result, player)
	}
	return result, nil
}
func (service *Service) IsPlayerOnline(uuid string) bool {
	players, err := service.getOnlinePlayerList()
	if err != nil {
		return false
	}
	for _, player := range players {
		if player.UUID == uuid {
			return true
		}
	}
	return false
}
