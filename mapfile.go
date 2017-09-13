package main

import (
	"io/ioutil"
	"log"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
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

	Organisms map[string]Organism `mapstructure:"organisms"`
}

// ParseMapfile reads and parses a Mapfile given a path.
func ParseMapfile(path string) (mapfile Mapfile) {
	v := viper.New()
	v.SetTypeByDefaultValue(true)

	v.SetDefault("defaults.empty_tile", " ")
	v.SetDefault("defaults.use_map_symbols", false)

	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal(errors.Wrapf(err, "error reading Mapfile '%s'", path))
	}

	err = v.UnmarshalExact(&mapfile)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error unmarshaling config"))
	}

	err = mapfile.Validate()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (m Mapfile) Validate() (err error) {
	var mapText string
	mapRawGiven := len(m.Atlas.Map.Raw) > 0
	mapLinkGiven := len(m.Atlas.Map.Link) > 0
	if !(mapRawGiven || mapLinkGiven) {
		return errors.New("one of ``atlas.map.raw`` or ``atlas.map.link`` must be present")
	}
	if mapRawGiven && mapLinkGiven {
		return errors.New("``atlas.map.raw`` and ``atlas.map.link`` cannot both be present")
	}
	if mapRawGiven {
		mapText = m.Atlas.Map.Raw
	} else {
		bytes, err := ioutil.ReadFile(m.Atlas.Map.Link)
		if err != nil {
			return errors.Wrap(err, "error reading ``atlas.map.link``")
		}
		mapText = string(bytes)
	}

	if len(m.Atlas.Legend) == 0 {
		return errors.New("``atlas.legend`` must have at least one entry")
	}
	err = m.validateMap(mapText)
	if err != nil {
		return
	}

	if len(m.Organisms) == 0 {
		return errors.New("``organisms`` must have at least one entry")
	}
	err = m.validateLegend()
	if err != nil {
		return
	}

	return
}

func (m Mapfile) validateMap(mapText string) error {
	for _, char := range strings.Split(mapText, "") {
		_, ok := m.Atlas.Legend[char]
		if !ok {
			return errors.Errorf("map symbol '%s' not found in ``atlas.legend``", char)
		}
	}
	return nil
}

func (m Mapfile) validateLegend() error {
	for _, key := range m.Atlas.Legend {
		_, ok := m.Organisms[key]
		if !ok {
			return errors.Errorf("'%s' is referenced in ``atlas.legend``, but no entry is found in ``organisms``")
		}
	}
	return nil
}

func (m Mapfile) validateOrganisms() error {
	var result error
	for _, organism := range m.Organisms {
		if err := vStringMinLen(organism.Stats.Name, 2, "name"); err != nil {
			result = multierror.Append(result, err)
		}
		if err := vIntMinVal(organism.Stats.Energy, 1, "energy"); err != nil {
			result = multierror.Append(result, err)
		}
		if err := vIntMinVal(organism.Stats.Size, 1, "size"); err != nil {
			result = multierror.Append(result, err)
		}
		if err := vIntMinVal(organism.Stats.Mass, 1, "mass"); err != nil {
			result = multierror.Append(result, err)
		}
		if organism.Classes != nil {
			for _, class := range organism.Classes {
				if err := vStringMinLen(string(class), 2, "classes"); err != nil {
					result = multierror.Append(result, err)
				}
			}
		}
	}
	return result
}

func vStringMinLen(val string, min int, key string) (err error) {
	if len(val) < min {
		err = errors.Errorf("organism attribute \"%s\" must have %d or more characters", key, min)
	}
	return
}

func vIntMinVal(val int, min int, key string) (err error) {
	if val < min {
		err = errors.Errorf("organism attribute \"%s\" must be %d or greater", key, min)
	}
	return
}

func (m Mapfile) ToWorld() *World {
	return nil
}
