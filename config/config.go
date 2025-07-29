package config

func GetValue(key string, defaultKey string) string {
	viper := NewViperConfig()
	value := viper.GetString(key)
	if value == "" {
		value = viper.GetString(defaultKey)
	}

	return value
}

var (
	RedisHost             = GetValue("redis.host", "")
	RedisPort             = GetValue("redis.port", "")
	RedisPass             = GetValue("redis.password", "")
	JwtSecret             = GetValue("jwt.secret", "")
	JwtTokenAccessExpire  = GetValue("jwt.access_expire", "")
	JwtTokenRefreshExpire = GetValue("jwt.refresh_expire", "")
	DbHost                = GetValue("database.host", "")
	DbPort                = GetValue("database.port", "")
	DbUser                = GetValue("database.user", "")
	DbPassword            = GetValue("database.password", "")
	DbName                = GetValue("database.name", "")
	DbSSLMode             = GetValue("database.sslmode", "")
	DbTimezone            = GetValue("database.timezone", "")
	DbMaxConnections      = GetValue("database.max_connections", "")
	DbIdleConnections     = GetValue("database.max_idle_connections", "")
	AwsUrl                = GetValue("aws_base_url", "")
	IsRunningCron         = GetValue("isRunningCron", "false")
	WhatsAppUrl           = GetValue("whatsappUrl", "")
	WhatsAppToken         = GetValue("whatsappToken", "")
)
