package main

import (
	"chek-list/api"
	"chek-list/core/db"
	"chek-list/internal/features/tasks/delivery"
	"chek-list/internal/features/tasks/repository"
	"chek-list/internal/features/tasks/service"
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"

	// Используем алиас, чтобы имя пакета не конфликтовало с библиотекой grpc
	grpcDelivery "chek-list/internal/features/tasks/delivery/grpc"
)

func main() {
	ctx_signal, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	connString := os.Getenv("DB_CONN_STRING")
	if connString == "" {
		log.Fatal("Переменная DB_CONN_STRING не задана")
	}
	ctx := context.Background()
	pool, err := db.NewPostgresPool(ctx, connString)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	//Kafka.
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "task-events",
		Balancer: &kafka.LeastBytes{},
	}
	defer func(kafkaWriter *kafka.Writer) {
		err := kafkaWriter.Close()
		if err != nil {

		}
	}(kafkaWriter)

	repo := repository.NewPostgresRepository(pool)
	s := service.NewTaskService(repo)

	handler := delivery.NewHandler(s)

	grpcHandler := grpcDelivery.NewHandler(s, kafkaWriter)
	grpcServer := grpc.NewServer()
	api.RegisterTaskServiceServer(grpcServer, grpcHandler)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		log.Println("grpc server listening on :50051")
		err := grpcServer.Serve(lis)
		if err != nil {
			return
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /create", handler.CreateHTTP)
	mux.HandleFunc("GET /list", handler.GetAllHTTP)
	mux.HandleFunc("DELETE /delete/{id}", handler.DeleteHTTP)
	mux.HandleFunc("PUT /done/{id}", handler.DoneHTTP)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx_signal.Done()
	log.Println("Получен сигнал. Останавливаем сервер...")

	grpcServer.GracefulStop()
	log.Println("gRPC сервер остановлен")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	log.Println("Сервер успешно остановлен")
}
