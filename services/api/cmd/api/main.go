package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"devflowos/api/internal/config"
	"devflowos/api/internal/database"
	"devflowos/api/internal/handler"
	"devflowos/api/internal/middleware"
	"devflowos/api/internal/repository"
	"devflowos/api/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		logger.Error("DATABASE_URL required")
		os.Exit(1)
	}
	if cfg.JWTSecret == "" {
		logger.Error("JWT_SECRET required")
		os.Exit(1)
	}

	ctx := context.Background()
	db, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	ideaRepo := repository.NewIdeaRepository(db)
	leetcodeRepo := repository.NewLeetCodeRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	oppRepo := repository.NewOpportunityRepository(db)
	financeRepo := repository.NewFinanceRepository(db)
	codingLogRepo := repository.NewCodingLogRepository(db)

	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret, 0)
	taskSvc := service.NewTaskService(taskRepo)
	ideaSvc := service.NewIdeaService(ideaRepo)
	leetcodeSvc := service.NewLeetCodeService(leetcodeRepo)
	sessionSvc := service.NewSessionService(sessionRepo)
	oppSvc := service.NewOpportunityService(oppRepo)
	financeSvc := service.NewFinanceService(financeRepo)
	codingLogSvc := service.NewCodingLogService(codingLogRepo)
	aiSvc := service.NewAIContentService(cfg.GeminiAPIKey)

	authH := handler.NewAuthHandler(authSvc)
	taskH := handler.NewTaskHandler(taskSvc)
	ideaH := handler.NewIdeaHandler(ideaSvc)
	leetcodeH := handler.NewLeetCodeHandler(leetcodeSvc)
	sessionH := handler.NewSessionHandler(sessionSvc)
	oppH := handler.NewOpportunityHandler(oppSvc)
	financeH := handler.NewFinanceHandler(financeSvc)
	codingLogH := handler.NewCodingLogHandler(codingLogSvc)
	aiH := handler.NewAIHandler(aiSvc)

	r := chi.NewRouter()
	r.Use(middleware.CORS(cfg.CORSOrigins))
	r.Use(middleware.Logger(logger))

	r.Post("/auth/signup", authH.Signup)
	r.Post("/auth/login", authH.Login)

	r.Route("/", func(r chi.Router) {
		r.Use(middleware.Auth(authSvc))
		r.Get("/tasks/today", taskH.GetToday)
		r.Get("/tasks", taskH.GetByDate)
		r.Post("/tasks", taskH.Create)
		r.Patch("/tasks/{id}", taskH.Update)
		r.Delete("/tasks/{id}", taskH.Delete)
		r.Post("/ideas", ideaH.Create)
		r.Get("/ideas", ideaH.List)
		r.Patch("/ideas/{id}", ideaH.Update)
		r.Post("/leetcode", leetcodeH.Create)
		r.Get("/leetcode", leetcodeH.List)
		r.Post("/sessions/start", sessionH.Start)
		r.Post("/sessions/end", sessionH.End)
		r.Get("/sessions/active", sessionH.GetActive)
		r.Post("/opportunities", oppH.Create)
		r.Get("/opportunities", oppH.List)
		r.Patch("/opportunities/{id}", oppH.Update)
		r.Post("/finances", financeH.Create)
		r.Get("/finances", financeH.List)
		r.Post("/coding-logs", codingLogH.Create)
		r.Get("/coding-logs", codingLogH.List)
		r.Post("/generate-content", aiH.GenerateContent)
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "error", err)
	}
	logger.Info("server stopped")
}
