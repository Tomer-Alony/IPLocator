package config

import "github.com/spf13/viper"

type Config struct {
	Port string
	RateLimit int
	MaxTokenDuration int
	StorePath string
	StoreKey string
}

func LoadConfig() Config {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("MAX_RATE_LIMIT", "5")
	viper.SetDefault("MAX_TOKEN_DURATION", "500")
	viper.SetDefault("STORE_PATH", "http://api.ipstack.com/%s?access_key=%s")
	viper.SetDefault("STORE_KEY", "90256bbc5e718058dcc02e6a8e38dfa6")

	viper.ReadInConfig()
	viper.AutomaticEnv()

	return Config{
		Port: viper.GetString("PORT"),
		RateLimit: viper.GetInt("MAX_RATE_LIMIT"),
		MaxTokenDuration: viper.GetInt("MAX_TOKEN_DURATION"),
		StorePath: viper.GetString("STORE_PATH"),
		StoreKey: viper.GetString("STORE_KEY"),
	}
}