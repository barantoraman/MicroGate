package config

const (
	defaultEnv = "dev"
)

type Config struct {
	ApiGateway ApiGatewayServiceConfigurations `mapstructure:"api_gateway_service" json:"api_gateway_service"`
	Auth       AuthServiceConfigurations       `mapstructure:"auth_service" json:"auth_service"`
	Task       TaskServiceConfigurations       `mapstructure:"task_service" json:"task_service"`
}

type ApiGatewayServiceConfigurations struct {
	TaskServicePort string `mapstructure:"task_service_port" json:"task_service_port"`
	TaskServiceHost string `mapstructure:"task_service_host" json:"task_service_host"`
	AuthServicePort string `mapstructure:"auth_service_port" json:"auth_service_port"`
	AuthServiceHost string `mapstructure:"auth_service_host" json:"auth_service_host"`
	HTTPPort        string `mapstructure:"http_port" json:"http_port"`
	HTTPHost        string `mapstructure:"http_host" json:"http_host"`
}

type AuthServiceConfigurations struct {
	DBType       string `mapstructure:"db_type" json:"db_type"`
	DBUser       string `mapstructure:"db_user" json:"db_user"`
	DBPass       string `mapstructure:"db_pass" json:"db_pass"`
	DBHost       string `mapstructure:"db_host" json:"db_host"`
	DBPort       int    `mapstructure:"db_port" json:"db_port"`
	DBName       string `mapstructure:"db_name" json:"db_name"`
	DSN          string `mapstructure:"db_dsn" json:"db_dsn"`
	MaxOpenConns int    `mapstructure:"db_max_open_conns" json:"db_max_open_conns"`
	MaxIdleConns int    `mapstructure:"db_max_idle_conns" json:"db_max_idle_conns"`
	MaxIdleTime  string `mapstructure:"db_max_idle_time" json:"db_max_idle_time"`
	SSLMode      string `mapstructure:"db_ssl_mode" json:"db_ssl_mode"`
	StoreType    string `mapstructure:"store_type" json:"store_type"`
	RedisUrl     string `mapstructure:"redis_url" json:"redis_url"`
	RedisPass    string `mapstructure:"redis_pass" json:"redis_pass"`
	GRPCPort     string `mapstructure:"grpc_port" json:"grpc_port"`
	GRPCHost     string `mapstructure:"grpc_host" json:"grpc_host"`
}

type TaskServiceConfigurations struct {
	DBType       string `mapstructure:"db_type" json:"db_type"`
	DBUser       string `mapstructure:"db_user" json:"db_user"`
	DBPass       string `mapstructure:"db_pass" json:"db_pass"`
	DBPort       int    `mapstructure:"db_port" json:"db_port"`
	DBName       string `mapstructure:"db_name" json:"db_name"`
	URI          string `mapstructure:"db_uri" json:"db_uri"`
	MaxOpenConns int    `mapstructure:"db_max_open_conns" json:"db_max_open_conns"`
	MaxIdleConns int    `mapstructure:"db_max_idle_conns" json:"db_max_idle_conns"`
	MaxIdleTime  string `mapstructure:"db_max_idle_time" json:"db_max_idle_time"`
	SSLMode      string `mapstructure:"db_ssl_mode" json:"db_ssl_mode"`
	GRPCPort     string `mapstructure:"grpc_port" json:"grpc_port"`
	GRPCHost     string `mapstructure:"grpc_host" json:"grpc_host"`
}
