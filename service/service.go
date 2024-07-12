package service

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sushi/model"
	"sushi/utils/DB"
	"sushi/utils/config"
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
	service.db.DB.Model(&record).Where("uid = ?", uuid).Select("balance")
	return record.Balance
}
func (service *Service) GetBalanceRanking(page int) []model.Xconomy {
	var record []model.Xconomy
	service.db.DB.Model(&record).Order("balance desc").Offset((page - 1) * 20).Limit(20).Find(&record)
	return record
}
func (service *Service) GetPlayerKillStatsSortByTotal(page int) []model.PlayerKillStats {
	var record []model.PlayerKillStats
	service.db.DB.Model(&record).Order(gorm.Expr("warden_kills + ender_dragon_kills + wither_kills + piglin_kills + phantom_kills + ancient_guardian_kills DESC")).Offset((page - 1) * 20).Limit(20).Find(&record)
	return record
}
func (service *Service) GetPlayerKillStatsSortByWarden(page int) []model.PlayerKillStats {
	var record []model.PlayerKillStats
	service.db.DB.Model(&record).Order("warden_kills desc").Offset((page - 1) * 20).Limit(20).Find(&record)
	return record
}
func (service *Service) GetPlayerKillStatsSortByEnderDragon(page int) []model.PlayerKillStats {
	var record []model.PlayerKillStats
	service.db.DB.Model(&record).Order("ender_dragon_kills desc").Offset((page - 1) * 20).Limit(20).Find(&record)
	return record
}
func (service *Service) GetPlayerKillStatsSortByWither(page int) []model.PlayerKillStats {
	var record []model.PlayerKillStats
	service.db.DB.Model(&record).Order("wither_kills desc").Offset((page - 1) * 20).Limit(20).Find(&record)
	return record
}
func (service *Service) GetPlayerKillStatsSortByPiglin(page int) []model.PlayerKillStats {
	var record []model.PlayerKillStats
	service.db.DB.Model(&record).Order("piglin_kills desc").Offset((page - 1) * 20).Limit(20).Find(&record)
	return record
}
func (service *Service) GetPlayerKillStatsSortByPhantom(page int) []model.PlayerKillStats {
	var record []model.PlayerKillStats
	service.db.DB.Model(&record).Order("phantom_kills desc").Offset((page - 1) * 20).Limit(20).Find(&record)
	return record
}
func (service *Service) GetPlayerKillStatsSortByAncientGuardian(page int) []model.PlayerKillStats {
	var record []model.PlayerKillStats
	service.db.DB.Model(&record).Order("ancient_guardian_kills desc").Offset((page - 1) * 20).Limit(20).Find(&record)
	return record
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

func (service *Service) decodeFishData(data model.CustomfishingData) (error, *model.CustomfishingDataDecoded) {
	var record model.CustomfishingData
	service.db.DB.Model(&record).Where("uuid = ?", data.Uuid).First(&record)
	var playerfishdata PlayerFishData
	err := json.Unmarshal([]byte(record.Data), &playerfishdata)
	if err != nil {
		return err, nil
	}
	decodedRecord := model.CustomfishingDataDecoded{
		Uuid:    record.Uuid,
		Amount:  playerfishdata.Stats.Amount,
		MaxSize: playerfishdata.Stats.Size,
	}
	return nil, &decodedRecord
}
func (service *Service) getFishData() []model.CustomfishingData {
	var record []model.CustomfishingData
	service.db.DB.Model(&record).Find(&record)
	return record
}
func (service *Service) GetDecodedFishData() (error, []model.CustomfishingDataDecoded) {
	record := service.getFishData()
	var decodedRecords []model.CustomfishingDataDecoded
	for _, data := range record {
		err, decodedRecord := service.decodeFishData(data)
		if err != nil {
			return err, nil
		}
		decodedRecords = append(decodedRecords, *decodedRecord)
	}
	return nil, decodedRecords
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

// GetPlayerProfileByUUID
func (service *Service) GetPlayerProfileByName(name string) model.PlayerProfile {
	var record model.PlayerProfile
	uuid := service.GetUUIDByName(name)
	service.db.DB.Model(&model.Xconomy{}).Where("uid = ?", uuid).First(&record.Xconomy)
	service.db.DB.Model(&model.PlayerKillStats{}).Where("player_id = ?", uuid).First(&record.KillStats)
	var FishData model.CustomfishingData
	service.db.DB.Model(&model.CustomfishingData{}).Where("uuid = ?", uuid).First(&FishData)
	err, decodedFishData := service.decodeFishData(FishData)
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
