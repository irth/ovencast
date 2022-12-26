package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

type Config struct {
	path string

	StreamKey    string `yaml:"stream-key"`
	PasswordHash string `yaml:"password-hash"`

	sync.RWMutex `yaml:"-"`
}

func randStr(length int) (string, error) {
	u := make([]byte, length/2)
	_, err := rand.Read(u)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(u)[:length], nil
}

func NewConfig(path string) (*Config, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	conf := Config{
		path: path,
	}

	err = conf.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &conf, conf.Initialize()
		}
		return nil, err
	}

	return &conf, nil
}

func (c *Config) Initialize() error {
	if err := c.GenerateStreamKey(); err != nil {
		return err
	}

	if err := c.SetPassword("admin"); err != nil {
		return err
	}

	return c.Save()
}

func (c *Config) Load() error {
	f, err := os.Open(c.path)
	if err != nil {
		return fmt.Errorf("open %s: %w", c.path, err)
	}
	defer f.Close()

	return yaml.NewDecoder(f).Decode(c)
}

func (c *Config) Save() error {
	f, err := os.Create(c.path)
	if err != nil {
		return fmt.Errorf("open %s: %w", c.path, err)
	}
	defer f.Close()

	return yaml.NewEncoder(f).Encode(c)
}

func (c *Config) GenerateStreamKey() error {
	streamKey, err := randStr(64)
	if err != nil {
		return fmt.Errorf("generating stream key: %w", err)
	}
	c.StreamKey = streamKey
	return nil
}

func (c *Config) SetPassword(password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing default password: %w", err)
	}
	c.PasswordHash = string(passwordHash)
	return nil
}

func (c *Config) CheckPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(c.PasswordHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
