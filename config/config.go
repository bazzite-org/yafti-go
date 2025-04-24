package config

// Action represents an toggable script to be executed on the final screen
type Action struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
	Script      string `json:"script"`
}

type Screen struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Actions     []Action `json:"actions"`
}

// Unmarshaled config file
type Config struct {
	Screens []Screen `json:"screens"`
}
