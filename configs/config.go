package configs

import (
	"sync"

	"github.com/DmitryVesenniy/goconfig/ini"
)

var (
	config *ConfigApp
	once   sync.Once
)

type ConfigApp struct {
	PathList  string `ini:"PATH"`
	PathOut   string `ini:"OUT_PATH"`
	SkipExist bool   `ini:"SKIP"`
}

// Get Exports
func Get(fileIni string) (*ConfigApp, error) {
	var err error
	once.Do(func() {
		funcGetConf := ini.Get(&ConfigApp{}, fileIni)
		var cfg interface{}
		cfg, err = funcGetConf()
		config, _ = cfg.(*ConfigApp)

	})
	return config, err
}
