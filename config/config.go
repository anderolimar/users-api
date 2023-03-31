// This package has all configurations necessary to run Users API
package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port            int
	DBURI           string
	Database        string
	RateLimit       int
	RateLimitTokens int
	ApiUser         string
	ApiPass         string
	ApiHost         string
}

func NewConfig() Config {
	var port int = getIntValue("PORT", 3000)
	return Config{
		Port:            port,
		DBURI:           os.Getenv("MONGODB_URI"),
		Database:        os.Getenv("MONGODB_DATABASE"),
		RateLimit:       getIntValue(os.Getenv("RATE_LIMIT"), 1),
		RateLimitTokens: getIntValue(os.Getenv("RATE_LIMIT_TOKENS"), 5),
		ApiUser:         getStringValue("API_USER", "apiuser"),
		ApiPass:         getStringValue("API_PASS", "apipass"),
		ApiHost:         getStringValue("API_HOST", fmt.Sprintf("localhost:%d", port)),
	}
}

func (c *Config) Validate() {
	if c.DBURI == "" {
		fmt.Println("Invalid MONGODB_URI environment variable")
		os.Exit(0)
	}

	if c.Database == "" {
		fmt.Println("Invalid MONGODB_DATABASE environment variable")
		os.Exit(0)
	}
}

func getIntValue(envName string, defaultValue int) int {
	valueStr := os.Getenv(envName)
	if valueInt, err := strconv.Atoi(valueStr); err == nil {
		return valueInt
	}
	return defaultValue
}

func getStringValue(envName string, defaultValue string) string {
	valueStr := os.Getenv(envName)
	if valueStr != "" {
		return valueStr
	}
	return defaultValue
}
