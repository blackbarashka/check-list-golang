package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func main() {
	logFile, err := os.OpenFile("logs/events.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть файл логов: %v", err)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {

		}
	}(logFile)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"kafka:9092"},
		Topic:     "task-events",
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})
	defer func(r *kafka.Reader) {
		err := r.Close()
		if err != nil {

		}
	}(r)

	fmt.Println("Воркер успешно запущен и слушает топик 'task-events'...")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Ошибка при чтении сообщения: %v", err)
			break
		}

		logLine := fmt.Sprintf("[KAFKA LOG] Событие: %s\n", string(m.Value))

		fmt.Print(logLine)

		if _, err := logFile.WriteString(logLine); err != nil {
			log.Printf("Ошибка записи в файл: %v", err)
		}
	}
}
