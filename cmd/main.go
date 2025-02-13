package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "rashpor.com/todolist/docs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	"rashpor.com/todolist/internal/config"
	"rashpor.com/todolist/internal/handler"
	"rashpor.com/todolist/internal/storage"
)

// @title TODO API
// @version 1.0
// @description API для управления задачами
// @host localhost:8083
// @BasePath /
func main() {

	sLogger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	config.LoadEnvironment()
	appConfig := config.NewDBConfig()

	db, err := storage.ConnectDB(appConfig)
	if err != nil {
		sLogger.Error("Ошибка подключения к БД", slog.Any("error", err))
		return
	}
	defer db.Close()

	runMigrations(db, sLogger)

	repo := storage.NewTaskRepository(db, sLogger)
	taskHandler := handler.NewTaskHandler(repo, sLogger)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(compress.New())

	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	app.Post("/tasks", taskHandler.CreateTask)
	app.Get("/tasks", taskHandler.GetAllTasks)
	app.Put("/tasks/:id", taskHandler.UpdateTask)
	app.Delete("/tasks/:id", taskHandler.DeleteTask)

	sLogger.Info("Запускаем сервер на порту :8083...")
	log.Fatal(app.Listen(":8083"))
}

func runMigrations(db *pgxpool.Pool, logger *slog.Logger) {
	logger.Info("Начинаем миграции")
	sqlDB, err := sql.Open("pgx", db.Config().ConnString())
	if err != nil {
		logger.Error("Ошибка создания *sql.DB из *pgxpool.Pool", slog.Any("error", err))
		return
	}
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		logger.Error("Ошибка создания драйвера миграции", slog.Any("error", err))
		return
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/storage/migrations",
		"postgres", driver,
	)
	if err != nil {
		logger.Error("Ошибка инициализации миграций", slog.Any("error", err))
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("Ошибка применения миграций", slog.Any("error", err))
		return
	}

	logger.Info("Миграции успешно применены")
}
