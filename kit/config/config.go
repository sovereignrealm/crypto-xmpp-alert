package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	XmppUsername     string `required:"true" split_words:"true"`
	XmppPassword     string `required:"true" split_words:"true"`
	XmppServer       string `required:"true" split_words:"true"`
	XmppRecipientJid string `required:"true" split_words:"true"`
}

func findProjectRoot(startDir, markerFile string) (string, error) {
	// Start at the current directory and move up until root is reached or the marker file is found
	for dir := startDir; dir != "/"; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, markerFile)); err == nil {
			return dir, nil
		}
	}
	return "", fmt.Errorf("project root with marker file '%s' not found", markerFile)
}

func NewConfig(envFile string) *Config {
	_, filename, _, _ := runtime.Caller(0)
	startDir := filepath.Dir(filename)
	rootPath, err := findProjectRoot(startDir, "go.mod")
	if err != nil {
		panic(err) // or handle error more gracefully
	}

	if envFile != "" {
		err := godotenv.Load(filepath.Join(rootPath, envFile))
		if err != nil {
			panic(fmt.Errorf("godotenv.Load: %w", err))
		}
	}

	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		panic(fmt.Errorf("envconfig.Process: %w", err))
	}

	return cfg
}
