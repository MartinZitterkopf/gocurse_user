package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MartinZitterkopf/gocurse_user/internal/user"
	"github.com/MartinZitterkopf/gocurse_user/pkg/bootstrap"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	router := mux.NewRouter()
	_ = godotenv.Load()
	l := bootstrap.InitLogger()

	instanceDB, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	pagLimDefault := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDefault == "" {
		l.Fatal("paginator limit default is required")
	}

	userRepo := user.NewRepo(l, instanceDB)
	userService := user.NewService(l, userRepo)
	userEndpoint := user.MakeEndpoints(userService, user.Config{PageLimDefault: pagLimDefault})

	router.HandleFunc("/users", userEndpoint.Create).Methods("POST")
	router.HandleFunc("/users", userEndpoint.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEndpoint.GetByID).Methods("GET")
	router.HandleFunc("/users/{id}", userEndpoint.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEndpoint.Delete).Methods("DELETE")

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		Handler:      router,
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
