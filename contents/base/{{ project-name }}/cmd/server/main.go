package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"{{ module_path }}/internal/config"
{% if has_persistence and persistence == "PostgreSQL" %}	// "github.com/jackc/pgx/v5/pgxpool"
{% endif %}{% if has_persistence and persistence == "MySQL" %}	// "database/sql"
	// _ "github.com/go-sql-driver/mysql"
{% endif %}{% if has_cache %}	// "github.com/redis/go-redis/v9"
{% endif %}{% if messaging == "Kafka" %}	// "github.com/IBM/sarama"
{% endif %}{% if messaging == "Pulsar" %}	// "github.com/apache/pulsar-client-go/pulsar"
{% endif %}{% if has_s3 or has_azure_blob %}	"{{ module_path }}/internal/storage"
{% endif %})

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "error", err)
		os.Exit(1)
	}

	// Structured logging: JSON in production, text in development
	var logHandler slog.Handler
	if cfg.LoggingJSON {
		logHandler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, nil)
	}
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

{% if has_persistence and persistence == "PostgreSQL" %}	// pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	// if err != nil {
	// 	slog.Error("db init failed", "error", err); os.Exit(1)
	// }
	// defer pool.Close()

{% endif %}{% if has_persistence and persistence == "MySQL" %}	// db, err := sql.Open("mysql", cfg.DatabaseURL)
	// if err != nil {
	// 	slog.Error("db init failed", "error", err); os.Exit(1)
	// }
	// defer db.Close()

{% endif %}{% if has_cache %}	// rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisURL})
	// defer rdb.Close()

{% endif %}{% if messaging == "Kafka" %}	// if err := messaging.Init(cfg.KafkaBrokers, cfg.KafkaUsername, cfg.KafkaPassword, cfg.KafkaSASLMechanism); err != nil {
	// 	slog.Error("kafka init failed", "error", err); os.Exit(1)
	// }
	// defer messaging.Close()

{% endif %}{% if messaging == "Pulsar" %}	// if err := messaging.Init(cfg.PulsarBrokerURL, cfg.PulsarTopic, cfg.PulsarJWTToken, cfg.PulsarSubscriptionName); err != nil {
	// 	slog.Error("pulsar init failed", "error", err); os.Exit(1)
	// }
	// defer messaging.Close()

{% endif %}{% if has_s3 %}	if err := storage.InitS3(cfg.S3); err != nil {
		slog.Error("storage s3 init failed", "error", err)
		os.Exit(1)
	}
{% endif %}{% if has_azure_blob %}	if err := storage.InitAzureBlob(cfg.Azure); err != nil {
		slog.Error("storage azure-blob init failed", "error", err)
		os.Exit(1)
	}
{% endif %}
	grpcSrv := grpc.NewServer()
	// TODO: register {{ PrefixName }}{{ SuffixName }}Service after running 'make proto':
	// pb.Register{{ PrefixName }}{{ SuffixName }}Server(grpcSrv, &service.{{ PrefixName }}{{ SuffixName }}Service{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		slog.Error("listen failed", "error", err)
		os.Exit(1)
	}

	mgmtMux := http.NewServeMux()
	mgmtMux.HandleFunc("/health/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok"}`)
	})
	mgmtMux.HandleFunc("/health/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok"}`)
	})
	mgmtMux.Handle("/metrics", promhttp.Handler())
	mgmt := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.ManagementPort),
		Handler: mgmtMux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("management server starting", "port", cfg.ManagementPort)
		if err := mgmt.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("management server error", "error", err)
		}
	}()

	go func() {
		slog.Info("grpc server starting", "port", cfg.Port)
		if err := grpcSrv.Serve(lis); err != nil {
			slog.Error("grpc server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down...")

	grpcSrv.GracefulStop()

	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = mgmt.Shutdown(shutCtx)
}
