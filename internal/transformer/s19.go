package transformer

import (
	"errors"
	"github.com/FedoraTipper/AntHive/internal/constants"
	"github.com/FedoraTipper/AntHive/internal/models"
)

type S19Transformer struct{}

func (*S19Transformer) ConvertStatsPayloadToMiner(friendlyMinerName, crawlerId string, i interface{}) (*models.Miner, error) {

	statsPayload, ok := i.(*constants.S19MinerStats)

	if !ok {
		return nil, errors.New("Unable to cast interface to S19MinerStats struct")
	}

	var fans []models.Fan
	var hashboards []models.Hashboard

	stats := statsPayload.STATS[0]
	info := statsPayload.MinerInfo
	status := statsPayload.Status

	for i, f := range stats.Fan {
		fans = append(fans, models.Fan{
			FanNumber: i,
			RPM:       f,
		})
	}

	for _, hb := range stats.Chain {
		var malfunctioningChipsList []int

		for i, chipStatus := range hb.Asic {
			if string(chipStatus) != " " && string(chipStatus) != constants.Antminer_OKChipStatus {
				malfunctioningChipsList = append(malfunctioningChipsList, i+1)
			}
		}

		hashboards = append(hashboards, models.Hashboard{
			BoardNumber:             hb.Index,
			NoOfChips:               hb.AsicNum,
			MalfunctioningChipsList: malfunctioningChipsList,
			HWErrors:                hb.Hw,
			ChipFrequency:           hb.FreqAvg,
			CurrentHashRate:         hb.RateReal,
			RatedHashRate:           hb.RateIdeal,
			PICTemperature:          hb.TempPic,
			PCBTemperature:          hb.TempPcb,
			ChipTemperature:         hb.TempChip,
			SerialNumber:            hb.Sn,
		})
	}

	return &models.Miner{
		MinerName:    friendlyMinerName,
		CrawlerName:  crawlerId,
		MinerModel:   info.Model,
		Status:       getStatusId(status.StatusCode),
		MinerVersion: info.Version,
		Fans:         fans,
		HashBoards:   hashboards,
	}, nil
}

func getStatusId(s19Status string) constants.MinerStatus {
	if s19Status == constants.S19_OKStatus {
		return constants.MinerOK
	}
	return constants.MinerStopped
}
