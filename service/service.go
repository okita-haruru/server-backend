package service

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"sort"
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
	//service.db.DB.Model(&record).Order(gorm.Expr("warden_kills + ender_dragon_kills + wither_kills + piglin_kills + phantom_kills + ancient_guardian_kills desc")).Offset((page - 1) * 20).Limit(20).Find(&record)
	service.db.DB.Model(&record).
		Select("*, (warden_kills + ender_dragon_kills + wither_kills + piglin_kills + phantom_kills + ancient_guardian_kills) as total_kills").
		Order("total_kills desc").
		Offset((page - 1) * 20).
		Limit(20).
		Find(&record)
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
func (service *Service) GetPlayTime(page int) []model.PlayTime {
	var record []model.PlayTime
	service.db.DB.Model(&record).Order("play_time desc").Offset((page - 1) * 20).Limit(20).Find(&record)
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
	var playerfishdata PlayerFishData
	err := json.Unmarshal([]byte(data.Data), &playerfishdata)
	if err != nil {
		return err, nil
	}
	decodedRecord := model.CustomfishingDataDecoded{
		Uuid:    data.Uuid,
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
func (service *Service) getDecodedFishData() (error, []model.CustomfishingDataDecoded) {
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

func (service *Service) sortFishDataByAmount(fish string) (error, []model.CustomfishingDataDecoded) {
	err, records := service.getDecodedFishData()
	if err != nil {
		return err, nil
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Amount[fish] > records[j].Amount[fish]
	})
	return nil, records
}

func (service *Service) GetTotalAmount(record model.CustomfishingDataDecoded) int {
	var totalAmount int
	for _, amount := range record.Amount {
		totalAmount += amount
	}
	return totalAmount
}

func (service *Service) sortFishDataByTotalAmount() (error, []model.CustomfishingDataDecoded) {
	err, records := service.getDecodedFishData()
	if err != nil {
		return err, nil
	}
	sort.Slice(records, func(i, j int) bool {
		return service.GetTotalAmount(records[i]) > service.GetTotalAmount(records[j])
	})
	return nil, records
}

func (service *Service) sortFishDataBySize(fish string) (error, []model.CustomfishingDataDecoded) {
	err, records := service.getDecodedFishData()
	if err != nil {
		return err, nil
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].MaxSize[fish] > records[j].MaxSize[fish]
	})
	return nil, records

}
func (service *Service) GetFishRankingByAmount(fish string, page int) (error, []model.CustomfishingDataDecoded) {
	err, records := service.sortFishDataByAmount(fish)
	if err != nil {
		return err, nil
	}
	return nil, records[(page-1)*20 : min(20*page-1, len(records))]
}
func (service *Service) GetFishRankingByTotalAmount(page int) (error, []model.CustomfishingDataDecoded) {
	err, records := service.sortFishDataByTotalAmount()
	if err != nil {
		return err, nil
	}
	return nil, records[(page-1)*20 : min(20*page-1, len(records))]
}
func (service *Service) GetFishRankingBySize(fish string, page int) (error, []model.CustomfishingDataDecoded) {
	err, records := service.sortFishDataBySize(fish)
	if err != nil {
		return err, nil
	}
	return nil, records[(page-1)*20 : min(20*page-1, len(records))]
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
