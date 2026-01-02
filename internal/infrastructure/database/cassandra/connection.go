package cassandra

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

type Config struct {
	Hosts    []string
	Keyspace string
	Username string
	Password string
}

var Session *gocql.Session

// NewSession creates a new Cassandra session
func NewSession(cfg Config) (*gocql.Session, error) {
	cluster := gocql.NewCluster(cfg.Hosts...)
	cluster.Keyspace = cfg.Keyspace
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Timeout = 10 * time.Second

	if cfg.Username != "" && cfg.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: cfg.Username,
			Password: cfg.Password,
		}
	}

	// Connection pooling
	cluster.NumConns = 2
	cluster.SocketKeepalive = 10 * time.Second

	// Retry policy
	cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{
		NumRetries: 3,
		Min:        time.Second,
		Max:        10 * time.Second,
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Cassandra: %w", err)
	}

	Session = session
	log.Println("Cassandra session created successfully")
	return session, nil
}

// GetSession returns the global session instance
func GetSession() *gocql.Session {
	if Session == nil {
		log.Fatal("Cassandra session not initialized. Call NewSession() first.")
	}
	return Session
}

// Close closes the Cassandra session
func Close() {
	if Session != nil {
		Session.Close()
		log.Println("Cassandra session closed")
	}
}

