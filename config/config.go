package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Database struct {
		User     string `json:"user"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
	} `json:"database"`
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`
	SessionMicroservice struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"session_microservice"`
	TrackMicroservice struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"track_microservice"`
}

func (c *Config) GetDbConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name,
	)
}

func (c *Config) GetServerConnString() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func (c *Config) GetSessionMicroserviceConnString() string {
	return fmt.Sprintf("%s:%d", c.SessionMicroservice.Host, c.SessionMicroservice.Port)
}

func (c *Config) GetTrackMicroserviceConnString() string {
	return fmt.Sprintf("%s:%d", c.TrackMicroservice.Host, c.TrackMicroservice.Port)
}

func LoadConfig(configFileName string) (*Config, error) {
	file, err := os.Open(configFileName)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	config := new(Config)

	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
