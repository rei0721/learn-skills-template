package app

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/rei0721/go-scaffold/internal/handler"
	"github.com/rei0721/go-scaffold/internal/middleware"
	"github.com/rei0721/go-scaffold/internal/repository"
	"github.com/rei0721/go-scaffold/internal/router"
	"github.com/rei0721/go-scaffold/internal/service"
	authService "github.com/rei0721/go-scaffold/internal/service/auth"
	rbacService "github.com/rei0721/go-scaffold/internal/service/rbac"
	"github.com/rei0721/go-scaffold/pkg/dbtx"
)

func (app *App) initBusiness() error {
	authRepo := repository.NewAuthRepository(app.DB.DB())
	authSvc := authService.NewAuthService(authRepo)
	authSvcAsService, ok := authSvc.(service.Service)
	if !ok {
		return fmt.Errorf("auth service does not implement service.Service")
	}
	if _, err := app.setServiceAll(authSvcAsService); err != nil {
		return err
	}
	authHandler := handler.NewAuthHandler(authSvc, app.Logger)

	var rbacHandler *handler.RBACHandler
	if app.RBAC != nil {
		rbacRepo := repository.NewRBACRepository(app.RBAC)
		rbacSvc := rbacService.NewRBACService(rbacRepo)
		if _, err := app.setServiceAll(rbacSvc); err != nil {
			return err
		}
		rbacHandler = handler.NewRBACHandler(rbacSvc, app.Logger)
	}

	// 初始化 router
	r := router.New(authHandler, rbacHandler, app.Logger, app.JWT)

	// Set Gin mode based on config
	if app.Config.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if app.Config.Server.Mode == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Setup router with middleware
	middlewareCfg := middleware.DefaultMiddlewareConfig()
	app.Router = r.Setup(middlewareCfg)

	return nil
}

// 注入 app 到 Service 层
func (app *App) setServiceAll(services ...service.Service) (*App, error) {
	for _, s := range services {
		if app.Executor != nil {
			s.SetExecutor(app.Executor)
			app.Logger.Debug("executor injected into user service")
		}

		// ⭐ 延迟注入 cache 到 Service 层
		if app.Cache != nil {
			s.SetCache(app.Cache)
			app.Logger.Debug("cache injected into user service")
		}

		// ⭐ 延迟注入 logger 到 Service 层
		if app.Logger != nil {
			s.SetLogger(app.Logger)
			app.Logger.Debug("logger injected into user service")
		}

		// ⭐ 延迟注入 JWT 到 Service 层
		if app.JWT != nil {
			s.SetJWT(app.JWT)
			app.Logger.Debug("JWT injected into user service")
		}

		// ⭐ 延迟注入 RBAC 到 Service 层
		if app.RBAC != nil {
			s.SetRBAC(app.RBAC)
			app.Logger.Debug("RBAC injected into user service")
		}

		// ⭐ 延迟注入 IDGenerator 到 Service 层
		if app.IDGenerator != nil {
			s.SetIDGenerator(app.IDGenerator)
			app.Logger.Debug("IDGenerator injected into user service")
		}

		// ⭐ 延迟注入 Crypto 到 Service 层
		if app.Crypto != nil {
			s.SetCrypto(app.Crypto)
			app.Logger.Debug("Crypto injected into service")
		}

		// ⭐ 延迟注入 Database 到 Service 层
		if app.DB != nil {
			s.SetDB(app.DB)
			app.Logger.Debug("Database injected into service")
		}

		// ⭐ 创建并注入 TxManager 到 Service 层
		if app.DB != nil {
			// 为每个Service创建事务管理器
			txManager, err := dbtx.NewManager(app.DB.DB(), app.Logger)
			if err != nil {
				app.Logger.Error("failed to create TxManager", "error", err)
			} else {
				s.SetTxManager(txManager)
				app.Logger.Debug("TxManager injected into service")
			}
		}
	}
	return app, nil
}
