package test_test

import (
	"testing"

	"dav_converter/configs"
)

func TestGetConfig(t *testing.T) {
	config, err := configs.Get("settings.txt")
	if err != nil {
		t.Error(err)
	}

	if config.PathList != "/path/to/file" {
		t.Error("[!] error config.PathList")
	}

	if !config.SkipExist {
		t.Error("[!] error config.SkipExist")
	}
}
