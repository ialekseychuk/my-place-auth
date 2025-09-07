package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/ialekseychuk/my-place-identity/internal/config"
	"github.com/ialekseychuk/my-place-identity/internal/handler"
	"github.com/ialekseychuk/my-place-identity/internal/infrastructure"
	"github.com/ialekseychuk/my-place-identity/internal/interceptor"
	"github.com/ialekseychuk/my-place-identity/internal/repository"
	"github.com/ialekseychuk/my-place-identity/internal/usecase"

	identityv1 "github.com/ialekseychuk/my-place-proto/gen/go/identity/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pool, err := pgxpool.New(ctx, config.POSTGRES_DSN)
	if err != nil {
		logrus.Fatalf("unable to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logrus.Fatalf("unable to ping database: %v", err)
	}
	logrus.Println("Connected to the database")

	// repositories
	userRepo := repository.NewUserRepository(pool)
	tokenRepo := repository.NewTokenRepository(pool)

	// jwt
	jwtManager := infrastructure.NewJWTManager(config.JWT_SECRET)

	// usecases
	loginUC := usecase.NewLogin(userRepo, tokenRepo, jwtManager, config.AccessTTL, config.RefreshTTL)
	registerUC := usecase.NewRegister(userRepo, tokenRepo, jwtManager, config.AccessTTL, config.RefreshTTL)
	refreshUC := usecase.NewRefresh(userRepo, tokenRepo, jwtManager, config.AccessTTL, config.RefreshTTL)
	validateUC := usecase.NewValidateToken(userRepo, jwtManager)
	logoutUC := usecase.NewLogout(tokenRepo)
	getMeUC := usecase.NewGetMe(userRepo)

	// services
	identityHandler := handler.NewIdentityHandler(
		loginUC,
		registerUC,
		refreshUC,
		validateUC,
		logoutUC,
		getMeUC,
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.Auth(config.JWT_SECRET),
		),
	)
	
	identityv1.RegisterIdentityServer(grpcServer, identityHandler)

	go func() {
		logrus.Printf("Starting gRPC server on :%d", config.Port)
		if err := grpcServer.Serve(lis); err != nil {
			logrus.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	logrus.Println("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcServer.GracefulStop()
	logrus.Println("gRPC server stopped")

	select {
	case <-ctxShutdown.Done():
		logrus.Println("gRPC server stopped")
	}
}
