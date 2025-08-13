package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	instance *Config
	once     sync.Once
)

// MustNew returns the application configuration.
// Note that config.Config is a singleton
func MustNew() Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadEnv(instance); err != nil {
			panic("failed to read application configuration")
		}
	})
	return *instance
}
