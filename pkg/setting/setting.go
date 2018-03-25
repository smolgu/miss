package setting

import (
	"os"

	"gopkg.in/yaml.v2"
)

const (
	configFilePath = "./conf/app.yaml"
)

var App struct {
	Vk struct {
		ServiceToken string
	} `yaml:"vk"`
}

var DataPath = "./data"

func NewContext(configFilePaths ...string) (err error) {
	cnfFilePath := configFilePath
	if len(configFilePaths) > 0 {
		cnfFilePath = configFilePaths[0]
	}
	f, err := os.Open(cnfFilePath)
	dec := yaml.NewDecoder(f)
	err = dec.Decode(&App)
	return
}
