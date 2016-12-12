package cmd

import "testing"

var (
)

func TestConfig_Get(t *testing.T) {
	config, err := NewConfig("../conf/fake_config.conf")
	if err == nil {
		t.Error("should throw error when read wrong config path.")
	}

	config, err = NewConfig("../conf/whale.conf")
	if err != nil {
		t.Errorf("error when read config: %v.", config)
	}

	masterPort := config.Get("master_port")
	if masterPort != "8001" {
		t.Error("read config value error.")
	}

	emptyValue := config.Get("nil_value")
	if emptyValue != "" {
		t.Error("should return nil for not exists key.")
	}
}