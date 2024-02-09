package util

import (
	"log"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application
// The values are read by viper from a config file or environment variable
type Config struct {
	HashSecretKey               string `mapstructure:"HASH_SECRET"`
	FrontEndDomain              string `mapstructure:"FRONTENT_DOMAIN"`
	ServerHost                  string `mapstructure:"SERVER_HOST"`
	ServerPort                  int    `mapstructure:"SERVER_PORT"`
	ServerUsername              string `mapstructure:"SERVER_USERNAME"`
	ServerPassword              string `mapstructure:"SERVER_PASSWORD"`
	FromEmail                   string `mapstructure:"FROM_EMAIL"`
	VerifyEmailSubject          string `mapstructure:"VERIFY_EMAIL_SUBJECT"`
	VerifyEmailTemplate         string `mapstructure:"VERIFY_EMAIL_TEMPLATE"`
	ForgotPasswordEmailSubject  string `mapstructure:"FORGOT_PASSWORD_EMAIL_SUBJECT"`
	ForgotPasswordEmailTemplate string `mapstructure:"FORGOT_PASSWORD_EMAIL_TEMPLATE"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
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
