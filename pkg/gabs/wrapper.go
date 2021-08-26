package gabs

import (
	"strconv"

	"github.com/Jeffail/gabs/v2"
)

type GabsWrapper struct {
	Container *gabs.Container
}

func (g *GabsWrapper) GetInt(path string) int {
	return int(g.GetFloat64(path))
}

func (g *GabsWrapper) GetString(path string) string {
	value, _ := g.Container.Path(path).Data().(string)

	return value
}

func (g *GabsWrapper) GetFloat64(path string) float64 {
	value, _ := g.Container.Path(path).Data().(float64)

	return value
}

func (g *GabsWrapper) GetFloat64FromString(path string) (float64, error) {
	value := g.GetString(path)

	return strconv.ParseFloat(value, 64)
}
