package config

import "github.com/spf13/viper"

const (
	EnvDev  = "dev"
	EnvProd = "prod"
	EnvTest = "test"
)

type Config struct {
	DBSource        string `mapstructure:"DB_SOURCE"`
	Addr            string `mapstructure:"ADDR"`
	Env             string `mapstructure:"ENV"`
	AWSBucketName   string `mapstructure:"AWS_BUCKET_NAME"`
	AWSBucketRegion string `mapstructure:"AWS_BUCKET_REGION"`
	CDNAddress      string `mapstructure:"CDN_ADDRESS"`
	RedisAddr       string `mapstructure:"REDIS_ADDR"`
	RedisPass       string `mapstructure:"REDIS_PASS"`
	AppURL          string `mapstructure:"APP_URL"`
	LocalstackURL   string `mapstructure:"LOCALSTACK_URL"`
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
