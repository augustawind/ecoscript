package main

import (
	"io/ioutil"
	"log"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Mapfile is the Settings object that a Mapfile will be marshaled into.
type Mapfile struct {
	Defaults struct {
		EmptyTile     string `mapstructure:"empty_tile"`
		DisplayLegend bool   `mapstructure:"display_legend"`
	} `mapstructure:"defaults"`

	Atlas struct {
		Map struct {
			grid   [][][]string
			Inline []string `mapstructure:"inline"`
			Files  []string `mapstructure:"files"`
		} `mapstructure:"map"`

		Legend map[string]string `mapstructure:"legend"`
	} `mapstructure:"atlas"`

	Ecology struct {
		Classes []Class `mapstructure:"classes"`
	} `mapstructure:"ecology"`

	Organisms map[string]OrganismData `mapstructure:"organisms"`
}

type OrganismData struct {
	*Organism
	Abilities []*Ability `mapstructure:"abilities"`
}

// ParseMapfile reads and parses a Mapfile given a path.
func ParseMapfile(path string) (mapfile Mapfile) {
	v := viper.New()
	v.SetTypeByDefaultValue(true)

	v.SetDefault("defaults.empty_tile", ".")
	v.SetDefault("defaults.display_legend", false)

	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatal(errors.Wrapf(err, "error reading Mapfile '%s'", path))
	}

	if err := v.Unmarshal(&mapfile); err != nil {
		log.Fatal(errors.Wrap(err, "error unmarshaling config"))
	}

	if err := mapfile.Sanitize(); err != nil {
		log.Fatal(err)
	}
	return
}

func (m Mapfile) Sanitize() (err error) {
	// Assert exactly one map source is provided
	mapRawGiven := len(m.Atlas.Map.Inline) > 0
	mapLinkGiven := len(m.Atlas.Map.Files) > 0
	if !(mapRawGiven || mapLinkGiven) {
		return errors.New("one of ``atlas.map.inline`` or ``atlas.map.files`` must be present")
	}
	if mapRawGiven && mapLinkGiven {
		return errors.New("``atlas.map.inline`` and ``atlas.map.files`` cannot both be present")
	}

	// Read map into layered grid
	var mapLayers []string
	if mapRawGiven {
		mapLayers = m.Atlas.Map.Inline
	} else {
		mapLayers := make([]string, len(m.Atlas.Map.Files))
		for z, file := range m.Atlas.Map.Files {
			bytes, err := ioutil.ReadFile(file)
			if err != nil {
				return errors.Wrap(err, "error reading ``atlas.map.files``")
			}
			mapLayers[z] = string(bytes)
		}
	}
	m.Atlas.Map.grid = gridify(mapLayers)

	// Validate map/legend relationship
	if len(m.Atlas.Legend) == 0 {
		return errors.New("``atlas.legend`` must have at least one entry")
	}
	if err = m.validateMapLegend(mapLayers); err != nil {
		return
	}

	// Validate legend/organism relationship
	if len(m.Organisms) == 0 {
		return errors.New("``organisms`` must have at least one entry")
	}
	if err = m.validateLegendOrganisms(); err != nil {
		return
	}

	err = m.validateClasses()
	return
}

func gridify(layers []string) [][][]string {
	stack := make([][][]string, len(layers))
	for z, layer := range layers {
		rows := strings.Split(strings.TrimSpace(layer), "\n")
		grid := make([][]string, len(rows))
		for y, row := range rows {
			grid[y] = strings.Split(strings.TrimSpace(row), "")
		}
		stack[z] = grid
	}
	return stack
}

func (m Mapfile) validateMapLegend(mapLayers []string) error {
	for _, layer := range m.Atlas.Map.grid {
		for _, row := range layer {
			for _, char := range row {
				if char == m.Defaults.EmptyTile {
					continue
				}
				_, ok := m.Atlas.Legend[char]
				if !ok {
					return errors.Errorf("map symbol '%s' not found in ``atlas.legend``", char)
				}
			}
		}
	}
	return nil
}

func (m Mapfile) validateLegendOrganisms() error {
	for _, key := range m.Atlas.Legend {
		_, ok := m.Organisms[key]
		if !ok {
			return errors.Errorf("'%s' is referenced in ``atlas.legend``, but no entry is found in ``organisms``")
		}
	}
	return nil
}

func (m Mapfile) validateClasses() error {
	var result error
	classes := m.Ecology.Classes
	if classes != nil && len(classes) > 0 {
		for _, class := range classes {
			if err := vStringMinLen(string(class), 2, "classes"); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}
	return result
}

func (m Mapfile) validateOrganismAttrs() error {
	var result error
	for _, organism := range m.Organisms {
		if err := vStringMinLen(organism.Attrs.Name, 2, "name"); err != nil {
			result = multierror.Append(result, err)
		}
		if err := vIntMinVal(organism.Attrs.Energy, 1, "energy"); err != nil {
			result = multierror.Append(result, err)
		}
		if err := vIntMinVal(organism.Attrs.Size, 1, "size"); err != nil {
			result = multierror.Append(result, err)
		}
		if err := vIntMinVal(organism.Attrs.Mass, 1, "mass"); err != nil {
			result = multierror.Append(result, err)
		}
		if organism.Classes != nil && len(organism.Classes) > 0 {
			for _, orgClass := range organism.Classes {
				ok := false
				for _, ecoClass := range m.Ecology.Classes {
					if orgClass == ecoClass {
						ok = true
						break
					}
				}
				if !ok {
					err := errors.Errorf("class \"%s\" not found in ``ecology.classes``", orgClass)
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
