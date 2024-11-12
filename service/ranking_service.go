package service

import (
	"sort"
	"sushi/model"
)

func (service *Service) GetBalanceRanking(page int) []model.Xconomy {
	var record []model.Xconomy
	service.db.DB.Model(&record).Order("balance desc").Offset((page - 1) * 20).Limit(20).Find(&record)
	return record
}
func (service *Service) GetPlayerKillStatsSortByTotal(page int) []model.PlayerKillStats {
	var record []model.PlayerKillStats
	//service.db.DB.Model(&record).Order(gorm.Expr("warden_kills + ender_dragon_kills + wither_kills + piglin_kills + phantom_kills + ancient_guardian_kills desc")).Offset((page - 1) * 20).Limit(20).Find(&record)
	service.db.DB.Model(&record).
		Select("*, (warden_kills + ender_dragon_kills + wither_kills + piglin_brute_kills + phantom_kills + ancient_guardian_kills) as total_kills").
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
	service.db.DB.Model(&record).Order("piglin_brute_kills desc").Offset((page - 1) * 20).Limit(20).Find(&record)
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
func (service *Service) sortFishDataByAmount(fish string) ([]model.CustomfishingDataDecoded, error) {
	records, err := service.getDecodedFishData()
	if err != nil {
		return nil, err
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Amount[fish] > records[j].Amount[fish]
	})
	return records, nil
}
