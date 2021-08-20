package crawler

import "encoding/json"
import "github.com/FedoraTipper/AntHive/internal/constants"

func ParseStats(data []byte) (*constants.MinerStats, error) {
	var statsPayload constants.MinerStats

	err := json.Unmarshal(data, &statsPayload)

	if err != nil {
		return nil, err
	}

	return &statsPayload, nil
}
