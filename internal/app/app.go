package app

import (
	"avito-test-task/config"
	"avito-test-task/internal/delivery/handlers"
	mdware "avito-test-task/internal/delivery/middleware"
	"avito-test-task/internal/ports"
	"avito-test-task/internal/repo"
	"avito-test-task/internal/usecase"
	"avito-test-task/pkg"
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"time"
)

func Run(cfg *config.Config) {
	lg, err := pkg.CreateLogger(cfg.LogFile, "prod")
	if err != nil {
		log.Fatal("can't create logger")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password,
		cfg.Host, cfg.Port, cfg.Db.Db)
	pool, err := pgx.Connect(ctx, connString)
	defer pool.Close(context.Background())
	if err != nil {
		log.Fatalf("can't connect to postgresql: %v", err.Error())
	}

	retryAdapter := pkg.NewPostgresRetryAdapter(pool, 3)

	notifyRepo := repo.NewPostgresNotifyRepo(pool)
	notifySender := ports.NewSender()

	done := make(chan bool, 1)
	houseRepo := repo.NewPostgresHouseRepo(pool)
	houseUsecase := usecase.NewHouseUsecase(houseRepo, notifySender, notifyRepo, done,
		5*time.Second, time.Duration(cfg.DbTimeoutSec)*time.Second, lg)
	houseHandler := handlers.NewHouseHandler(houseUsecase, time.Duration(cfg.DbTimeoutSec)*time.Second, lg)

	userRepo := repo.NewPostrgesUserRepo(pool)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handlers.NewUserHandler(userUsecase, time.Duration(cfg.DbTimeoutSec)*time.Second, lg)

	flatRepo := repo.NewPostgresFlatRepo(pool, retryAdapter)
	flatUsecase := usecase.NewFlatUsecase(flatRepo)
	flatHandler := handlers.NewFlatHandler(flatUsecase, time.Duration(cfg.DbTimeoutSec)*time.Second, lg)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Post("/house/create", mdware.AuthMiddleware(mdware.AccessMiddleware(houseHandler.Create)))
	r.Get("/house/{id}", mdware.AuthMiddleware(houseHandler.GetFlatsByID))
	r.Get("/dummyLogin", userHandler.DummyLogin)
	r.Post("/register", userHandler.Register)
	r.Post("/login", userHandler.Login)
	r.Post("/flat/update", mdware.AuthMiddleware(mdware.AccessMiddleware(flatHandler.Update)))
	r.Post("/flat/create", mdware.AuthMiddleware(flatHandler.Create))
	r.Post("/house/{id}/subscribe", mdware.AuthMiddleware(houseHandler.Subscribe))

	fmt.Println("done")
	err = http.ListenAndServe(":8081", r)
	if err != nil {
		fmt.Println(err)
	}
}
