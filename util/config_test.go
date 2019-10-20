package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitCfg(t *testing.T) {
	err := InitCfg("./config.test.yml")
	assert.Nil(t, err)

	e := InitCfg("./config.test.yml.1")
	assert.NotNil(t, e)
}

func TestGetConfig(t *testing.T) {
	err := InitCfg("./config.test.yml")
	if err == nil {
		cfg := GetConfig()
		assert.Equal(t, 1, len(cfg.Input), "Unmarshal yaml is right")
		assert.Equal(t, true, cfg.Input[0].TailFiles, "Unmarshal yaml value right")
	}
}
