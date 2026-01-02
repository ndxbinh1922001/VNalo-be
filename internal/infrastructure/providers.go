package infrastructure

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ndxbinh1922001/VNalo-be/internal/infrastructure/cache"
	"github.com/ndxbinh1922001/VNalo-be/internal/infrastructure/database/cassandra"
	"github.com/ndxbinh1922001/VNalo-be/internal/infrastructure/websocket"
)

// Config interface for infrastructure providers
type Config interface {
	GetPostgresConfig() PostgresConfig
	GetCassandraConfig() CassandraConfig
	GetRedisConfig() RedisConfig
	GetServerMode() string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

type CassandraConfig struct {
	Hosts    []string
	Keyspace string
	Username string
	Password string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// Module provides all infrastructure dependencies
var Module = fx.Options(
	// Databases
	fx.Provide(ProvidePostgresGORM),
	fx.Provide(ProvideCassandra),
	fx.Provide(ProvideRedis),

	// WebSocket
	fx.Provide(ProvideWebSocketHub),
	fx.Invoke(runWebSocketHub),
)

// ProvidePostgresGORM provides PostgreSQL GORM connection
func ProvidePostgresGORM(cfg Config) (*gorm.DB, error) {
	log.Println("🗄️  Initializing PostgreSQL (GORM)...")
	
	pgCfg := cfg.GetPostgresConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pgCfg.Host,
		pgCfg.Port,
		pgCfg.User,
		pgCfg.Password,
		pgCfg.Database,
		pgCfg.SSLMode,
	)

	var logLevel logger.LogLevel
	if cfg.GetServerMode() == "release" {
		logLevel = logger.Silent
	} else {
		logLevel = logger.Info
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PostgreSQL with GORM: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("✅ PostgreSQL (GORM) connected")
	return gormDB, nil
}

// ProvideCassandra provides Cassandra session
func ProvideCassandra(cfg Config) *gocql.Session {
	log.Println("🗄️  Initializing Cassandra...")
	
	cassCfg := cfg.GetCassandraConfig()
	cassSession, err := cassandra.NewSession(cassandra.Config{
		Hosts:    cassCfg.Hosts,
		Keyspace: cassCfg.Keyspace,
		Username: cassCfg.Username,
		Password: cassCfg.Password,
	})
	if err != nil {
		log.Printf("⚠️  Cassandra connection failed: %v (message features will not work)", err)
		return nil
	}

	log.Println("✅ Cassandra connected")
	return cassSession
}

// ProvideRedis provides Redis client
func ProvideRedis(cfg Config) *redis.Client {
	log.Println("🗄️  Initializing Redis...")
	
	redisCfg := cfg.GetRedisConfig()
	redisClient, err := cache.NewRedisClient(cache.Config{
		Host:     redisCfg.Host,
		Port:     redisCfg.Port,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})
	if err != nil {
		log.Printf("⚠️  Redis connection failed: %v (cache and presence features will not work)", err)
		return nil
	}

	log.Println("✅ Redis connected")
	return redisClient
}

// ProvideWebSocketHub provides WebSocket Hub
func ProvideWebSocketHub() *websocket.Hub {
	log.Println("🔌 Creating WebSocket Hub...")
	return websocket.NewHub()
}

func runWebSocketHub(hub *websocket.Hub) {
	log.Println("▶️  Starting WebSocket Hub...")
	go hub.Run()
}

