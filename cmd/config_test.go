package cmd

import (
	"testing"
)

func TestConfig_Get(t *testing.T) {
	config, err := NewConfig("../conf/fake_config.conf")
	if err == nil {
		t.Error("should throw error when read wrong config path.")
	}

	config, err = NewConfig("../conf/whale.conf.template")
	if err != nil {
		t.Errorf("error when read config: %v.", config)
	}

	masterPort := config.Int("master_port")
	if masterPort != 8001 {
		t.Error("read int value error.")
	}

	masterAddr := config.String("master_ip")
	if masterAddr != "127.0.0.1" {
		t.Error("read string value error.")
	}

	emptyInt := config.Int("nil_int")
	if emptyInt != -1 {
		t.Error("should return -1 for non-exist key int.")
	}

	emptyString := config.String("nil_string")
	if emptyString != "" {
		t.Error("should '' nil for not exists string.")
	}
}