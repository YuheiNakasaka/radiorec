package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config struct
type Config struct {
	List *viper.Viper
}

// Init initializes struct
func (c *Config) Init() error {
	b, err := exec.Command("go", "env", "GOPATH").CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to find path to database: %v", err)
	}

	configDir := ""
	for _, p := range filepath.SplitList(strings.TrimSpace(string(b))) {
		p = filepath.Join(p, filepath.FromSlash("/src/github.com/YuheiNakasaka/radiorec/config/"))
		if _, err = os.Stat(p); err == nil {
			configDir = p
			break
		}
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(configDir)
	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Fatal error config file: %v", err)
	}

	c.List = viper.GetViper()
	return err
}
