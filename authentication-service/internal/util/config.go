package util

import (
	"log"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application
// The values are read by viper from a config file or environment variable
type ConfigVars struct {
	HashSecretKeyVerifyEmail    string `mapstructure:"HASH_SECRET_VERIFY_EMAIL"`
	HashSecretKeyForgotPassword string `mapstructure:"HASH_SECRET_FORGOT_PASSWORD"`
	FrontEndDomain              string `mapstructure:"FRONTENT_DOMAIN"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config ConfigVars, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		log.Println(err)
		return
	}

	err = viper.Unmarshal(&config)
	return
}
