package contract

type Loader interface {
	GetConfigByKey(key string, config any) error
}
