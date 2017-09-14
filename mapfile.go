package main

import (
	"io/ioutil"
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
			layers     [][][]string
			layerNames []string
			Inline     []*atlasLayer `mapstructure:"inline"`
			Files      []*atlasLayer `mapstructure:"files"`
		} `mapstructure:"map"`

		Legend map[string]string `mapstructure:"legend"`
	} `mapstructure:"atlas"`

	Ecology struct {
		Classes []Class `mapstructure:"classes"`
	} `mapstructure:"ecology"`

	Organisms map[string]*Organism `mapstructure:"organisms"`
}

type atlasLayer struct {
	Name string `mapstructure:"name"`
	Grid string `mapstructure:"grid"`
}

// ParseMapfile reads and parses a Mapfile given a path.
func ParseMapfile(path string) (mapfile *Mapfile, err error) {
	v := viper.New()
	v.SetTypeByDefaultValue(true)

	v.SetDefault("defaults.empty_tile", ".")
	v.SetDefault("defaults.display_legend", false)

	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	if err = v.ReadInConfig(); err != nil {
		err = errors.Wrapf(err, "error reading Mapfile '%s'", path)
		return
	}

	if err = v.Unmarshal(&mapfile); err != nil {
		err = errors.Wrap(err, "error unmarshaling config")
		return
	}

	if err = mapfile.Sanitize(); err != nil {
		return
	}
	return
}

func (m *Mapfile) Sanitize() (err error) {
	// Assert exactly one map source is provided
	mapSourceInline := len(m.Atlas.Map.Inline) > 0
	mapSourceFiles := len(m.Atlas.Map.Files) > 0
	if !(mapSourceInline || mapSourceFiles) {
		return errors.New("one of ``atlas.map.inline`` or ``atlas.map.files`` must be present")
	}
	if mapSourceInline && mapSourceFiles {
		return errors.New("``atlas.map.inline`` and ``atlas.map.files`` cannot both be present")
	}

	// Read map into layered grid
	var depth int
	var layers []string
	var layerNames []string

	if mapSourceInline {
		depth = len(m.Atlas.Map.Inline)
		layers = make([]string, depth)
		layerNames = make([]string, depth)

		for z := range m.Atlas.Map.Inline {
			data := m.Atlas.Map.Inline[z]
			layers[z] = data.Grid
			layerNames[z] = data.Name
		}

	} else {
		depth = len(m.Atlas.Map.Files)
		layers = make([]string, depth)
		layerNames = make([]string, depth)

		for z := range m.Atlas.Map.Files {
			data := m.Atlas.Map.Files[z]
			bytes, err := ioutil.ReadFile(data.Grid)
			if err != nil {
				return errors.Wrap(err, "error reading ``atlas.map.files``")
			}
			layers[z] = string(bytes)
			layerNames[z] = data.Name
		}
	}
	m.Atlas.Map.layers = gridify(layers)
	m.Atlas.Map.layerNames = layerNames

	// Validate map/legend relationship
	if len(m.Atlas.Legend) == 0 {
		return errors.New("``atlas.legend`` must have at least one entry")
	}
	if err = m.validateMapLegend(layers); err != nil {
		return
	}

	// Validate legend/organism relationship
	if len(m.Organisms) == 0 {
		return errors.New("``organisms`` must have at least one entry")
	}
	if err = m.validateLegendOrganisms(); err != nil {
		return
	}

	if err = m.validateClasses(); err != nil {
		return
	}
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

func (m *Mapfile) validateMapLegend(mapLayers []string) error {
	for _, layer := range m.Atlas.Map.layers {
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

func (m *Mapfile) validateLegendOrganisms() error {
	for _, key := range m.Atlas.Legend {
		_, ok := m.Organisms[key]
		if !ok {
			return errors.Errorf("'%s' is referenced in ``atlas.legend``, but no entry is found in ``organisms``")
		}
	}
	return nil
}

func (m *Mapfile) validateClasses() error {
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

func (m *Mapfile) validateOrganismAttrs() error {
	var result error
	for _, orgData := range m.Organisms {
		if err := vStringMinLen(orgData.Attrs.Name, 2, "name"); err != nil {
			result = multierror.Append(result, err)
		}
		if err := vIntMinVal(orgData.Attrs.Energy, 1, "energy"); err != nil {
			result = multierror.Append(result, err)
		}
		if err := vIntMinVal(orgData.Attrs.Size, 1, "size"); err != nil {
			result = multierror.Append(result, err)
		}
		if err := vIntMinVal(orgData.Attrs.Mass, 1, "mass"); err != nil {
			result = multierror.Append(result, err)
		}
		if orgData.Classes != nil && len(orgData.Classes) > 0 {
			for _, orgClass := range orgData.Classes {
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

func (m *Mapfile) ToWorld() *World {
	atlasLayers := m.Atlas.Map.layers
	layerNames := m.Atlas.Map.layerNames

	height := len(atlasLayers[0])
	width := len(atlasLayers[0][0])
	world := NewWorld(width, height, layerNames)

	for z := range atlasLayers {
		layer := world.Layer(z)

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				symbol := atlasLayers[z][y][x]
				if symbol == m.Defaults.EmptyTile {
					continue
				}

				key := m.Atlas.Legend[symbol]
				data := m.Organisms[key]
				org := NewOrganism(data.Attrs).
					AddClasses(data.Classes...).
					AddAbilities(data.Abilities...).
					Init()

				layer.Add(org, Vec2D(x, y))
			}
		}
	}
	return world
}
