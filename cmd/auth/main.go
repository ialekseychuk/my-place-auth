package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/ialekseychuk/my-place-auth/internal/auth"
	"github.com/ialekseychuk/my-place-auth/internal/repository"
	authv1 "github.com/ialekseychuk/my-place-proto/gen/go/auth/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {

	ctx := context.Background()

	dsn := os.Getenv("POSTGRES_DSN")
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logrus.Fatalf("unable to connect to database: %v", err)

	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logrus.Fatalf("unable to ping database: %v", err)
	}
	logrus.Println("Connected to the database")

	userRepo := repository.NewRepository(pool)
	secret := os.Getenv("JWT_SECRET")
	authService := auth.NewAuthService(userRepo, secret)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServer(grpcServer, authService)

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
