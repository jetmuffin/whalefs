package cmd

import (
	"os"
	"bufio"
	"strings"
	"io"
)

type Config struct {
	filepath string
	data map[string]string
}

func NewConfig(filepath string) (*Config, error) {
	config := &Config {
		filepath: filepath,
		data: make(map[string]string),
	}

	if err := config.read(); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) Get(key string) string {
	if value, exists := c.data[key]; exists {
		return value
	}
	return ""
}

func (c *Config) read() error {
	file, err := os.Open(c.filepath)
	if err != nil {
		return err
	}

	data := make(map[string]string)
	buf := bufio.NewReader(file)
	for {
		l, err := buf.ReadString('\n')
		line := strings.TrimSpace(l)
		if err != nil {
			if err != io.EOF {
				return err
			}
			if len(line) == 0 {
				break
			}
		}
		switch {
		case len(line) == 0:
		case string(line[0]) == "#":
		default:
			i := strings.IndexAny(line, "=")
			value := strings.TrimSpace(line[i+1 : len(line)])
			key := strings.TrimSpace(line[0:i])
			data[key] = value
		}
	}
	c.data = data
	return nil
}