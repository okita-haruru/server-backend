package service

import (
	"encoding/json"
	"sushi/model"
)

func (service *Service) decodeFishData(data model.CustomfishingData) (*model.CustomfishingDataDecoded, error) {
	var playerfishdata PlayerFishData
	err := json.Unmarshal([]byte(data.Data), &playerfishdata)
	if err != nil {
		return nil, err
	}
	decodedRecord := model.CustomfishingDataDecoded{
		Uuid:    data.Uuid,
		Amount:  playerfishdata.Stats.Amount,
		MaxSize: playerfishdata.Stats.Size,
	}
	return &decodedRecord, nil
}
func (service *Service) getFishData() []model.CustomfishingData {
	var record []model.CustomfishingData
	service.db.DB.Model(&record).Find(&record)
	return record
}
func (service *Service) getDecodedFishData() ([]model.CustomfishingDataDecoded, error) {
	record := service.getFishData()
	var decodedRecords []model.CustomfishingDataDecoded
	for _, data := range record {
		decodedRecord, err := service.decodeFishData(data)
		if err != nil {
			return nil, err
		}
		decodedRecords = append(decodedRecords, *decodedRecord)
	}
	return decodedRecords, nil
}
