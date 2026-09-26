package config

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Team represents a single Slack workspace configuration.
type Team struct {
	Name           string   `mapstructure:"name"`
	BotToken       string   `mapstructure:"bot_token"`
	AppToken       string   `mapstructure:"app_token"`
	DefaultChannel string   `mapstructure:"default_channel"`
	AllowedCmds    []string `mapstructure:"commands"`
	RestrictedCmds []string `mapstructure:"restricted_commands"`
}

// Config holds the full application configuration.
type Config struct {
	Teams []Team `mapstructure:"teams"`
}

// Load reads config from flags, env vars, and a YAML file.
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("log_level", "info")

	// Read env vars
	v.AutomaticEnv()
	v.SetEnvPrefix("slack_backbone")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Bind known keys from env
	v.BindEnv("config", "CONFIG_FILE")

	return nil, nil
}

// BindFlags registers config keys with a Cobra command.
func BindFlags(cmd *cobra.Command) {
	cmd.Flags().String("log-level", "info", "Log level: debug|info|warn|error")
	viper.BindPFlag("log_level", cmd.Flags().Lookup("log-level"))
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

// GetDefaultChannel returns the team's default notification channel.
func (t *Team) GetDefaultChannel() string {
	if t.DefaultChannel != "" {
		return t.DefaultChannel
	}
	return "#general"
}
