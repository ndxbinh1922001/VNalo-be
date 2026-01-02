# 🎉 VNalo Chat - Telegram-like Chat Application

Real-time chat application với **Uber Fx**, **GORM**, **Goose**, **PostgreSQL**, **Cassandra**, **Redis**, **WebSocket** theo **DDD + Clean Architecture**.

[![Build](https://img.shields.io/badge/build-passing-brightgreen)]()
[![Go](https://img.shields.io/badge/Go-1.25.3-blue)]()
[![Uber Fx](https://img.shields.io/badge/Uber%20Fx-1.23-orange)]()
[![GORM](https://img.shields.io/badge/GORM-1.31.1-red)]()

---

## ✨ **Đặc Điểm Nổi Bật**

### 🏗️ Enterprise Architecture
- ✅ **Uber Fx** - Professional dependency injection với auto-wiring
- ✅ **Auto Route Registration** - Zero router updates khi thêm module mới
- ✅ **Self-Contained Modules** - Mỗi module tự quản lý providers
- ✅ **DDD** - Domain-Driven Design
- ✅ **Clean Architecture** - 4 layers separation
- ✅ **Polyglot Persistence** - PostgreSQL + Cassandra

### 🚀 Production-Ready
- ✅ **GORM + PostgreSQL** - Relational data
- ✅ **Goose** - Database migrations
- ✅ **Cassandra** - High-throughput messages
- ✅ **Redis** - Caching & Pub/Sub
- ✅ **WebSocket** - Real-time communication
- ✅ **Docker Compose** - Easy deployment

---

## 📊 **Tech Stack**

```
Backend:        Go 1.25.3
DI Framework:   Uber Fx 1.23.0
Web Framework:  Gin 1.11.0
ORM:            GORM 1.31.1
Migration:      Goose 3.26.0
Databases:      PostgreSQL, Cassandra, Redis
Real-time:      Gorilla WebSocket 1.5.3
Auth:           JWT 5.3.0
API Docs:       Swagger
Container:      Docker Compose
```

---

## 🏗️ **Architecture**

### Modular Architecture với Auto Route Registration
```
┌─────────────────────────────────────────────────────────┐
│              Go API Server (Gin + WebSocket)            │
│              Uber Fx Auto-Wiring + Auto Routes          │
└─────────────────────────────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
┌───────▼──────┐  ┌───────▼──────┐  ┌──────▼──────┐
│ Infrastructure│  │   Modules     │  │   Router    │
│   Module      │  │  (Self-Cont.) │  │ (Auto-Reg.) │
│               │  │               │  │             │
│ • PostgreSQL  │  │ • User        │  │ • Collects  │
│ • Cassandra   │  │ • Message     │  │   all route │
│ • Redis       │  │ • Conversation│  │   functions │
│ • WebSocket   │  │ • Contact     │  │ • Auto-reg. │
└──────────────┘  └───────────────┘  └─────────────┘
        │                 │                 │
    ┌───┴───┬──────────────┴──────────────┬─┴───┐
    │       │              │               │     │
┌───▼───┐ ┌─▼──────┐  ┌───▼────┐  ┌──────▼──┐ ┌─▼───┐
│PgSQL  │ │Cassandra│  │ Redis  │  │WebSocket│ │MinIO│
│(GORM) │ │ (gocql) │  │(cache) │  │  Hub    │ │(S3) │
└───────┘ └─────────┘  └────────┘  └─────────┘ └─────┘
```

### Auto Route Registration Flow
```
Module provides route function
    ↓
Fx tags với group:"routes"
    ↓
Router auto-collects tất cả
    ↓
Routes registered tự động! ✨
```

---

## 🚀 **Quick Start**

### 1. Start Infrastructure
```bash
cd /Users/lap14945/Desktop/VNalo-be
make docker-up
sleep 60
```

### 2. Initialize Databases
```bash
# Cassandra
docker cp migrations/cassandra/schema.cql vnalo_cassandra:/schema.cql
docker exec -it vnalo_cassandra cqlsh -f /schema.cql

# PostgreSQL (auto-runs on app start)
make migrate-up
```

### 3. Run Application
```bash
make run
```

### 4. Test
```bash
curl http://localhost:8080/health
open http://localhost:8080/swagger/index.html
```

---

## 📂 **Project Structure**

```
VNalo-be/
├── cmd/                      # Entry points
│   ├── api/main.go          # Fx app (35 lines!)
│   └── migrate/main.go      # Migration CLI
│
├── internal/
│   ├── infrastructure/      # Infrastructure Module (Self-Contained)
│   │   ├── providers.go    # Infrastructure providers (DB, Cache, WS)
│   │   ├── database/       # PostgreSQL, Cassandra
│   │   ├── cache/          # Redis
│   │   └── websocket/      # WebSocket Hub & Client
│   │
│   ├── initialize/         # Uber Fx App Module
│   │   ├── providers.go    # Main module registry
│   │   ├── config.go       # Config provider
│   │   ├── migration.go    # Migration provider
│   │   └── router.go       # Router với Auto Route Registration
│   │
│   ├── modules/            # Business Modules (DDD - Self-Contained)
│   │   ├── user/          # ✅ Complete (100%)
│   │   │   ├── providers.go        # Module providers + route registration
│   │   │   ├── domain/             # Entities, Value Objects, Repositories
│   │   │   ├── application/        # Services, DTOs
│   │   │   └── presentation/       # Handlers, Routers
│   │   └── message/       # ⚠️ Backend ready (70%)
│   │       ├── providers.go        # Module providers
│   │       ├── domain/
│   │       ├── application/
│   │       └── presentation/
│   │
│   └── middleware/        # HTTP middlewares
│
├── migrations/            # Database migrations
│   ├── postgres/         # Goose migrations
│   └── cassandra/        # CQL schema
│
├── pkg/                  # Shared utilities
├── config/               # Configuration
└── docs/                 # Documentation
```

### 🎯 **Key Architecture Features**

**1. Self-Contained Modules**
- Mỗi module có `providers.go` riêng
- Module tự quản lý dependencies
- Module tự register routes

**2. Auto Route Registration**
- Router tự động collect tất cả route functions
- Zero router updates khi thêm module mới
- Sử dụng Uber Fx function groups

**3. Infrastructure Module**
- Tách riêng infrastructure concerns
- PostgreSQL, Cassandra, Redis, WebSocket
- Có `providers.go` riêng

---

## 🎯 **Features**

### ✅ Working Now
- User Management (CRUD, VIP, Password)
- Health Check API
- Swagger Documentation
- Database Connections (All 3)
- WebSocket Hub (Running)
- Migrations (Goose)
- Docker Compose Setup

### ⏳ In Development
- Message sending/receiving
- WebSocket endpoints
- JWT Authentication
- Conversation management
- Contact management

---

## 🛠️ **Commands**

```bash
# Development
make run              # Run application
make build            # Build binaries
make test             # Run tests

# Database
make migrate-up       # Run migrations
make migrate-down     # Rollback
make migrate-status   # Check status

# Docker
make docker-up        # Start services
make docker-down      # Stop services
make postgres-shell   # PostgreSQL CLI
make cassandra-shell  # Cassandra CLI

# Utilities
make swagger          # Generate API docs
make help             # Show all commands
```

---

## 📚 **Documentation**

### Getting Started
- 📖 `IMPLEMENTATION_COMPLETE.txt` - Quick summary
- 📖 `UBER_FX_IMPLEMENTATION.md` - Complete Fx guide
- 📖 `docs/UBER_FX_GUIDE.md` - Quick reference

### For Developers
- 📖 `PROJECT_STRUCTURE.md` - Visual structure
- 📖 `SUCCESS_SUMMARY.txt` - Statistics

---

## 🎓 **Adding New Module**

Với **Auto Route Registration**, thêm module mới chỉ cần **2 bước**!

### Step 1: Create Module với Providers (Self-Contained)
```go
// internal/modules/conversation/providers.go
package conversation

import (
	"log"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
	
	// ... imports
)

// Module provides conversation module dependencies
var Module = fx.Options(
	// Dependencies
	fx.Provide(provideRepository),
	fx.Provide(provideService),
	fx.Provide(provideHandler),
	
	// Route Registration (Auto-collected by router!)
	fx.Provide(
		fx.Annotate(
			provideRouteRegistration,
			fx.ResultTags(`group:"routes"`), // ← Tag để router auto-collect
		),
	),
)

func provideRepository(db *gorm.DB) repository.ConversationRepository {
	return repo.NewConversationRepository(db)
}

func provideService(repo repository.ConversationRepository) service.ConversationService {
	return svc.NewConversationService(repo)
}

func provideHandler(svc service.ConversationService) *handler.ConversationHandler {
	return handler.NewConversationHandler(svc)
}

// Route registration function (auto-called by router)
func provideRouteRegistration(h *handler.ConversationHandler) func(*gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		log.Println("✅ Registering conversation routes...")
		conversationRouter.RegisterRoutes(router, h)
	}
}
```

### Step 2: Register Module (1 line)
```go
// internal/initialize/providers.go
import "github.com/.../conversation"

var AppModule = fx.Options(
	infrastructure.Module,
	user.Module,
	message.Module,
	conversation.Module, // ← Add here!
	RouterModule,
)
```

### ✅ **Done!**
```
✅ router.go: NO CHANGES NEEDED!
✅ Routes auto-registered!
✅ Fx auto-wires everything!
```

**Benefits:**
- 🎯 **Zero Router Updates** - Router tự động collect routes
- 🎯 **Self-Contained** - Module tự quản lý mọi thứ
- 🎯 **Type-Safe** - Compile-time checks
- 🎯 **Scalable** - Thêm bao nhiêu module cũng được!

---

## 📊 **Progress**

```
Infrastructure:      ████████████████████ 100% ✅
User Module:         ████████████████████ 100% ✅
Message Module:      ██████████████░░░░░░  70% ⚠️
Auto Route Reg:      ████████████████████ 100% ✅
Documentation:       ████████████████████ 100% ✅
Uber Fx:             ████████████████████ 100% ✅

OVERALL:             ████████████████░░░░  85% ✅
```

### 🎯 **Architecture Improvements**
- ✅ **Self-Contained Modules** - Mỗi module có providers.go riêng
- ✅ **Infrastructure Module** - Tách riêng infrastructure concerns
- ✅ **Auto Route Registration** - Zero router updates
- ✅ **Function Groups** - Uber Fx advanced features

---

## 🔐 **API Endpoints**

### User Management
```
POST   /api/v1/users                   Create user
GET    /api/v1/users                   List users
GET    /api/v1/users/:id               Get user
PUT    /api/v1/users/:id               Update user
DELETE /api/v1/users/:id               Delete user
POST   /api/v1/users/:id/promote-vip   Promote to VIP
POST   /api/v1/users/:id/demote-vip    Demote from VIP
POST   /api/v1/users/:id/change-password  Change password
POST   /api/v1/users/:id/activate      Activate
POST   /api/v1/users/:id/deactivate    Deactivate
```

### System
```
GET    /health                         Health check
GET    /swagger/*any                   API documentation
```

---

## 🧪 **Testing**

```bash
# Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@vnalo.com",
    "password": "password123",
    "username": "testuser"
  }'

# Get users
curl http://localhost:8080/api/v1/users

# Swagger UI
open http://localhost:8080/swagger/index.html
```

---
**Made with ❤️ using Go, Uber Fx, GORM, Goose, PostgreSQL, Cassandra, Redis & WebSocket**

🚀 **Happy Coding!** 🚀
