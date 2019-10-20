package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUUID(t *testing.T) {
	assert.NotEqual(t, UUID(), UUID())
}

func TestIsInclude(t *testing.T) {
	pattern := []string{"\\d{4}"}
	text := []byte("1234")
	anotherText := []byte("abcd")
	assert.Equal(t, true, IsInclude(text, pattern))
	assert.NotEqual(t, true, IsInclude(anotherText, pattern))
}

func TestGetTempDir(t *testing.T) {
	p := GetTempDir()
	assert.NotEmpty(t, p)
}

func TestIsDir(t *testing.T) {
	p := "/Users/schoeu/Downloads/git/nma"
	assert.Equal(t, true, IsDir(p))
}

func TestGetAbsPath(t *testing.T) {
	p := "./config.test.yml"
	absPath := GetAbsPath("", p)
	assert.Equal(t, "/", absPath[:1])
}

func TestGetCwd(t *testing.T) {
	p := GetCwd()
	assert.NotEmpty(t, p)
}

func TestPathExist(t *testing.T) {
	p := "../config.test.yml"
	ok, err := PathExist(p)
	assert.Empty(t, err)
	assert.Equal(t, true, ok)
}
