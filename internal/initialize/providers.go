package initialize

import (
	"go.uber.org/fx"

	"github.com/ndxbinh1922001/VNalo-be/internal/infrastructure"
	"github.com/ndxbinh1922001/VNalo-be/internal/modules/message"
	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user"
)

// AppModule combines all application modules
var AppModule = fx.Options(
	// Configuration
	fx.Provide(LoadConfig),
	fx.Provide(AsInfrastructureConfig),

	// Infrastructure (Databases, WebSocket, Cache)
	infrastructure.Module,

	// Migrations
	MigrationModule,

	// Business modules (from internal/modules)
	// TODO: Add more modules as they are implemented
	user.Module,
	message.Module,

	// Router (must be last)
	RouterModule,
)

// AsInfrastructureConfig converts *Config to infrastructure.Config interface
func AsInfrastructureConfig(cfg *Config) infrastructure.Config {
	return cfg
}
