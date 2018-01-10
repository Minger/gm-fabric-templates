package {{.ConfigPackageName}}

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Initialize loads configuration information
func Initialize() error {
	var configFilePath string
	var envPrefix string
	var err error

	pflag.StringVar(
		&configFilePath,
		"config",
		"",
		"path to config file",
	)

	pflag.StringVar(
		&envPrefix,
		"prefix",
		"",
		"prefix to be stripped from environment variables",
	)
	pflag.Parse()

	if configFilePath != "" {
		viper.SetConfigName(strings.Split(filepath.Base(configFilePath), ".")[0])
		viper.AddConfigPath(filepath.Dir(configFilePath))
	} else {
		viper.SetConfigName("settings")
		viper.AddConfigPath("/etc/{{.ServiceName}}/")
		viper.AddConfigPath("$HOME/.{{.ServiceName}}")
		viper.AddConfigPath(".")
	}

	if err = viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "viper.ReadInConfig")
	}

	if envPrefix != "" {
		viper.SetEnvPrefix(envPrefix)
	}
	viper.AutomaticEnv()

	return nil
}
