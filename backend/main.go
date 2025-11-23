package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"backend/api"
	"backend/server"
	"backend/service"
	"backend/service/auth"

	jwtPackage "backend/pkg/jwt"
	storageAuth "backend/storage/auth"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	logger *logrus.Logger
	config *viper.Viper

	port = "8090"
)

func init() {
	config = viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	logger = logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
}

func main() {
	privateKeyBase64 := config.GetString("jwt.privatePEM")
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		log.Fatal("fail to decode private key", err.Error())
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyBytes))
	if err != nil {
		log.Fatal("fail to parse jwt private key", err.Error())
	}

	publicKeyBase64 := config.GetString("jwt.publicPEM")
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		log.Fatal("fail to decode public key", err.Error())
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyBytes))
	if err != nil {
		log.Fatal("fail to parse jwt public key", err.Error())
	}

	issuer := config.GetString("jwt.issuer")
	if issuer == "" {
		issuer = "boilerplate-v2"
	}

	tokenDurationMinutes := config.GetInt("jwt.tokenDurationMinutes")
	if tokenDurationMinutes == 0 {
		tokenDurationMinutes = 15
	}

	j := jwtPackage.New(&jwtPackage.NewJWTOptions{
		PrivateKey:    privateKey,
		PublicKey:     publicKey,
		Issuer:        issuer,
		TokenDuration: time.Duration(tokenDurationMinutes) * time.Minute,
	})

	storage, err := storageAuth.NewStorage(logger, config)
	if err != nil {
		logger.Fatal("error initializing postgres storage", err.Error())
	}

	authService := auth.New(logger, storage)
	service := service.New(service.ServiceInitParams{
		Auth: authService,
	})

	server := server.New(server.ServerInitParams{
		Service: service,
		JWT:     j,
	})

	r := chi.NewMux()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(server.WithAuth)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	h := api.HandlerFromMux(server, r)

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	s := &http.Server{
		Handler: h,
		Addr:    addr,
	}

	log.Fatal(s.ListenAndServe())
}
