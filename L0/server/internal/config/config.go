package config

type Config struct {
	Env      EnvType `env-default:"local"`
	
	Database DatabaseConfig
}
