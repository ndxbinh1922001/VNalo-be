package message

import (
	"log"

	"github.com/gocql/gocql"
	"go.uber.org/fx"

	messageRepo "github.com/ndxbinh1922001/VNalo-be/internal/modules/message/infrastructure/persistence/cassandra"
)

// TODO: Uncomment when handlers are ready
// import "github.com/gin-gonic/gin"

// Module provides message module dependencies
var Module = fx.Options(
	fx.Provide(provideRepository),
	// TODO: Add when service and handlers are implemented
	// fx.Provide(provideService),
	// fx.Provide(provideHandler),
	// fx.Provide(provideWebSocketHandler),
	// fx.Provide(
	// 	fx.Annotate(
	// 		provideRouteRegistration,
	// 		fx.ResultTags(`group:"routes"`),
	// 	),
	// ),
)

func provideRepository(cassandra *gocql.Session) interface{} {
	if cassandra == nil {
		log.Println("⚠️  Message repository not available (Cassandra not connected)")
		return nil
	}

	log.Println("📦 Creating message repository...")
	return messageRepo.NewMessageRepository(cassandra)
}

// TODO: Uncomment when message service is implemented
/*
func provideService(
	repo repository.MessageRepository,
	hub *websocket.Hub,
	redis *redis.Client,
) service.MessageService {
	log.Println("⚙️  Creating message service...")
	return messageService.NewMessageService(repo, hub, redis)
}

func provideHandler(svc service.MessageService) *handler.MessageHandler {
	log.Println("🎯 Creating message handler...")
	return messageHandler.NewMessageHandler(svc)
}

func provideWebSocketHandler(
	hub *websocket.Hub,
	svc service.MessageService,
) *handler.WebSocketHandler {
	log.Println("🔌 Creating WebSocket handler...")
	return messageHandler.NewWebSocketHandler(hub, svc)
}
*/
