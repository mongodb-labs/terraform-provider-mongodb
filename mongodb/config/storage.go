package config

// Storage MongoDB configuration for storage parameters
type Storage struct {
	DBPath     string      `yaml:"dbPath"`
	Engine     string      `yaml:"engine"`
	Journal    *Journal    `yaml:"journal"`
	WiredTiger *WiredTiger `yaml:"wiredTiger"`
}

// Journal MongoDB configuration for storage.journal parameters
type Journal struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

// WiredTiger Wired Tiger configuration params
type WiredTiger struct {
	EngineConfig *EngineConfig `yaml:"engineConfig,omitempty"`
}

// EngineConfig Wired Tiger engine configuration params
type EngineConfig struct {
	CacheSizeGB float64 `yaml:"cacheSizeGB,omitempty"`
}
