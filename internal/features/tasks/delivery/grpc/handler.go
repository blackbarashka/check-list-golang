package grpc

import (
	"chek-list/api"
	"chek-list/internal/features/tasks/service"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	service *service.TaskService
	api.UnimplementedTaskServiceServer
	kafkaWriter *kafka.Writer
}

func NewHandler(server *service.TaskService, kafka *kafka.Writer) *Handler {
	return &Handler{service: server, kafkaWriter: kafka}
}

func (h *Handler) CreateTask(ctx context.Context, in *api.CreateTaskRequest) (*api.CreateTaskResponse, error) {
	t, err := h.service.Create(ctx, in.GetTitle(), in.GetDescription())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event := &Event{
		Action: "create task",
		Time:   time.Now(),
	}
	events, _ := json.Marshal(event)

	go func() {
		msg := &kafka.Message{Value: events}
		if err := h.kafkaWriter.WriteMessages(context.Background(), *msg); err != nil {
			log.Println(err)
		}
	}()
	return &api.CreateTaskResponse{
		Task: &api.Task{
			Id:          int64(t.ID),
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status,
			AddTime:     timestamppb.New(t.AddTime),
		},
	}, nil
}

func (h *Handler) GetTasks(ctx context.Context, in *api.GetTaskRequest) (*api.GetTaskResponse, error) {
	tasks, err := h.service.GetAll(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	var pbTasks []*api.Task
	for _, task := range tasks {
		pdTask := &api.Task{
			Id:          int64(task.ID),
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			AddTime:     timestamppb.New(task.AddTime),
		}
		pbTasks = append(pbTasks, pdTask)
	}
	return &api.GetTaskResponse{
		TasksList: pbTasks,
	}, nil
}

func (h *Handler) MarkTaskDone(ctx context.Context, in *api.MarkTaskDoneRequest) (*api.MarkTaskDoneResponse, error) {
	t, err := h.service.Done(ctx, int(in.GetId()))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &api.MarkTaskDoneResponse{
		Task: &api.Task{
			Id:          int64(t.ID),
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status,
			AddTime:     timestamppb.New(t.AddTime),
		},
	}, nil
}

func (h *Handler) DeleteTask(ctx context.Context, in *api.DeleteTaskRequest) (*api.DeleteTaskResponse, error) {
	err := h.service.Delete(ctx, int(in.GetId()))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &api.DeleteTaskResponse{}, nil
}

func (h *Handler) GetTaskByID(ctx context.Context, in *api.GetTaskByIDRequest) (*api.GetTaskByIDResponse, error) {
	taskNew, err := h.service.GetTaskByID(ctx, int(in.Id))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	pbTask := &api.Task{
		Id:          int64(taskNew.ID),
		Title:       taskNew.Title,
		Description: taskNew.Description,
		Status:      taskNew.Status,
		AddTime:     timestamppb.New(taskNew.AddTime),
	}

	// 2. Безопасно проверяем указатель на nil
	if taskNew.CompletedTime != nil {
		pbTask.CompletedTime = timestamppb.New(*taskNew.CompletedTime)
	}

	return &api.GetTaskByIDResponse{
		Task: pbTask,
	}, nil
}
