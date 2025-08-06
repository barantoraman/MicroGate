package config

import (
	"os"

	loaderContract "github.com/barantoraman/microgate/pkg/config/contract"
	loader "github.com/barantoraman/microgate/pkg/config/loader"
)

func GetLoader() loaderContract.Loader {
	env := getEnv("APP_ENV", defaultEnv)
	return loader.NewLoader(env)
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
