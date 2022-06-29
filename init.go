package main

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	// Setup logging
	log, _ = zap.NewProduction()
	defer log.Sync()

	// Setup flags and defaults
	pflag.Int("port", defaultPort, "The port to listen on")
	pflag.String("address", defaultAddress, "The address to listen on")
	pflag.String("auth-model", "./config/model.conf", "The path to the auth model file")
	pflag.String("auth-policy", "./config/policy.csv", "The path to the auth policy file")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// Set the environment prefix for incoming variables
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	// Set the configuration file name
	viper.SetConfigName(configFileName)
	// Add the default path(s) to look for the configuration file
	for _, path := range strings.Split(configPaths, ",") {
		log.Info("Added config path", zap.String("path", path))
		viper.AddConfigPath(path)
	}

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("No config file found")
		} else {
			log.Error("Error reading config file", zap.Error(err))
		}
	} else {
		log.Info("Loaded config file", zap.String("file", viper.ConfigFileUsed()))
	}

	// Note that this probably is not a good idea for production, but it's fine for a demo.
	// If the server is rebooted, the tokens will be lost.
	// This won't scale beyond a single server, either.
	log.Info("Initializing maps")
	tokens = make(map[string]string)
	passwords = make(map[string]string)
}
