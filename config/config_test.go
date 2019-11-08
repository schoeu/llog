package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitCfg(t *testing.T) {
	err := InitCfg("./test/config.test.yml")
	assert.Nil(t, err)
}

func TestGetConfig(t *testing.T) {
	err := InitCfg("./test/config.test.yml")
	if err == nil {
		cfg := GetConfig()
		assert.Equal(t, 2, len(cfg.Input))
		assert.Equal(t, true, cfg.Input[0].TailFiles)
	}
}
