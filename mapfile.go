package ecoscript

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
			layers     [][][]string
			layerNames []string
			Inline     []*layerEntry `mapstructure:"inline"`
			Files      []*layerEntry `mapstructure:"files"`
		} `mapstructure:"map"`

		RawLegend []*legendEntry `mapstructure:"legend"`
		Legend    map[string]string
	} `mapstructure:"atlas"`

	Entities map[string]*Entity `mapstructure:"entities"`
}

type layerEntry struct {
	Name string `mapstructure:"name"`
	Grid string `mapstructure:"grid"`
}

type legendEntry struct {
	Symbol    string `mapstructure:"symbol"`
	EntityKey string `mapstructure:"entity"`
}

// ParseMapfile reads and parses a Mapfile at the given file path.
//
// It uses Viper to read and unmarshal the Mapfile into a Mapfile struct.
// The Mapfile specification is in the docs (TODO).
//
// After reading and unmarshaling, the Mapfile struct is then validated and
// modified with the Mapfile#clean() function.
func ParseMapfile(filePath string) (mapfile *Mapfile, err error) {
	v := viper.New()
	v.SetTypeByDefaultValue(true)

	v.SetDefault("defaults.empty_tile", ".")
	v.SetDefault("defaults.display_legend", false)

	v.SetConfigFile(filePath)
	v.SetConfigType("yaml")
	if err = v.ReadInConfig(); err != nil {
		err = errors.Wrapf(err, "error reading Mapfile '%s'", filePath)
		return
	}

	if err = v.Unmarshal(&mapfile); err != nil {
		err = errors.Wrap(err, "error unmarshaling config")
		return
	}

	if err = mapfile.clean(); err != nil {
		return
	}
	return
}

// Clean validates a Mapfile, preparing it to be converted to a World just
// enough to facilitate validation.
//
// Validate
// --------
// - Assert all required params are set and defined correctly.
// - Assert exactly one map source is provided (inline or file).
// - Assert all symbols used in map are defined in legend
// - Assert no symbol occurs more than once in legend.
// - Assert all entities used in legend are defined in entities.
// - Validate entity attributes.
//
// Prepare
// -------
// - Import world map.
// - Convert raw map into grid of characters ([][]string).
// - Convert raw legend into key/value map.
func (m *Mapfile) clean() (err error) {
	// Validate map input sources
	mapSourceInline := len(m.Atlas.Map.Inline) > 0
	mapSourceFiles := len(m.Atlas.Map.Files) > 0
	if !(mapSourceInline || mapSourceFiles) {
		return errors.New("one of ``atlas.map.inline`` or ``atlas.map.files`` must be present")
	}
	if mapSourceInline && mapSourceFiles {
		return errors.New("``atlas.map.inline`` and ``atlas.map.files`` cannot both be present")
	}

	// Read world map
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

	// Validate and read legend
	if len(m.Atlas.RawLegend) == 0 {
		return errors.New("``atlas.legend`` must have at least one entry")
	}

	m.Atlas.Legend = make(map[string]string)
	for i := range m.Atlas.RawLegend {
		entry := m.Atlas.RawLegend[i]
		_, exists := m.Atlas.Legend[entry.Symbol]
		if exists {
			err = errors.Errorf(
				"symbol '%s' occurs more than once in `atlas.legend``",
				entry.Symbol,
			)
			return errors.WithMessage(err, "symbol must be unique")
		}
		m.Atlas.Legend[entry.Symbol] = entry.EntityKey
	}

	// Validate map/legend relationship
	if err = m.cleanMapLegend(); err != nil {
		return
	}

	// Validate legend/entity relationship
	if len(m.Entities) == 0 {
		return errors.New("``entities`` must have at least one entry")
	}
	if err = m.cleanLegendEntities(); err != nil {
		return
	}

	return
}

func (m *Mapfile) cleanMapLegend() error {
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

func (m *Mapfile) cleanLegendEntities() error {
	for _, key := range m.Atlas.Legend {
		_, ok := m.Entities[key]
		if !ok {
			return errors.Errorf("'%s' is referenced in ``atlas.legend``, but no entry is found in ``entities``", key)
		}
	}
	return nil
}

func (m *Mapfile) cleanEntityAttrs() error {
	var result error
	for _, ent := range m.Entities {
		if err := vStringMinLen(ent.Name, 2, "name"); err != nil {
			result = multierror.Append(result, err)
		}
		if err := vIntMinVal(ent.Attrs.Energy, 1, "energy"); err != nil {
			result = multierror.Append(result, err)
		}
		if err := vIntMinVal(ent.Attrs.Size, 1, "size"); err != nil {
			result = multierror.Append(result, err)
		}
		if err := vIntMinVal(ent.Attrs.Mass, 1, "mass"); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}

func vStringMinLen(val string, min int, key string) (err error) {
	if len(val) < min {
		err = errors.Errorf("entity attribute \"%s\" must have %d or more characters", key, min)
	}
	return
}

func vIntMinVal(val int, min int, key string) (err error) {
	if val < min {
		err = errors.Errorf("entity attribute \"%s\" must be %d or greater", key, min)
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

// ToWorld creates a World from a Mapfile.
//
// Steps
// -----
// - Determine World dimensions.
// - Initialize World.
// - For each tile in each Layer:
//   - Get map symbol for that tile.
//   - If symbol is the empty tile symbol:
//     - Leave the tile blank.
//   - Otherwise:
//     - Look up Entity data in the legend.
//     - Create a new Entity with that data.
//     - Add the Entity to the layer.
// - Return the World.
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
				// Skip empty tiles.
				symbol := atlasLayers[z][y][x]
				if symbol == m.Defaults.EmptyTile {
					continue
				}

				// Create new Entity.
				key := m.Atlas.Legend[symbol]
				data := m.Entities[key]

				abilities := make([]*Ability, len(data.Abilities))
				for i := range data.Abilities {
					rawAbility := data.Abilities[i]
					behavior := Behaviors[rawAbility.Name]
					properties := rawAbility.Properties
					if properties == nil {
						properties = make(Properties)
					}
					abilities[i] = behavior.Ability(properties)
				}

				ent := NewEntity(data.Name, data.Symbol, data.Attrs).
					AddClasses(data.Traits...).
					AddAbilities(abilities...)

				// Add Entity to Layer.
				exec, ok := layer.Add(ent, Vec2D(x, y))
				if !ok {
					log.Printf("couldn't add an entity to a layer")
				}
				exec()
			}
		}
	}
	return world
}
