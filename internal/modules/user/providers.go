package user

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/application/service"
	userService "github.com/ndxbinh1922001/VNalo-be/internal/modules/user/application/service"
	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/domain/repository"
	userRepo "github.com/ndxbinh1922001/VNalo-be/internal/modules/user/infrastructure/persistence/repository"
	userHandler "github.com/ndxbinh1922001/VNalo-be/internal/modules/user/presentation/http/handler"
	userRouter "github.com/ndxbinh1922001/VNalo-be/internal/modules/user/presentation/http/router"
)

// Module provides user module dependencies
var Module = fx.Options(
	fx.Provide(provideRepository),
	fx.Provide(provideService),
	fx.Provide(provideHandler),
	fx.Provide(
		fx.Annotate(
			provideRouteRegistration,
			fx.ResultTags(`group:"routes"`), // ← Tag để Fx auto-collect
		),
	),
)

func provideRepository(db *gorm.DB) repository.UserRepository {
	log.Println("📦 Creating user repository...")
	return userRepo.NewUserRepository(db)
}

func provideService(repo repository.UserRepository) service.UserService {
	log.Println("⚙️  Creating user service...")
	return userService.NewUserService(repo)
}

func provideHandler(svc service.UserService) *userHandler.UserHandler {
	log.Println("🎯 Creating user handler...")
	return userHandler.NewUserHandler(svc)
}

// provideRouteRegistration returns a function to register user routes
// Fx will collect this function và router sẽ tự động gọi nó! ✨
func provideRouteRegistration(h *userHandler.UserHandler) func(*gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		log.Println("✅ Registering user routes...")
		// Dùng trực tiếp function RegisterUserRoutes có sẵn! ✨
		userRouter.RegisterUserRoutes(router, h)
	}
}
