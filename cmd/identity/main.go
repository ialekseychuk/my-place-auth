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
	"github.com/ialekseychuk/my-place-identity/internal/repository"
	grpcSvc "github.com/ialekseychuk/my-place-identity/internal/transport/grpc"
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

	userRepo := repository.NewRepository(pool)
	identityHandler := grpcSvc.NewIdentityHandler(userRepo, config.JWT_SECRET)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Port))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	identityv1.RegisterAuthServer(grpcServer, identityHandler)

	go func() {
		logrus.Println("Starting gRPC server on :50051")
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
