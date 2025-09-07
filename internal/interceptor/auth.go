package interceptor

import (
	"context"
	"strings"

	"github.com/ialekseychuk/my-place-identity/internal/infrastructure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ContextKey string

const UserIDKey ContextKey = "user_id"

func Auth(jwtSecret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		public := map[string]struct{}{
			"/identity.Identity/Register": {},
			"/identity.Identity/Login":     {},
		}

		if _, ok := public[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok || len(md.Get("authorization")) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}
		tokenStr := strings.TrimPrefix(md.Get("authorization")[0], "Bearer ")

		jwtMgr := infrastructure.NewJWTManager(jwtSecret)
		claims, err := jwtMgr.ValidateAccessToken(tokenStr)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		return handler(ctx, req)
	}
}
