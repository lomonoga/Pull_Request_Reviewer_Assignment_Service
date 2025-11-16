package router

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"pull_requests_service/internal/config"
	"pull_requests_service/internal/handler"
	"pull_requests_service/internal/repository"
	"pull_requests_service/internal/service"
	"time"
)

func SetupRouter(cfg *config.Config) http.Handler {

	db, err := sql.Open("postgres", cfg.GetDBConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	db.SetMaxOpenConns(cfg.DBMaxConns)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	teamRepo := repository.NewTeamRepository(db)
	userRepo := repository.NewUserRepository(db)
	prRepo := repository.NewPRRepository(db)

	teamService := service.NewTeamService(teamRepo, userRepo)
	userService := service.NewUserService(userRepo, prRepo)
	prService := service.NewPRService(prRepo, userRepo)

	teamHandler := handler.NewTeamHandler(teamService)
	userHandler := handler.NewUserHandler(userService)
	prHandler := handler.NewPRHandler(prService)

	mux := http.NewServeMux()

	// Team
	mux.HandleFunc("POST /team/add", teamHandler.AddTeam)
	mux.HandleFunc("GET /team/get", teamHandler.GetTeam)

	// User
	mux.HandleFunc("POST /users/setIsActive", userHandler.SetUserActive)
	mux.HandleFunc("GET /users/getReview", userHandler.GetUserReviews)

	// PR
	mux.HandleFunc("POST /pullRequest/create", prHandler.CreatePR)
	mux.HandleFunc("POST /pullRequest/merge", prHandler.MergePR)
	mux.HandleFunc("POST /pullRequest/reassign", prHandler.ReassignPR)

	// Health check
	mux.HandleFunc("GET /health", healthHandler)

	handler := applyMiddleware(mux)

	return handler
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "healthy"}`))
}

func applyMiddleware(handler http.Handler) http.Handler {
	handler = loggingMiddleware(handler)
	return handler
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		for key, values := range recorder.Header() {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())

		duration := time.Since(start)
		log.Printf(
			"[%s] %s %s - %d %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			recorder.Code,
			duration,
		)
	})
}
