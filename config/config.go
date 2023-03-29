package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port  int
	DBURI string
}

func NewConfig() Config {
	return Config{
		Port:  getIntValue("PORT", 3000),
		DBURI: os.Getenv("MONGODB_URI"),
	}
}

func getIntValue(envName string, defaultValue int) int {
	valueStr := os.Getenv(envName)
	if valueInt, err := strconv.Atoi(valueStr); err == nil {
		return valueInt
	}
	return defaultValue
}
