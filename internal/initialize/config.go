package initialize

import (
	"log"

	"github.com/spf13/viper"

	"github.com/ndxbinh1922001/VNalo-be/internal/infrastructure"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Postgres  PostgresConfig  `mapstructure:"postgres"`
	Cassandra CassandraConfig `mapstructure:"cassandra"`
	Redis     RedisConfig     `mapstructure:"redis"`
	WebSocket WebSocketConfig `mapstructure:"websocket"`
	JWT       JWTConfig       `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release, test
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"sslmode"`
}

type CassandraConfig struct {
	Hosts    []string `mapstructure:"hosts"`
	Keyspace string   `mapstructure:"keyspace"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type WebSocketConfig struct {
	ReadBufferSize  int `mapstructure:"read_buffer_size"`
	WriteBufferSize int `mapstructure:"write_buffer_size"`
	MaxMessageSize  int `mapstructure:"max_message_size"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration int    `mapstructure:"expiration"`
}

// Implement infrastructure.Config interface
func (c *Config) GetPostgresConfig() infrastructure.PostgresConfig {
	return infrastructure.PostgresConfig{
		Host:     c.Postgres.Host,
		Port:     c.Postgres.Port,
		User:     c.Postgres.User,
		Password: c.Postgres.Password,
		Database: c.Postgres.Database,
		SSLMode:  c.Postgres.SSLMode,
	}
}

func (c *Config) GetCassandraConfig() infrastructure.CassandraConfig {
	return infrastructure.CassandraConfig{
		Hosts:    c.Cassandra.Hosts,
		Keyspace: c.Cassandra.Keyspace,
		Username: c.Cassandra.Username,
		Password: c.Cassandra.Password,
	}
}

func (c *Config) GetRedisConfig() infrastructure.RedisConfig {
	return infrastructure.RedisConfig{
		Host:     c.Redis.Host,
		Port:     c.Redis.Port,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	}
}

func (c *Config) GetServerMode() string {
	return c.Server.Mode
}

// LoadConfig loads configuration from file and environment
func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Set defaults for Server
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")

	// Set defaults for PostgreSQL
	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", "5432")
	viper.SetDefault("postgres.user", "postgres")
	viper.SetDefault("postgres.password", "postgres")
	viper.SetDefault("postgres.database", "vnalo_db")
	viper.SetDefault("postgres.sslmode", "disable")

	// Set defaults for Cassandra
	viper.SetDefault("cassandra.hosts", []string{"localhost"})
	viper.SetDefault("cassandra.keyspace", "vnalo_chat")
	viper.SetDefault("cassandra.username", "")
	viper.SetDefault("cassandra.password", "")

	// Set defaults for Redis
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// Set defaults for WebSocket
	viper.SetDefault("websocket.read_buffer_size", 1024)
	viper.SetDefault("websocket.write_buffer_size", 1024)
	viper.SetDefault("websocket.max_message_size", 524288) // 512KB

	// Set defaults for JWT
	viper.SetDefault("jwt.secret", "change-this-secret-key-in-production")
	viper.SetDefault("jwt.expiration", 86400) // 24 hours

	// Enable reading from environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("⚠️  Warning: Error reading config file: %v. Using defaults and environment variables.", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("❌ Unable to decode config into struct: %v", err)
	}

	log.Println("✅ Configuration loaded successfully")
	return &config
}
