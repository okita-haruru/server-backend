package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sort"
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
	var newRecords []model.CustomfishingDataDecoded
	for _, record := range records {
		if record.MaxSize[fish] != 0 {
			newRecords = append(newRecords, record)
		}
	}
	return nil, newRecords[(page-1)*20 : min(20*page-1, len(newRecords))]
}
func (service *Service) GetFishRankingByTotalAmount(page int) (error, []model.CustomfishingDataDecoded) {
	err, records := service.sortFishDataByTotalAmount()
	if err != nil {
		return err, nil
	}
	for i, record := range records {
		fmt.Println(record)
		fmt.Println(service.GetTotalAmount(record))
		if service.GetTotalAmount(record) == 0 {
			records = records[:i]
			break
		}
	}
	fmt.Println(len(records))
	return nil, records[(page-1)*20 : min(20*page-1, len(records))]
}
func (service *Service) GetFishRankingBySize(fish string, page int) (error, []model.CustomfishingDataDecoded) {
	err, records := service.sortFishDataBySize(fish)
	if err != nil {
		return err, nil
	}
	var newRecords []model.CustomfishingDataDecoded
	for _, record := range records {
		if record.MaxSize[fish] != 0 {
			newRecords = append(newRecords, record)
		}
	}
	return nil, newRecords[(page-1)*20 : min(20*page-1, len(newRecords))]
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

type Fish struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func (service *Service) GetFish() ([]Fish, error) {
	var fishes []Fish

	file, err := os.ReadFile("new.csv")
	if err != nil {
		fmt.Println(err)
	}
	//解决读取csv中文乱码的问题
	//reader := csv.NewReader(transform.NewReader(bytes.NewReader(file), simplifiedchinese.GBK.NewDecoder()))

	reader := csv.NewReader(bytes.NewReader(file))
	reader.Read()
	for {
		csvdata, err := reader.Read() // 按行读取数据,可控制读取部分
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error:", err.Error())
			return nil, err
		}

		if csvdata[0] != "" && csvdata[3] != "" {
			fishes = append(fishes, Fish{
				Name: csvdata[0],
				Key:  csvdata[3],
			})
		}
	}
	return fishes, nil
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

	// 创建一个map，以便于快速查找日期
	recordCountMap := make(map[string]int)
	for _, v := range recordCount {
		recordCountMap[v.LoginTime.Format("2006-01-02")] = v.Count
	}

	// 获取日期范围
	startDate := recordCount[0].LoginTime
	endDate := recordCount[len(recordCount)-1].LoginTime

	// 遍历日期范围
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
