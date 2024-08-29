package foree_config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPasswd               string
	DBAddr                 string
	DBName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
}

var Envs = initConfig()

func initConfig() Config {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	configPath := filepath.Join(ex, "../../deploy/.local_env")
	err = godotenv.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}
	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port:       getEnv("PORT", "8080"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPasswd:   getEnv("DB_PASSWD", ""),
		DBAddr:     fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:     getEnv("DB_NAME", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

// func getEnvAsInt(key string, fallback int64) int64 {
// 	if value, ok := os.LookupEnv(key); ok {
// 		i, err := strconv.ParseInt(value, 10, 64)
// 		if err != nil {
// 			return fallback
// 		}

// 		return i
// 	}

// 	return fallback
// }