package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"taskforge/internal/task"
)

type Repository struct{
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository{
	return &Repository{
		pool: pool,
	}
}

func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error){
	pool, err := pgxpool.New(ctx,databaseURL)
	if err != nil {
		return nil, fmt.Errorf("pool create error: %w", err)
	}

	if err := pool.Ping(ctx); err!=nil{
		pool.Close()
		return nil, fmt.Errorf("pool ping error: %w", err)
	}
	return pool, nil
}

func (r *Repository) Post(ctx context.Context, task.PostInput) (task.Task, error){
	query := `
		INSERT INTO public.tasks (report_type, parameters, status)
		VALUES ($1, $2, $3)
		RETURNING id, report_type, parameters, status, result, error, attempts, created_at, started_at, finished_at
	`
	posted, err := scanTask(r.pool.QueryRow(ctx, query,
		))
}

func scanTask(row pgx.Row) (task.Task, error){
	var t task.Task
	err := row.Scan(
		&t.ID,
		&t.Report_type,
		&t.Parameters,
		&t.Status,
		&t.Result,
		&t.Error,
		&t.Attempts,
		&t.Created_at,
		&t.Started_at,
		&t.Finished_at,
	)
	if err!=nil{
		return task.Task{},err
	}
	return t,nil
}