package setting

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	configFilePath = "./conf/app.yaml"
)

// Dev describe is now running under dev
var Dev bool

// App is global setting
var App struct {
	Vk struct {
		ServiceToken string `yaml:"service_token"`
	} `yaml:"vk"`
	DataPath string `yaml:"data_path"`
	Secret   string
}

// NewContext opent conf file, parse it
func NewContext(configFilePaths ...string) (err error) {
	cnfFilePath := configFilePath
	if len(configFilePaths) > 0 {
		cnfFilePath = configFilePaths[0]
	}
	log.Printf("read config: %v", cnfFilePath)
	f, err := os.Open(cnfFilePath)
	if err != nil {
		return errors.Wrap(err, "open config file")
	}
	dec := yaml.NewDecoder(f)
	err = dec.Decode(&App)
	if err != nil {
		return errors.Wrap(err, "decode config file")
	}
	log.Printf("setting: vk_service_token: %s", App.Vk.ServiceToken)
	log.Printf("setting: data_path: %s", App.DataPath)
	return
}
