package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"log/slog"

	"telemetry-mend-server/internal/db"
	"telemetry-mend-server/internal/handlers"
	"telemetry-mend-server/internal/scm"
	"telemetry-mend-server/internal/ai"
	"telemetry-mend-server/internal/service"
	"telemetry-mend-server/internal/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Initialize Database
	dsn := "file:telemetry_mend.db?cache=shared&_journal_mode=WAL"
	database, err := db.InitDB(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize SCM and Services
	gitClient, err := scm.NewGitClient("./repos")
	if err != nil {
		log.Fatalf("failed to initialize git client: %v", err)
	}
	analyzer := service.NewAnalyzerService(database, gitClient)
	mockAI := &ai.MockProvider{}
	fixer := service.NewFixerService(database, analyzer, mockAI)

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	appHandler := handlers.NewAppHandler(database)
	logHandler := handlers.NewLogHandler(database)
	clusterHandler := handlers.NewClusterHandler(database, fixer)

	r.Route("/api", func(r chi.Router) {
		r.Post("/apps", appHandler.Create)
		r.Get("/apps", appHandler.List)
		r.Post("/logs/ingest", logHandler.Ingest)
		r.Get("/clusters", clusterHandler.List)
		r.Get("/clusters/{id}", clusterHandler.Get)
		r.Post("/clusters/{id}/fix", clusterHandler.GenerateFix)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: web.RegisterHandlers(r),
	}

	go func() {
		slog.Info("starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	slog.Info("server exited")
}
