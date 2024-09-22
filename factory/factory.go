package factory

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"

	authDelivery "simple-backend-nongki-go/features/auth/delivery"
	authRepository "simple-backend-nongki-go/features/auth/repository"
	authService "simple-backend-nongki-go/features/auth/service"
)

func InitFactory(router *httprouter.Router, pool *pgxpool.Pool, ctx context.Context) {
	authRepoInterface := authRepository.NewAuthRepository(pool, ctx)
	authServiceInterface := authService.NewAuthService(authRepoInterface)
	authDelivery.NewAuthDelivery(router, authServiceInterface)
}
