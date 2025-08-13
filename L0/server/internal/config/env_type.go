package config

import "fmt"

type EnvType string

const (
	EnvLocal EnvType = "local"
	EnvProd  EnvType = "prod"
)

// UnmarshalText implements text unmarshalling for EnvType
func (l *EnvType) UnmarshalText(text []byte) error {
	switch string(text) {
	case "local":
		*l = EnvLocal
	case "prod":
		*l = EnvProd
	default:
		return fmt.Errorf("invalid env type: %s", text)
	}
	return nil
}
