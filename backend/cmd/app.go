package cmd

import (
	"context"
	"finder/internal/cache"
	"finder/internal/config"
	"finder/internal/handler"
	"finder/internal/provider"
	"finder/internal/repository"
	"finder/internal/service"
	"finder/pkg/db"
	"finder/pkg/http/client"
	echohttp "finder/pkg/http/echo"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	_ "finder/docs"

	"github.com/labstack/echo/v5"
)

func Execute() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		return
	}

	pool, err := db.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return
	}
	defer pool.Close()

	redisClient, err := db.NewRedisClient(cfg.RedisURL)
	if err != nil {
		slog.Error("failed to connect to redis", "error", err)
		return
	}
	defer redisClient.Close()

	redisCache := cache.NewRedisCache(redisClient)

	repo := repository.NewPostgresRepository(pool)
	cachedRepo := repository.NewCachedRepository(repo, redisCache, cfg.CacheTTL)

	p1Client := client.New(client.Config{BaseURL: cfg.Provider1.BaseURL})
	p2Client := client.New(client.Config{BaseURL: cfg.Provider2.BaseURL})

	p1json := provider.NewJSONProvider(cfg.Provider1.Name, p1Client)
	p1cb := provider.NewCircuitBreakerProvider(cfg.Provider1.Name, p1json, cfg.CBTimeout, cfg.CBThreshold)
	p1rl := provider.NewRateLimitedProvider(cfg.Provider1.Name, p1cb, redisCache, cfg.RateLimitPerSec)
	p1 := provider.NewCachedProvider(cfg.Provider1.Name, p1rl, redisCache, cfg.ProviderCacheTTL)

	p2xml := provider.NewXMLProvider(cfg.Provider2.Name, p2Client)
	p2cb := provider.NewCircuitBreakerProvider(cfg.Provider2.Name, p2xml, cfg.CBTimeout, cfg.CBThreshold)
	p2rl := provider.NewRateLimitedProvider(cfg.Provider2.Name, p2cb, redisCache, cfg.RateLimitPerSec)
	p2 := provider.NewCachedProvider(cfg.Provider2.Name, p2rl, redisCache, cfg.ProviderCacheTTL)

	aggregator := provider.NewAggregator(p1, p2)

	svc := service.New(cachedRepo, aggregator)
	h := handler.New(svc)

	e := echohttp.New()
	h.RegisterRoutes(e)

	sc := echo.StartConfig{
		Address:         ":" + cfg.Port,
		GracefulTimeout: 10 * time.Second,
	}
	if err := sc.Start(ctx, e); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
