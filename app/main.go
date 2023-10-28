package main

import (
	"errors"
	"fmt"
	stdlog "log"
	"net/http"
	"os"

	"github.com/MauricioMilano/stock_app/config"
	"github.com/MauricioMilano/stock_app/routes"
	"github.com/MauricioMilano/stock_app/services/rabbitmq"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const (
	i = 7
	j
	k
)

func prin() {
	fmt.Print("Hello")
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}

}

func run() error {
	// Load env values
	err := godotenv.Load()
	if err != nil {
		stdlog.Println("Error loading .env file")
		return err
	}

	// Connect Rabbit MQ
	conn, ch := rabbitmq.InitilizeBroker()
	defer conn.Close()
	defer ch.Close()

	// JWT_SECRET must be set for Auth signing
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		stdlog.Println("JWT Secret not set")
		return errors.New("JWT Secret not set")
	}
	opts := config.ConfigOpts{}
	opts.ConnectDB()
	// Setup app routes
	r := mux.NewRouter()
	routes.RegisterAuthRoutes(r, opts)
	routes.RegisterChatRoutes(r, opts)
	routes.RegisterWebsocketRoute(r)

	// Start api server
	port := os.Getenv("APP_PORT")
	fmt.Println("App started")

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	return err
}
