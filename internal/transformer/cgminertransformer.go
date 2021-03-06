package transformer

import (
	"fmt"

	"github.com/FedoraTipper/AntHive/internal/constants"
	pkgConstants "github.com/FedoraTipper/AntHive/pkg/constants"
	"github.com/FedoraTipper/AntHive/pkg/converter"
	gabsWrapper "github.com/FedoraTipper/AntHive/pkg/gabs"
	"github.com/FedoraTipper/AntHive/pkg/hex"
	"github.com/FedoraTipper/AntHive/pkg/models"
	"github.com/Jeffail/gabs/v2"
)

const (
	hashboardCountField = "miner_count"
	elapsedField        = "Elapsed"
	rateUnitField       = "rate_unit"
	fanNumField         = "fan_num"
	fanSpeedField       = "fan%d"
	asicChipField       = "chain_acs%d"
	chipCountField      = "chain_acn%d"
	tempPCBField        = "temp_pcb%d"
	tempChipField       = "temp_chip%d"
	tempPICField        = "temp_pic%d"
	hardwareErrorField  = "chain_hw%d"
	hashboardRateField  = "chain_rate%d"
	frequencyField      = "freq%d"
)

type CGMinerTransformer struct{}

func (*CGMinerTransformer) ConvertStatsPayloadToMinerStats(friendlyMinerName, crawlerId string, data []byte) (*models.MinerStats, error) {
	jsonMap, err := gabs.ParseJSON(data)

	if err != nil {
		return nil, err
	}

	minerModel, _ := jsonMap.Search("STATS", "0", "Type").Data().(string)
	minerVersion, _ := jsonMap.Search("STATS", "0", "Miner").Data().(string)
	statsMap := jsonMap.Search("STATS", "1")

	wrapper := &gabsWrapper.GabsWrapper{
		Container: statsMap,
	}

	fans := parseFans(wrapper)
	hashboards, err := parseHashBoards(wrapper)

	if err != nil {
		return nil, err
	}

	return &models.MinerStats{
		MinerName:    friendlyMinerName,
		CrawlerName:  crawlerId,
		Status:       pkgConstants.MinerStatusOk,
		MinerModel:   minerModel,
		MinerVersion: minerVersion,
		HashRateUnit: wrapper.GetString(rateUnitField),
		Uptime:       wrapper.GetInt(elapsedField),
		Fans:         fans,
		HashBoards:   hashboards,
		Nonce:        hex.GenerateHexString(8),
	}, nil
}

func parseFans(gabsWrapper *gabsWrapper.GabsWrapper) []models.Fan {
	var fans []models.Fan

	for i := 1; i <= gabsWrapper.GetInt(fanNumField); i++ {
		fans = append(fans, models.Fan{
			FanNumber: i,
			RPM:       gabsWrapper.GetInt(fmt.Sprintf(fanSpeedField, i)),
		})
	}

	return fans
}

func parseHashBoards(gabsWrapper *gabsWrapper.GabsWrapper) ([]models.Hashboard, error) {
	var hashboards []models.Hashboard

	for i := 1; i <= gabsWrapper.GetInt(hashboardCountField); i++ {
		var malfunctioningChipsList []int

		for i, chipStatus := range gabsWrapper.GetString(fmt.Sprintf(asicChipField, i)) {
			if string(chipStatus) != " " && string(chipStatus) != constants.Antminer_OKChipStatus {
				malfunctioningChipsList = append(malfunctioningChipsList, i+1)
			}
		}

		picTempValue := gabsWrapper.GetString(fmt.Sprintf(tempPICField, i))
		pcbTempValue := gabsWrapper.GetString(fmt.Sprintf(tempPCBField, i))
		chipTempValue := gabsWrapper.GetString(fmt.Sprintf(tempChipField, i))

		currentHashRate, err := gabsWrapper.GetFloat64FromString(fmt.Sprintf(hashboardRateField, i))

		if err != nil {
			return nil, err
		}

		hashboards = append(hashboards, models.Hashboard{
			BoardNumber:             i,
			NoOfChips:               gabsWrapper.GetInt(fmt.Sprintf(chipCountField, i)),
			MalfunctioningChipsList: malfunctioningChipsList,
			HWErrors:                gabsWrapper.GetInt(fmt.Sprintf(hardwareErrorField, i)),
			ChipFrequency:           gabsWrapper.GetInt(fmt.Sprintf(frequencyField, i)),
			CurrentHashRate:         currentHashRate,
			PICTemperature:          converter.StringToIntSliceWithDashDelimiter(picTempValue),
			PCBTemperature:          converter.StringToIntSliceWithDashDelimiter(pcbTempValue),
			ChipTemperature:         converter.StringToIntSliceWithDashDelimiter(chipTempValue),
		})
	}

	return hashboards, nil
}
