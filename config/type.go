package config

type Config struct {
	Database DatabaseConfig `mapstructure:"database" yaml:"database"`
	Redis    RedisConfig    `mapstructure:"redis" yaml:"redis"`
	Log      LogConfig      `mapstructure:"log" yaml:"log"`
	//jwt签名密钥
	JwtSignKeyHex string `mapstructure:"jwt_sign_key_hex" yaml:"jwt_sign_key_hex"`
	//jwt加密密钥
	JwtSecretHex string `mapstructure:"jwt_secret_hex" yaml:"jwt_secret_hex"`
	//服务端口号
	ServerPort string `mapstructure:"server_port" yaml:"server_port"`
	//视频配置
	Vedio VedioConfig `mapstructure:"vedio" yaml:"vedio"`
	//推荐系统配置（预留）
	Recommendation RecommendConfig `mapstructure:"recommendation" yaml:"recommendation"`
}

// DatabaseConfig 数据库配置 (支持 PostgreSQL)
type DatabaseConfig struct {
	Driver   string `mapstructure:"driver" yaml:"driver"` // postgres, mysql
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Dbname   string `mapstructure:"dbname" yaml:"dbname"`
	SSLMode  string `mapstructure:"sslmode" yaml:"sslmode"` // PostgreSQL: disable, require, verify-full
	Timezone string `mapstructure:"timezone" yaml:"timezone"` // 时区，例如: Asia/Shanghai
	//"10s"
	Timeout string `mapstructure:"timeout" yaml:"timeout"`
	//  设置连接池中空闲连接的最大数量。
	MaxIdleConns int `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	//  设置打开数据库连接的最大数量。
	MaxOpenConns int `mapstructure:"max_open_conns" yaml:"max_open_conns"`
	// 连接最大生命周期 (e.g., "1h")
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"`
}

// MysqlConfig 兼容旧配置（已弃用，请使用 DatabaseConfig）
// Deprecated: Use DatabaseConfig instead
type MysqlConfig = DatabaseConfig

type LogConfig struct {
	Path         string `mapstructure:"path" yaml:"path"`
	PanicLogName string `mapstructure:"panic_log_name" yaml:"panic_log_name"`
	//trace,debug,info,warn,error,fatal,panic
	Level string `mapstructure:"level" yaml:"level"`
}

type JwtConfig struct {
	SignKeyHex string `mapstructure:"sign_key_hex" yaml:"sign_key_hex"`
	SecretHex  string `mapstructure:"secret_hex" yaml:"secret_hex"`
}

type VedioConfig struct {
	//视频存储的根目录
	BasePath string `mapstructure:"base_path" yaml:"base_path"`
	// e.g. static
	UrlPrefix string `mapstructure:"url_prefix" yaml:"url_prefix"`
	// e.g. http://localhost:8080
	Domain string `mapstructure:"domain" yaml:"domain"`
}

// RecommendConfig 推荐系统配置（预留）
type RecommendConfig struct {
	// 是否启用真实推荐系统
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`
	// Python推荐API地址（预留）
	APIURL string `mapstructure:"api_url" yaml:"api_url"`
	// 超时时间
	Timeout string `mapstructure:"timeout" yaml:"timeout"`
	// Mock配置
	Mock MockConfig `mapstructure:"mock" yaml:"mock"`
}

// MockConfig Mock数据配置
type MockConfig struct {
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Enabled  bool   `mapstructure:"enabled" yaml:"enabled"`   // 是否启用Redis
	Host     string `mapstructure:"host" yaml:"host"`         // Redis地址
	Port     int    `mapstructure:"port" yaml:"port"`         // Redis端口
	Password string `mapstructure:"password" yaml:"password"` // Redis密码（可选）
	DB       int    `mapstructure:"db" yaml:"db"`             // Redis数据库索引（0-15）
	// 连接池配置
	PoolSize     int    `mapstructure:"pool_size" yaml:"pool_size"`         // 最大连接数
	MinIdleConns int    `mapstructure:"min_idle_conns" yaml:"min_idle_conns"` // 最小空闲连接数
	MaxRetries   int    `mapstructure:"max_retries" yaml:"max_retries"`     // 最大重试次数
	// 缓存配置
	DefaultExpiration string `mapstructure:"default_expiration" yaml:"default_expiration"` // 默认过期时间（如 "1h", "30m"）
}
