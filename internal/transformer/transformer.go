package transformer

import (
	"errors"
	"fmt"

	"github.com/FedoraTipper/AntHive/internal/constants"
	"github.com/FedoraTipper/AntHive/pkg/models"
)

type ITransformer interface {
	ConvertStatsPayloadToMinerStats(friendlyMinerName string, crawlerName string, data []byte) (*models.MinerStats, error)
}

func GetTransformer(model constants.MinerSeries) (ITransformer, error) {
	switch model {
	case constants.X19:
		return &CGMinerTransformer{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown model: %s", model))
	}
}
