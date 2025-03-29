package conf

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config Variables //
var ServerConfig = InitApiServerConfig()
var DatabaseConfig = InitPostgresSqlConfig()

// Config structs //
type PostgreSqlConfig struct {
	DBUser     string
	DBPassword string
	DBAddress  string
	DBPort     string
	DBName     string
	SSLMode    string
}

type ApiServerConfig struct {
	PublicHost                    string
	Port                          string
	JWTSecret                     string
	JWTExpirationInSeconds        int64
	RefreshTokenSecret            string
	RefreshTokenExpirationInHours int64
}

// Configs Functions //
func InitPostgresSqlConfig() PostgreSqlConfig {
	godotenv.Load()

	return PostgreSqlConfig{
		DBUser:     getEnv("DB_USER", "psql_admin"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBAddress:  getEnv("DB_HOST", "127.0.0.1"),
		DBName:     getEnv("DB_NAME", "smartspend_db"),
		SSLMode:    getEnv("SSL_MODE", "disable"),
	}
}

func InitApiServerConfig() ApiServerConfig {
	godotenv.Load()

	return ApiServerConfig{
		PublicHost:                    getEnv("PUBLIC_HOST", "0.0.0.0"),
		Port:                          getEnv("PORT", "8080"),
		JWTSecret:                     getEnv("JWT_SECRET", "not-so-secret-now-is-it?"),
		JWTExpirationInSeconds:        getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*1),
		RefreshTokenSecret:            getEnv("REFRESH_TOKEN_SECRET", "not-so-secret-now-is-it?"),
		RefreshTokenExpirationInHours: getEnvAsInt("REFRESH_TOKEN_EXPIRATION_IN_HOURS", 30*24),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
