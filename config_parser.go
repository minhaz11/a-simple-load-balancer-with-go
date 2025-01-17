package main

import (
	"bufio"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port    int
	Servers []string
}

func LoadConfig() (*Config, error) {

	file, err := os.Open("config.yml")

	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileStat, err := file.Stat()

	if err != nil {
		return nil, err
	}

	data := make([]byte, fileStat.Size())
	reader := bufio.NewReader(file)

	_, err = reader.Read(data)

	if err != nil {
		return nil, err
	}

	config := &Config{}

	err = yaml.Unmarshal(data, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
