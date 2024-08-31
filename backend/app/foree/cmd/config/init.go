package foree_config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func Load(envFilePath string) error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}

	configPath := filepath.Join(ex, envFilePath)
	fmt.Println("*******", configPath)
	err = godotenv.Load(configPath)
	if err != nil {
		return err
	}
	return nil
}
