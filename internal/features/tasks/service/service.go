package service

import (
	"chek-list/internal/features/tasks"
	"chek-list/internal/features/tasks/repository"
	"context"
)

type TaskService struct {
	repo repository.PostgresRepository
}

func NewTaskService(repo repository.PostgresRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) Create(ctx context.Context, title, description string) (*tasks.Task, error) {
	newTask := tasks.NewTask(title, description)
	err := s.repo.Create(ctx, newTask)
	return newTask, err
}

func (s *TaskService) GetAll(ctx context.Context) ([]*tasks.Task, error) {
	tasksList, err := s.repo.GetList(ctx)
	if err != nil {
		return nil, err
	}
	return tasksList, nil
}

func (s *TaskService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	return err
}
func (s *TaskService) Done(ctx context.Context, id int) (*tasks.Task, error) {
	taskNew, err := s.repo.Done(ctx, id)
	if err != nil {
		return nil, err
	}
	return taskNew, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, id int) (*tasks.Task, error) {
	taskNew, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return taskNew, nil
}
