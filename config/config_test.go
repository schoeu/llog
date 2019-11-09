package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var configPath = "../test/test.yml"

func TestInitCfg(t *testing.T) {
	err := InitCfg(configPath)
	assert.Nil(t, err)
}

func TestGetConfig(t *testing.T) {
	err := InitCfg(configPath)
	if err == nil {
		cfg := GetConfig()
		assert.Equal(t, 2, len(cfg.Input))
		assert.Equal(t, true, cfg.Input[0].TailFiles)
	}
}
