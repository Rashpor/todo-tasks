package storage

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"rashpor.com/todolist/internal/storage/models"
)

type TaskRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewTaskRepository(db *pgxpool.Pool, logger *slog.Logger) *TaskRepository {
	return &TaskRepository{db: db, logger: logger}
}

func (r *TaskRepository) CreateNewTask(ctx context.Context, task *models.Task) error {
	query := `
	 INSERT INTO tasks (title, description, status, created_at, updated_at) 
	 VALUES ($1, $2, $3, $4, $5) 
	 RETURNING id, created_at, updated_at`
	r.logger.Debug("Выполняем SQL-запрос", slog.String("query", query))

	now := time.Now()
	err := r.db.QueryRow(ctx, query, task.Title, task.Description, task.Status, now, now).
		Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) GetAllTasks(ctx context.Context, page, limit int) ([]*models.Task, error) {
	offset := (page - 1) * limit
	query := "SELECT id, title, description, status, created_at, updated_at FROM tasks ORDER BY created_at DESC LIMIT $1 OFFSET $2"
	r.logger.Debug("Выполняем SQL-запрос", slog.String("query", query))

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("Ошибка выполнения SQL-запроса", slog.Any("error", err))
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			r.logger.Error("Ошибка при сканировании строки", slog.Any("error", err))
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) UpdateTaskOnID(ctx context.Context, id int, task *models.Task) error {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1)`
	r.logger.Debug("Выполняем SQL-запрос", slog.String("query", checkQuery))
	err := r.db.QueryRow(ctx, checkQuery, id).Scan(&exists)
	if err != nil || !exists {
		return models.ErrNoRecord
	}

	query := `UPDATE tasks SET title = $1, description = $2, status = $3, updated_at = now() WHERE id = $4`
	r.logger.Debug("Выполняем SQL-запрос", slog.String("query", query))
	_, err = r.db.Exec(ctx, query, task.Title, task.Description, task.Status, id)
	return err
}

func (r *TaskRepository) DeleteTaskOnID(ctx context.Context, id int) error {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1)`
	r.logger.Debug("Выполняем SQL-запрос", slog.String("query", checkQuery))
	err := r.db.QueryRow(ctx, checkQuery, id).Scan(&exists)
	if err != nil || !exists {
		return models.ErrNoRecord
	}
	query := `DELETE FROM tasks WHERE id = $1`
	r.logger.Debug("Выполняем SQL-запрос", slog.String("query", query))
	_, err = r.db.Exec(ctx, query, id)
	return err
}
