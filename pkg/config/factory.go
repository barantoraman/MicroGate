package config

import (
	loaderContract "github.com/barantoraman/microgate/pkg/config/contract"
	loader "github.com/barantoraman/microgate/pkg/config/loader"
)

func GetLoader(env string) loaderContract.Loader {
	switch env {
	case "":
		return loader.NewLoader(defaultEnv)
	default:
		return loader.NewLoader(env)
	}
}
