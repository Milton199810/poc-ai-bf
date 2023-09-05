package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alemelomeza/improved-octo-memory.git/internal/infra/services/llm"
	"github.com/alemelomeza/improved-octo-memory.git/internal/infra/web"
	custommiddleware "github.com/alemelomeza/improved-octo-memory.git/internal/infra/web/middleware"
	"github.com/alemelomeza/improved-octo-memory.git/internal/repository"
	"github.com/alemelomeza/improved-octo-memory.git/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

func main() {
	// Environment variables
	viper.SetConfigName(".env.local")
	viper.SetConfigType("env")
	viper.AddConfigPath("config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// DB connections
	dns := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%s;database=%s",
		viper.GetString("SQL_SERVER_HOST"),
		viper.GetString("SQL_SERVER_USER"),
		viper.GetString("SQL_SERVER_PASSWORD"),
		viper.GetString("SQL_SERVER_PORT"),
		viper.GetString("SQL_SERVER_DATABASE"),
	)
	db, err := sql.Open("mssql", dns)
	if err != nil {
		log.Fatalf("DB startup error: %v\n", err)
	}
	defer db.Close()

	// Repositories
	conversationRepo := repository.NewConversationRepositorySQLServer(db)
	summaryRepo := repository.NewSummaryRepositoryCSV("../storage/summaries.csv")

	// LLMs
	AzureSummaryLLM := llm.NewAzureOpenAI(
		viper.GetString("AZURE_OPENAI_KEY"),
		viper.GetString("AZURE_OPENAI_ENDPOINT"),
	)
	GCPSummaryLLM := llm.NewGCPVertex(
		viper.GetString("GCP_CREDENTIALS_PATH"),
		viper.GetString("GCP_VERTEX_API_ENDPOINT"),
		viper.GetString("GCP_VERTEX_PROJECT_ID"),
		viper.GetString("GCP_VERTEX_MODEL_ID"),
	)

	// Use Cases
	summaryUseCase := usecase.NewSummaryUseCase(conversationRepo, AzureSummaryLLM, GCPSummaryLLM)
	evaluateUseCase := usecase.NewEvaluationUseCase(summaryRepo)

	// Handlers/Controllers
	handlers := web.NewHandlers(summaryUseCase, evaluateUseCase)

	// Routes
	r := chi.NewRouter()

	r.Use(custommiddleware.BasicAuthMiddleware)

	r.Get("/summary", handlers.GetSummaryHandler)
	r.Post("/summary", handlers.PostSummaryHandler)
	r.Post("/evaluation", handlers.PostEvaluationHandler)

	// Server
	server := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server startup error: %v\n", err)
	}

	// Gracefull shutdown
	shutdownCtx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		s := <-sigCh
		log.Printf("Received signal %v, attempting graceful shutdown", s)
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("HTTP shutdown error: %v\n", err)
		}
		log.Println("Graceful shutdown complete")
		cancelCtx()
	}()
}
