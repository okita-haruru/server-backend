package service

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sushi/model"
)

func (service *Service) GetTotalAmount(record model.CustomfishingDataDecoded) int {
	var totalAmount int
	for _, amount := range record.Amount {
		totalAmount += amount
	}
	return totalAmount
}

func (service *Service) sortFishDataByTotalAmount() ([]model.CustomfishingDataDecoded, error) {
	records, err := service.getDecodedFishData()
	if err != nil {
		return nil, err
	}
	sort.Slice(records, func(i, j int) bool {
		return service.GetTotalAmount(records[i]) > service.GetTotalAmount(records[j])
	})
	return records, nil
}

func (service *Service) sortFishDataBySize(fish string) ([]model.CustomfishingDataDecoded, error) {
	records, err := service.getDecodedFishData()
	if err != nil {
		return nil, err
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].MaxSize[fish] > records[j].MaxSize[fish]
	})
	return records, nil
}
func (service *Service) GetFishRankingByAmount(fish string, page int) ([]model.CustomfishingDataDecoded, error) {
	records, err := service.sortFishDataByAmount(fish)
	if err != nil {
		return nil, err
	}
	var newRecords []model.CustomfishingDataDecoded
	for _, record := range records {
		if record.MaxSize[fish] != 0 {
			newRecords = append(newRecords, record)
		}
	}
	return newRecords[(page-1)*20 : min(20*page-1, len(newRecords))], nil
}
func (service *Service) GetFishRankingByAmountNoPage(fish string) ([]model.CustomfishingDataDecoded, error) {
	records, err := service.sortFishDataByAmount(fish)
	if err != nil {
		return nil, err
	}
	var newRecords []model.CustomfishingDataDecoded
	for _, record := range records {
		if record.MaxSize[fish] != 0 {
			newRecords = append(newRecords, record)
		}
	}
	return newRecords, nil
}
func (service *Service) GetFishRankingByTotalAmount(page int) ([]model.CustomfishingDataDecoded, error) {
	records, err := service.sortFishDataByTotalAmount()
	if err != nil {
		return nil, err
	}
	for i, record := range records {
		if service.GetTotalAmount(record) == 0 {
			records = records[:i]
			break
		}
	}
	return records[(page-1)*20 : min(20*page-1, len(records))], nil
}

func (service *Service) GetFishRankingByTotalAmountNoPage() ([]model.CustomfishingDataDecoded, error) {
	records, err := service.sortFishDataByTotalAmount()
	if err != nil {
		return nil, err
	}
	for i, record := range records {
		if service.GetTotalAmount(record) == 0 {
			records = records[:i]
			break
		}
	}
	return records, nil
}
func (service *Service) GetFishRankingBySize(fish string, page int) ([]model.CustomfishingDataDecoded, error) {
	records, err := service.sortFishDataBySize(fish)
	if err != nil {
		return nil, err
	}
	var newRecords []model.CustomfishingDataDecoded
	for _, record := range records {
		if record.MaxSize[fish] != 0 {
			newRecords = append(newRecords, record)
		}
	}
	return newRecords[(page-1)*20 : min(20*page-1, len(newRecords))], nil
}
func (service *Service) GetFishRankingBySizeNoPage(fish string) ([]model.CustomfishingDataDecoded, error) {
	records, err := service.sortFishDataBySize(fish)
	if err != nil {
		return nil, err
	}
	var newRecords []model.CustomfishingDataDecoded
	for _, record := range records {
		if record.MaxSize[fish] != 0 {
			newRecords = append(newRecords, record)
		}
	}
	return newRecords, nil
}

func (service *Service) GetFish() ([]Fish, error) {
	var fishes []Fish

	file, err := os.ReadFile("new.csv")
	if err != nil {
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
			return nil, err
		}
		var minSize, maxSize int
		if csvdata[5] != "" {
			parts := strings.Split(csvdata[5], "~")
			minSize, _ = strconv.Atoi(parts[0])
			maxSize, _ = strconv.Atoi(parts[1])
		} else {
			minSize = 0
			maxSize = 0
		}

		price, _ := strconv.Atoi(csvdata[7])
		if csvdata[0] != "" && csvdata[3] != "" {
			fishes = append(fishes, Fish{
				Name:         csvdata[0],
				Key:          csvdata[3],
				FormalName:   csvdata[1],
				LatinName:    csvdata[2],
				MinSize:      minSize,
				MaxSize:      maxSize,
				Price:        price,
				Rarity:       csvdata[10],
				Description:  csvdata[14],
				Distribution: csvdata[19],
			})
		}
	}
	return fishes, nil
}
func (service *Service) getFishDataByUUID(uuid string) model.CustomfishingData {
	var record model.CustomfishingData
	service.db.DB.Model(&record).Where("uuid = ?", uuid).First(&record)
	return record
}
func (service *Service) GetFishRanking(uuid string) ([]model.FishingRanking, error) {
	var result []model.FishingRanking
	decodedRecord, err := service.decodeFishData(service.getFishDataByUUID(uuid))
	if err != nil {
		return nil, err
	}
	fishes, err := service.GetFish()
	if err != nil {
		return nil, err
	}
	for _, fish := range fishes {
		amountRanking, err := service.GetFishRankingByAmountNoPage(fish.Key)
		if err != nil {
			return nil, err
		}
		sizeRanking, err := service.GetFishRankingBySizeNoPage(fish.Key)
		if err != nil {
			return nil, err
		}
		result = append(result, model.FishingRanking{
			Name: fish.Name,
			Amount: model.RankingInt{
				Rank:  getRank(amountRanking, uuid),
				Value: decodedRecord.Amount[fish.Key],
			},
			MaxSize: model.RankingFloat{
				Rank:  getRank(sizeRanking, uuid),
				Value: decodedRecord.MaxSize[fish.Key],
			},
		})
	}
	return result, nil
}
func getRank(data []model.CustomfishingDataDecoded, uuid string) int {
	for i, record := range data {
		if record.Uuid == uuid {
			return i + 1
		}
	}
	return 0
}
func (service *Service) GetAmountRanking(uuid string) (*model.RankingInt, error) {
	data, err := service.GetFishRankingByTotalAmountNoPage()
	if err != nil {
		return nil, err
	}
	return &model.RankingInt{
		Rank:  getRank(data, uuid),
		Value: service.GetTotalAmount(data[getRank(data, uuid)]),
	}, nil
}
