package config

import (
	"iter"
	"log"
	"os"
	"slices"
	"sync/atomic"

	"github.com/goccy/go-yaml"
)

func init() {
	if err := LoadConfig(); err != nil {
		log.Panicf("Failed to load config: %v", err)
	}
	Inhibit.Store(false)
}

var ConfStatus *Config

// Disallow closing the server even if heartbeat is not received
var Inhibit atomic.Bool

// Action represents a toggable script to be executed on the final screen
type Action struct {
	ID          string `json:"id,required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
	Script      string `json:"script"`
}

// GetActionByID searches for an Action with the given ID in the slice of Actions.
// It returns the found Action and its index, or an empty Action and -1 if not found.
func GetActionByID(ag []Action, id string) (action Action, idx int) {
	for i, a := range ag {
		if a.ID == id {
			return a, i
		}
	}
	return Action{}, -1
}

// AddAction adds or updates an Action in the slice of Actions.
// If an Action with the same ID already exists, it will be replaced.
// Returns the modified slice of Actions.
func AddAction(ag []Action, action Action) []Action {
	if _, i := GetActionByID(ag, action.ID); i != -1 {
		ag[i] = action
		return ag
	} else {
		return append(ag, action)
	}
}

// AddAllActions adds or updates multiple Actions in the slice of Actions.
// It calls AddAction for each provided Action.
// Returns the modified slice of Actions.
func AddAllActions(ag []Action, actions ...Action) []Action {
	for _, action := range actions {
		ag = AddAction(ag, action)
	}
	return ag
}

// RemoveActionByID removes an Action with the given ID from the slice of Actions.
// Returns the modified slice of Actions and a boolean indicating if an Action was removed.
func RemoveActionByID(ag []Action, id string) ([]Action, bool) {
	if _, i := GetActionByID(ag, id); i != -1 {
		return slices.Delete(ag, i, i+1), true
	}
	return ag, false
}

type Screen struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Actions     []Action `json:"actions,required"`
}

// Unmarshaled config file
type Config struct {
	Screens []Screen `json:"screens,required"`
}

func (c *Config) GetAllActions() iter.Seq[Action] {
	return func(yield func(Action) bool) {
		for _, screen := range c.Screens {
			for _, action := range screen.Actions {
				if !yield(action) {
					return
				}
			}
		}
	}
}

func (c *Config) GetActionsByIds(ids []string) ([]Action, bool) {
	res := []Action{}
	for act := range c.GetAllActions() {
		for _, id := range ids {
			if act.ID == id {
				res = append(res, act)
			}
		}
	}
	return res, len(res) > 0
}

func LoadConfig() error {
	var configPath string
	if envPath := os.Getenv("YAFTI_CONF"); envPath != "" {
		configPath = envPath
	} else {
		configPath = "/usr/share/yafti/yafti.yml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	// Log the contents of the config file for debugging
	log.Printf("Loaded config file contents: %s", string(data))

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	ConfStatus = &config
	return nil
}
