package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/thom151/rewardsHub/internal/database"
	"github.com/thom151/rewardsHub/internal/dropbox"
)

type apiConfig struct {
	db                       *database.Queries
	tokenSecret              string
	dropboxAccToken          string
	dropboxAccTokenExpiresAt time.Time
	dropboxRefreshToken      string
	dropboxClientID          string
	dropboxClientSecret      string
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

	dropboxRefreshToken := os.Getenv("DROPBOX_REFRESH_TOKEN")
	if dropboxRefreshToken == "" {
		log.Fatal("DROPBOX_REFRESH_TOKEN environment variable is not set")
	}

	dropboxClientID := os.Getenv("DROPBOX_CLIENT_ID")
	if dropboxClientID == "" {
		log.Fatal("DROPBOX_CLIENT_ID environment variable is not set")
	}

	dropboxClientSecret := os.Getenv("DROPBOX_CLIENT_SECRET")
	if dropboxClientSecret == "" {
		log.Fatal("DROPBOX_CLIENT_SECRET environment variable is not set")
	}

	dropboxAccToken, err := dropbox.GetNewAccessToken(dropboxRefreshToken, dropboxClientID, dropboxClientSecret)
	if err != nil {
		log.Fatal("DROPBOX_ACC_TOKEN cannot get")
	}

	fmt.Printf("ACC TOKEN: %s\n", dropboxAccToken.AccessToken)

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		db:                       dbQueries,
		tokenSecret:              tokenSecret,
		dropboxAccToken:          dropboxAccToken.AccessToken,
		dropboxAccTokenExpiresAt: time.Now().Add(time.Duration(dropboxAccToken.ExpiresIn) * time.Second),
		dropboxRefreshToken:      dropboxRefreshToken,
		dropboxClientID:          dropboxClientID,
		dropboxClientSecret:      dropboxClientSecret,
	}
	mux := http.NewServeMux()

	const filepathRoot = "./web"
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	mux.HandleFunc("GET /healthz", handlerReadiness)

	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	//AUTHORIZED USERS
	mux.Handle("POST /api/property", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerCreateProperty)))
	mux.Handle("POST /api/booking", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerCreateBooking)))

	//ADMINS ONLY
	mux.Handle("POST /api/platform/organization", apiCfg.authMiddleware(apiCfg.plaformAdminOnly(http.HandlerFunc(apiCfg.handlerAdminCreateOrganization))))
	mux.Handle("POST /api/platform/service", apiCfg.authMiddleware(apiCfg.plaformAdminOnly(http.HandlerFunc(apiCfg.handlerAdminCreateService))))
	mux.Handle("POST /api/platform/bookings/{booking_id}/confirm", apiCfg.authMiddleware(apiCfg.plaformAdminOnly(http.HandlerFunc(apiCfg.handlerAdminConfirmBooking))))

	//ADMIN OR ORG ADMIN
	mux.Handle("POST /api/org/approve_membership", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerApproveOrganizationMembership)))
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
