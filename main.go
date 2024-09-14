package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
	"github.com/knockbox/authentication/internal/handlers"
	"github.com/knockbox/authentication/pkg/middleware"
	"github.com/knockbox/authentication/pkg/utils"
	"os"
)

var bindAddress string
var useDotEnv bool

func init() {
	const (
		usageBindAddress = "the address to bind to, e.g. :9090"
		usageUseDotEnv   = "read variables from a .env file in running directory"
	)

	flag.StringVar(&bindAddress, "bindAddress", ":9090", usageBindAddress)
	flag.StringVar(&bindAddress, "addr", ":9090", usageBindAddress)

	flag.BoolVar(&useDotEnv, "dotenv", false, usageUseDotEnv)
	flag.BoolVar(&useDotEnv, "denv", false, usageUseDotEnv)
}

func main() {
	flag.Parse()

	l := hclog.Default()

	if useDotEnv {
		l.Info("using .env")

		if err := godotenv.Load(); err != nil {
			l.Error(".env", "error", err)
			os.Exit(1)
		}
	}

	sm := mux.NewRouter()
	sm.Use(middleware.UseLogging(l).Middleware)

	// /api grouping
	apiRouter := sm.PathPrefix("/api").Subrouter()
	//apiRouter.Use(middleware.UseCaching(l).Middleware)

	// Routes
	handlers.NewHealthcheck().Route(apiRouter)
	handlers.NewUser(l).Route(apiRouter)

	utils.StartServerWithGracefulShutdown(sm, bindAddress, l)
}
