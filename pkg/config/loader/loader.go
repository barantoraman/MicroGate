package viper

import (
	"fmt"

	loaderContract "github.com/barantoraman/microgate/pkg/config/contract"
	"github.com/spf13/viper"
)

type loader struct {
	configFile string
	v          *viper.Viper
}

func NewLoader(env string) loaderContract.Loader {
	configFile := fmt.Sprintf("service.%s.yaml", env)
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(configFile)
	v.AddConfigPath(".")
	return &loader{
		configFile: configFile,
		v:          v,
	}
}

func (l *loader) GetConfigByKey(key string, config any) error {
	if err := l.v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file %s", err.Error())
	}
	if err := l.v.UnmarshalKey(key, config); err != nil {
		return fmt.Errorf("failed to unmarshal config %s", err.Error())
	}
	return nil
}
