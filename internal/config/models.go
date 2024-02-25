package config

import (
	"fmt"
	"time"
)

type Config struct {
	Service          Service  `json:"service,omitempty"`
	Database         Database `json:"database,omitempty"`
	URL              string   `json:"url,omitempty"`
	AllCurrenciesURL string   `json:"allcurrenciesurl,omitempty"`
}

type Service struct {
	Name    string `json:"name,omitempty"`
	Host    string `json:"host,omitempty"`
	Port    string `json:"port,omitempty"`
	Timeout int    `json:"timeout,omitempty"`
}

type Database struct {
	Type     string `json:"type,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name,omitempty"`
	SslMode  string `json:"sslMode,omitempty"`
}

func (c *Config) Address() string {
	return c.Service.Host + ":" + c.Service.Port
}

func (c *Config) WriteTimeout() time.Duration {
	return time.Duration(c.Service.Timeout) * time.Second
}

func (c *Config) ReadTimeout() time.Duration {
	return time.Duration(c.Service.Timeout) * time.Second
}

func (c *Config) IdleTimeout() time.Duration {
	return time.Duration(c.Service.Timeout*4) * time.Second
}

func (c *Config) Wait() time.Duration {
	return time.Duration(c.Service.Timeout)
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SslMode,
	)
}
