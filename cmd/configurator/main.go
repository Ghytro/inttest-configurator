package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"

	_ "github.com/lib/pq"

	"configurator/internal/api"

	authApi "configurator/internal/api/auth"
	projectsApi "configurator/internal/api/projects"

	authRepository "configurator/internal/repository/auth"
	projectsRepository "configurator/internal/repository/projects"

	authUseCase "configurator/internal/usecase/auth"
	projectsUseCase "configurator/internal/usecase/projects"

	"configurator/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/pressly/goose"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
)

// @title IntTest configurator
// @version 2.0
// @description idk what to write here
// @description it's just a swagger
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1

// @securityDefinitions.cookieAuth ApiKeyAuth
// @in cookie
// @name jwt

func main() {
	// fixme: по нормальному через кобру переделать
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	sigIntListener(g)

	logger, err := initLogger()
	if err != nil {
		log.Fatal(errors.Wrap(err, "error in logger init"))
	}

	db, err := database.NewPostgres(
		ctx,
		os.Getenv("DATABASE_URL"), // separate the config
		logger,
		database.WithPgPoolSize(8),
	)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error connecting to db"))
	}
	if os.Getenv("AUTO_MIGRATE") == "true" {
		if err := autoMigrate(os.Getenv("DATABASE_URL"), "up"); err != nil {
			log.Fatal(errors.Wrap(err, "error while migrating"))
		}
	}

	var (
		apis     []api.Handler
		aUseCase *authUseCase.UseCase
	)
	authRepo := authRepository.New(db, logger)
	projectRepo := projectsRepository.New(db, logger)

	{
		aUseCase = authUseCase.New(logger, authRepo)
		apis = append(apis, authApi.New(logger, aUseCase))
	}
	{
		useCase := projectsUseCase.New(projectRepo, logger)
		apis = append(apis, projectsApi.New(useCase, logger))
	}
	app := initFiberRouter(apis, aUseCase)
	ln, err := initFiberListener(":8080")
	if err != nil {
		log.Fatal(err)
	}
	listenFiber(ctx, g, app, ln)
}

func sigIntListener(g *errgroup.Group) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	g.Go(func() error {
		for range sigs {
			return errors.New("sigint caught, graceful shutdown")
		}
		return nil
	})
}

func initLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment(zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func initFiberRouter(
	apis []api.Handler,
	// для валидатора через постгрес
	uc api.AuthMiddlewareUseCase,
) *fiber.App {
	fiberApi := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
		ErrorHandler: api.ErrorHandler,

		BodyLimit:                100 * 1024 * 1024, // 300mb
		EnableSplittingOnParsers: true,
	})
	apiV1 := fiberApi.Group("/api/v1")
	apiV1.Use(logger.New())
	authenticator := api.NewAuthenticator(uc)
	for _, a := range apis {
		a.Register(apiV1, authenticator)
	}
	return fiberApi
}

func initFiberListener(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}

func listenFiber(ctx context.Context, g *errgroup.Group, app *fiber.App, ln net.Listener) {
	g.Go(func() error {
		return app.Listener(ln)
	})
	<-ctx.Done()
	ln.Close()
}

const (
	pgDialect       = "postgres"
	sqlMigrationDir = "./internal/migrations"
)

func autoMigrate(dbUrl, cmd string) error {
	if err := goose.SetDialect(pgDialect); err != nil {
		return err
	}

	db, err := sql.Open(pgDialect, dbUrl)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic("failed to close goose dbconn")
		}
	}()
	if err := goose.Run(cmd, db, sqlMigrationDir); err != nil {
		return err
	}
	return nil
}
