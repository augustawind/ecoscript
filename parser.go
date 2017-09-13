package main

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Mapfile is the Settings object that a Mapfile will be marshalled into.
type Mapfile struct {
	Defaults struct {
		EmptyTile     string `mapstructure:"empty_tile"`
		DisplayLegend bool   `mapstructure:"display_legend"`
	} `mapstructure:"defaults"`

	Atlas struct {
		Map struct {
			Raw  string `mapstructure:"raw"`
			Link string `mapstructure:"link"`
		} `mapstructure:"map"`

		Legend map[string]string `mapstructure:"legend"`
	} `mapstructure:"atlas"`

	Ecology struct {
		Classes []string `mapstructure:"classes"`
	} `mapstructure:"ecology"`

	Organisms map[string]string `mapstructure:"organisms"`
}

// ParseMapfile reads and parses a Mapfile given a path.
func ParseMapfile(path string) (mapfile Mapfile) {
	v := viper.New()
	v.SetTypeByDefaultValue(true)

	v.SetDefault("defaults.empty_tile", " ")
	v.SetDefault("defaults.use_map_symbols", false)

	v.SetConfigFile(path)
	v.SetConfigType("toml")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal(errors.Errorf("error reading Mapfile '%s': %v", path, err))
	}

	err = v.UnmarshalExact(&mapfile)
	if err != nil {
		log.Fatal(errors.Errorf("error unmarshaling config: %v", err))
	}
	return
}
