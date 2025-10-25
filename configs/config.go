package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type RedisConfig struct {
	Host string `mapstructure:"REDIS_HOST"`
	Port string `mapstructure:"REDIS_PORT"`
}

type RabbitMQConfig struct {
	Host string `mapstructure:"RABBITMQ_HOST"`
	Port string `mapstructure:"RABBITMQ_PORT"`
	User string `mapstructure:"RABBITMQ_USERNAME"`
	Pass string `mapstructure:"RABBITMQ_PASSWORD"`
}

type DatabaseConfig struct {
	Host string `mapstructure:"DB_HOST"`
	Port string `mapstructure:"DB_PORT"`
	User string `mapstructure:"DB_USERNAME"`
	Pass string `mapstructure:"DB_PASSWORD"`
	Name string `mapstructure:"DB_NAME"`
}

type Configuration struct {
	Port              string         `mapstructure:"PORT"`
	DatabaseConfig    DatabaseConfig `mapstructure:",squash"`
	RedisConfig       RedisConfig    `mapstructure:",squash"`
	RabbitMQConfig    RabbitMQConfig `mapstructure:",squash"`
	ProductServiceUrl string         `mapstructure:"PRODUCT_SERVICE_URL"`
}

func LoadConfig() (*Configuration, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Configuration
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &config, nil
}
