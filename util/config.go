package util

import "github.com/spf13/viper"

type Config struct {
	DBSource string `mapstructure:"DB_SOURCE"`
	ServerAddress string	`mapstructure:"SERVER_ADDRESS"`
	GrpcServerAddress string	`mapstructure:"GRPC_SERVER_ADDRESS"`
	SymmetricKey string	`mapstructure:"SYMMETRIC_KEY"`
}

func LoadConfig(path string)(config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.AutomaticEnv()

	err =viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	err  = viper.Unmarshal(&config)
	return config, err
}