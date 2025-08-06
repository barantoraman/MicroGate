package config

const (
	defaultEnv = "dev"
)

type Config struct {
	ApiGateway ApiGatewayServiceConfigurations `mapstructure:"api_gateway_service"`
	Auth       AuthServiceConfigurations       `mapstructure:"auth_service"`
	Task       TaskServiceConfigurations       `mapstructure:"task_service"`
}

type ApiGatewayServiceConfigurations struct {
	TaskServicePort string `mapstructure:"task_service_port"`
	TaskServiceHost string `mapstructure:"task_service_host"`
	AuthServicePort string `mapstructure:"auth_service_port"`
	AuthServiceHost string `mapstructure:"auth_service_host"`
	HTTPPort        string `mapstructure:"http_port"`
	HTTPHost        string `mapstructure:"http_host"`
}

type AuthServiceConfigurations struct {
	DBType       string `mapstructure:"db_type"`
	DBUser       string `mapstructure:"db_user"`
	DBPass       string `mapstructure:"db_pass"`
	DBHost       string `mapstructure:"db_host"`
	DBPort       int    `mapstructure:"db_port"`
	DBName       string `mapstructure:"db_name"`
	DSN          string `mapstructure:"db_dsn"`
	MaxOpenConns int    `mapstructure:"db_max_open_conns"`
	MaxIdleConns int    `mapstructure:"db_max_idle_conns"`
	MaxIdleTime  string `mapstructure:"db_max_idle_time"`
	SSLMode      string `mapstructure:"db_ssl_mode"`
	StoreType    string `mapstructure:"store_type"`
	RedisUrl     string `mapstructure:"redis_url"`
	RedisPass    string `mapstructure:"redis_pass"`
	GRPCPort     string `mapstructure:"grpc_port"`
	GRPCHost     string `mapstructure:"grpc_host"`
}

type TaskServiceConfigurations struct {
	DBType       string `mapstructure:"db_type"`
	DBUser       string `mapstructure:"db_user"`
	DBPass       string `mapstructure:"db_pass"`
	DBPort       int    `mapstructure:"db_port"`
	DBName       string `mapstructure:"db_name"`
	URI          string `mapstructure:"db_uri"`
	MaxOpenConns int    `mapstructure:"db_max_open_conns"`
	MaxIdleConns int    `mapstructure:"db_max_idle_conns"`
	MaxIdleTime  string `mapstructure:"db_max_idle_time"`
	SSLMode      string `mapstructure:"db_ssl_mode"`
	GRPCPort     string `mapstructure:"grpc_port"`
	GRPCHost     string `mapstructure:"grpc_host"`
}
