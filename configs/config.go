package configs

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func Init() {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = "envs/.env.dev"
	}

	viper.SetConfigFile(envFile)
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf(
			"config: no config file loaded (%v), using ENV or defaults",
			err,
		)
	}

	viper.AutomaticEnv()

	viper.SetDefault("SERVICE_NAME", "transaction-manager")
	viper.SetDefault("HTTP_PORT", "3000")
	viper.SetDefault("DATABASE_URL", "postgres://user:password@postgres:5432/app_db?sslmode=disable")
	viper.SetDefault("DB_MAX_CONNS", 5)
	viper.SetDefault("DB_MAX_CONN_LIFETIME", "300s")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_HOST", "redis")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("KAFKA_ADDRESS", "kafka:9092")
	viper.SetDefault("KAFKA_TOPIC", "transactions")
	viper.SetDefault("KAFKA_GROUP_ID", "transaction-group")
}
