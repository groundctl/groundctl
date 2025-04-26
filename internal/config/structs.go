package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

var (
	RabbitMQ struct {
		Host        string `mapstructure:"host"`
		Port        int    `mapstructure:"port"`
		Credentials map[string]struct {
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
		} `mapstructure:"credentials"`
	}
)

func init() {
	// Set up viper to read the config.yml file
	viper.SetConfigName("config")                   // Config file name without extension
	viper.SetConfigType("yml")                      // Config file type
	viper.AddConfigPath("/etc/groundctl/")          // Look for the config file in the etc dir
	viper.AddConfigPath("$HOME/.config/groundctl/") // Look for the config file in the .config dir
	viper.AddConfigPath(".")                        // Look for the config file in the current directory

	viper.AutomaticEnv()
	viper.SetEnvPrefix("env")                              // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // this is useful e.g. want to use . in Get() calls, but environmental variables to use _ delimiters (e.g. app.port -> APP_PORT)

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	viper.UnmarshalKey("rabbitmq", &RabbitMQ)
}
