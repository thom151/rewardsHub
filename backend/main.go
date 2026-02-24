package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/thom151/rewardsHub/internal/database"
)

type apiConfig struct {
	db          *database.Queries
	tokenSecret string
}

func main() {

	const port = "8080"
	godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	tokenSecret := os.Getenv("SECRET")
	if tokenSecret == "" {
		log.Fatal("SECRET must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		db:          dbQueries,
		tokenSecret: tokenSecret,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlerReadiness)

	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	//ADMINS ONLY
	mux.Handle("POST /api/platform/create_organization", apiCfg.authMiddleware(apiCfg.plaformAdminOnly(http.HandlerFunc(apiCfg.handlerAdminCreateOrganization))))
	mux.Handle("POST /api/platform/create_service", apiCfg.authMiddleware(apiCfg.plaformAdminOnly(http.HandlerFunc(apiCfg.handlerAdminCreateService))))
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port : %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
