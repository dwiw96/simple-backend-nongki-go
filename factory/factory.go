package factory

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"

	authCache "simple-backend-nongki-go/features/auth/cache"
	authDelivery "simple-backend-nongki-go/features/auth/delivery"
	authRepository "simple-backend-nongki-go/features/auth/repository"
	authService "simple-backend-nongki-go/features/auth/service"
)

func InitFactory(router *httprouter.Router, pool *pgxpool.Pool, rdClient *redis.Client, ctx context.Context) {
	authRepoInterface := authRepository.NewAuthRepository(pool, ctx)
	authCacheInterface := authCache.NewAuthCache(rdClient, ctx)
	authServiceInterface := authService.NewAuthService(authRepoInterface, authCacheInterface)
	authDelivery.NewAuthDelivery(router, authServiceInterface, pool, rdClient, ctx)
}
