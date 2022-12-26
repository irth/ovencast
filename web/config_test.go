package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tmp := t.TempDir()
	configPath := filepath.Join(tmp, "./testNewConfig.yaml")

	_, err := os.Stat(configPath)
	assert.ErrorIs(t, err, os.ErrNotExist)

	conf, err := NewConfig(configPath)
	assert.NoError(t, err)

	assert.NotEmpty(t, conf.PasswordHash)
	assert.NotEmpty(t, conf.StreamKey)
	assert.GreaterOrEqual(t, len(conf.StreamKey), 32)

	ok, err := conf.CheckPassword("admin")
	assert.NoError(t, err)
	assert.True(t, ok)

	conf2, err := NewConfig(configPath)
	assert.NoError(t, err)
	assert.Equal(t, conf.StreamKey, conf2.StreamKey)
	assert.Equal(t, conf.PasswordHash, conf2.PasswordHash)

	conf2.StreamKey = "test"
	err = conf2.Save()
	assert.NoError(t, err)

	err = conf.Load()
	assert.NoError(t, err)
	assert.Equal(t, "test", conf.StreamKey)
}

func TestConfigSetAndCheckPassword(t *testing.T) {
	tmp := t.TempDir()
	configPath := filepath.Join(tmp, "./testConfigSetAndCheckPassword.yaml")

	conf, err := NewConfig(configPath)
	assert.NoError(t, err)

	conf.SetPassword("old password")
	ok, err := conf.CheckPassword("old password")
	assert.NoError(t, err)
	assert.True(t, ok)

	conf.SetPassword("new password")
	ok, err = conf.CheckPassword("old password")
	assert.NoError(t, err)
	assert.False(t, ok)
	ok, err = conf.CheckPassword("new password")
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestGenerateStreamKey(t *testing.T) {
	tmp := t.TempDir()
	configPath := filepath.Join(tmp, "./testConfigSetAndCheckPassword.yaml")

	conf, err := NewConfig(configPath)
	assert.NoError(t, err)

	oldStreamKey := conf.StreamKey
	err = conf.GenerateStreamKey()
	assert.NoError(t, err)
	assert.NotEqual(t, oldStreamKey, conf.StreamKey)
	assert.NotEmpty(t, conf.StreamKey)
}
