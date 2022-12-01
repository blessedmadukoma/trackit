package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB_DIALECT      string `mapstructure:"DB_DIALECT"`
	DB_DRIVER      string `mapstructure:"DB_DRIVER"`
	DB_USER        string `mapstructure:"DB_USER"`
	DB_PASSWORD    string `mapstructure:"DB_PASSWORD"`
	DB_NAME        string `mapstructure:"DB_NAME"`
	DB_PORT        string `mapstructure:"DB_PORT"`
	DB_HOST        string `mapstructure:"DB_HOST"`
	DB_TIMEZONE    string `mapstructure:"DB_TIMEZONE"`
	SERVER_ADDRESS string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (cfg Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&cfg)
	return

}
