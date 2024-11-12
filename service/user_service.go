package service

import "sushi/model"

func (service *Service) GetPlayerWardenKillRanking(uuid string) model.RankingInt {
	var record model.PlayerKillStats
	service.db.DB.Model(&record).Where("player_id = ?", uuid).First(&record)
	killCount := record.WardenKills
	var count int64
	service.db.DB.Model(&record).Where("warden_kills > ?", killCount).Count(&count)
	return model.RankingInt{
		Rank:  int(count) + 1,
		Value: killCount,
	}
}
func (service *Service) GetPlayerEnderDragonKillRanking(uuid string) model.RankingInt {
	var record model.PlayerKillStats
	service.db.DB.Model(&record).Where("player_id = ?", uuid).First(&record)
	killCount := record.EnderDragonKills
	var count int64
	service.db.DB.Model(&record).Where("ender_dragon_kills > ?", killCount).Count(&count)
	return model.RankingInt{
		Rank:  int(count) + 1,
		Value: killCount,
	}
}
func (service *Service) GetPlayerWitherKillRanking(uuid string) model.RankingInt {
	var record model.PlayerKillStats
	service.db.DB.Model(&record).Where("player_id = ?", uuid).First(&record)
	killCount := record.WitherKills
	var count int64
	service.db.DB.Model(&record).Where("wither_kills > ?", killCount).Count(&count)
	return model.RankingInt{
		Rank:  int(count) + 1,
		Value: killCount,
	}
}
func (service *Service) GetPlayerPiglinBruteKillRanking(uuid string) model.RankingInt {
	var record model.PlayerKillStats
	service.db.DB.Model(&record).Where("player_id = ?", uuid).First(&record)
	killCount := record.PiglinBruteKills
	var count int64
	service.db.DB.Model(&record).Where("piglin_brute_kills > ?", killCount).Count(&count)
	return model.RankingInt{
		Rank:  int(count) + 1,
		Value: killCount,
	}
}
func (service *Service) GetPlayerPhantomKillRanking(uuid string) model.RankingInt {
	var record model.PlayerKillStats
	service.db.DB.Model(&record).Where("player_id = ?", uuid).First(&record)
	killCount := record.PhantomKills
	var count int64
	service.db.DB.Model(&record).Where("phantom_kills > ?", killCount).Count(&count)
	return model.RankingInt{
		Rank:  int(count) + 1,
		Value: killCount,
	}
}
func (service *Service) GetPlayerAncientGuardianKillRanking(uuid string) model.RankingInt {
	var record model.PlayerKillStats
	service.db.DB.Model(&record).Where("player_id = ?", uuid).First(&record)
	killCount := record.AncientGuardianKills
	var count int64
	service.db.DB.Model(&record).Where("ancient_guardian_kills > ?", killCount).Count(&count)
	return model.RankingInt{
		Rank:  int(count) + 1,
		Value: killCount,
	}

}
