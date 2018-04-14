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

// Init initializes struct.
// Default path is $GOPATH/src/github.com/YuheiNakasaka/radiorec/config/
func (c *Config) Init() error {
	configDir := os.Getenv("CONFIG_DIR")
	if configDir == "" {
		b, err := exec.Command("go", "env", "GOPATH").CombinedOutput()
		if err != nil {
			return fmt.Errorf("Failed to get GOPATH: %v", err)
		}

		for _, p := range filepath.SplitList(strings.TrimSpace(string(b))) {
			p = filepath.Join(p, filepath.FromSlash("/src/github.com/YuheiNakasaka/radiorec/config/"))
			if _, err = os.Stat(p); err == nil {
				configDir = p
				break
			}
		}
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(configDir)
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Fatal error config file: %v", err)
	}

	c.List = viper.GetViper()
	return err
}
