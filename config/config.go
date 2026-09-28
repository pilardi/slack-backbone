package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Team represents a single Slack workspace configuration.
type Team struct {
	Name     string `mapstructure:"name"`
	BotToken string `mapstructure:"bot_token"`
	AppToken string `mapstructure:"app_token"`

	AllowedCmds    []string `mapstructure:"commands"`
	RestrictedCmds []string `mapstructure:"restricted_commands"`
}

// Config holds the full application configuration.
type Config struct {
	Teams []Team `mapstructure:"teams"`
}

// Load reads config from flags, .env files, env vars, and a YAML file.
func Load(cmd *cobra.Command) (*Config, error) {
	v := viper.New()

	// Set defaults (lowest priority)
	v.SetDefault("log_level", "info")
	v.SetDefault("config", "teams.yaml")

	// Read .env files if they exist (.env.local overrides .env)
	_, _ = godotenv.Read()
	_, _ = godotenv.Read(".env.local")

	// Read env vars from the environment (higher priority than .env files)
	v.AutomaticEnv()
	v.SetEnvPrefix("slack_backbone")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Bind known keys from env
	_ = v.BindEnv("config", "CONFIG_FILE")

	// Read the config file if a path was provided (highest priority before defaults)
	cfgPath, _ := cmd.Flags().GetString("config")
	if cfgPath != "" {
		v.SetConfigFile(cfgPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file %q: %w", cfgPath, err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// BindFlags registers config keys with a Cobra command.
func BindFlags(cmd *cobra.Command) {
	viper.BindPFlags(cmd.Flags())
}

// Allowed returns true if the command is allowed for this team.
func (t *Team) Allowed(cmd string) bool {
	if len(t.AllowedCmds) == 0 && len(t.RestrictedCmds) == 0 {
		return true // no restrictions = all allowed
	}
	for _, restricted := range t.RestrictedCmds {
		if strings.EqualFold(restricted, cmd) {
			return false
		}
	}
	if len(t.AllowedCmds) > 0 {
		for _, allowed := range t.AllowedCmds {
			if strings.EqualFold(allowed, cmd) {
				return true
			}
		}
		return false
	}
	return true
}
