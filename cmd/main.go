package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MartinZitterkopf/gocurse_user/internal/user"
	"github.com/MartinZitterkopf/gocurse_user/pkg/bootstrap"
	"github.com/MartinZitterkopf/gocurse_user/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {

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

	ctx := context.Background()

	userRepo := user.NewRepo(l, instanceDB)
	userService := user.NewService(l, userRepo)
	userEndpoint := user.MakeEndpoints(userService, user.Config{PageLimDefault: pagLimDefault})

	h := handler.NewUserHTTPServer(ctx, userEndpoint)

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	errCh := make(chan error)
	go func() {
		l.Println("listen in ", address)
		errCh <- srv.ListenAndServe()
	}()

	err = <-errCh
	if err != nil {
		log.Fatal(err)
	}
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, HEAD, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
