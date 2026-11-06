package repository

import (
	"chek-list/internal/features/tasks"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository interface {
	Create(ctx context.Context, task *tasks.Task) error
	GetList(ctx context.Context) (task []*tasks.Task, err error)
	Delete(ctx context.Context, id int) error
	Done(ctx context.Context, id int) (*tasks.Task, error)
	GetByID(ctx context.Context, id int) (*tasks.Task, error)
}

type PostgresRepositoryImpl struct {
	conn *pgxpool.Pool
}

func NewPostgresRepository(conn *pgxpool.Pool) *PostgresRepositoryImpl {
	return &PostgresRepositoryImpl{
		conn: conn,
	}
}

func (r *PostgresRepositoryImpl) Create(ctx context.Context, task *tasks.Task) error {
	query := `
    INSERT INTO tasksTable(title, description, status, addTime)
    VALUES ($1, $2, $3, $4)
    RETURNING id
`
	err := r.conn.QueryRow(ctx, query, task.Title, task.Description, task.Status, task.AddTime).Scan(&task.ID)
	return err
}

func (r *PostgresRepositoryImpl) GetList(ctx context.Context) (task []*tasks.Task, err error) {
	query := `
    SELECT id, title, description, status, addTime, completedTime
    FROM tasksTable`

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	taskList := make([]*tasks.Task, 0)
	for rows.Next() {
		var task tasks.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.AddTime, &task.CompletedTime)
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, &task)
	}
	return taskList, nil
}

func (r *PostgresRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tasksTable WHERE id = $1`
	_, err := r.conn.Exec(ctx, query, id)
	return err
}

func (r *PostgresRepositoryImpl) Done(ctx context.Context, id int) (*tasks.Task, error) {
	query := `
	UPDATE tasksTable
    SET status = $1, completedTime = $2
	WHERE id = $3
	RETURNING id, title, description, status, addTime, completedTime`

	var task tasks.Task
	err := r.conn.QueryRow(ctx, query, true, time.Now(), id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.AddTime,
		&task.CompletedTime,
	)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *PostgresRepositoryImpl) GetByID(ctx context.Context, id int) (*tasks.Task, error) {
	query := `
    SELECT id, title, description, status, addTime, completedTime
    FROM tasksTable
    WHERE id = $1`
	var task tasks.Task
	err := r.conn.QueryRow(ctx, query, id).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.AddTime, &task.CompletedTime)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
