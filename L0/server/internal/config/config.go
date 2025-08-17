package config

import (
	"fmt"
	"net"
)

type Config struct {
	Env      EnvType `env:"ENV_TYPE"`
	HTTP     HTTPConfig
	Database DatabaseConfig
	Kafka    KafkaConfig
}

type HTTPConfig struct {
	Port         int      `env:"HTTP_PORT"`
	AllowOrigins []string `env:"HTTP_ALLOW_ORIGIN"`
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Dbname   string `env:"DB_NAME"`
}

type KafkaConfig struct {
	Brokers []string `env:"KAFKA_BROKERS"`
	Topic   string   `env:"KAFKA_TOPIC"`
}

func (c *DatabaseConfig) DSN() string {
	socket := net.JoinHostPort(c.Host, c.Port)
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		c.User,
		c.Password,
		socket,
		c.Dbname,
	)
}
