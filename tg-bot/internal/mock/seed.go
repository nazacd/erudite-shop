package mock

import (
	"encoding/json"
	"fmt"
	"os"

	"tg-bot/internal/backend"
)

type Seed struct {
	PickupAddress string             `json:"pickup_address"`
	Categories    []backend.Category `json:"categories"`
	Products      []backend.Product  `json:"products"`
}

func LoadSeed(path string) (Seed, error) {
	var seed Seed
	data, err := os.ReadFile(path)
	if err != nil {
		return seed, fmt.Errorf("cannot read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &seed); err != nil {
		return seed, fmt.Errorf("bad JSON in %s: %w", path, err)
	}
	return seed, nil
}
