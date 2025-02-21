package configs

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	Port                 string        `mapstructure:"PORT"`
	DbUsername           string        `mapstructure:"DB_USERNAME"`
	DbPassword           string        `mapstructure:"DB_PASSWORD"`
	DbPort               string        `mapstructure:"DB_PORT"`
	DbName               string        `mapstructure:"DB_NAME"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, err
}
