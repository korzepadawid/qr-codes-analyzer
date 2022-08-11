package config

import "github.com/spf13/viper"

type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	Addr     string `mapstructure:"ADDR"`
}

func Load(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
