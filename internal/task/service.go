package task

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)
var (
	ReportTypes = []string{"summary"}
	ErrInvalidPeriodDays = errors.New("invalid period days")
	ErrInvalidReportType = errors.New("invalid report type")
	ErrInvalidID 		 = errors.New("invalid ID format(need UUID)")
)

type Repository interface{
	Post(context.Context, PostInput) (Task, error)
	Get(context.Context, string) (Task,error)
}

type Service struct{
	repository Repository
}

func NewService(repo Repository) *Service{
	return &Service{
		repository: repo,
	}
}

func (s *Service) Post(ctx context.Context, input PostInput) (Task,error){
	if input.Parameters.PeriodDays < 1 || input.Parameters.PeriodDays > 365{
		return Task{}, ErrInvalidPeriodDays
	}

	if !slices.Contains(ReportTypes, input.ReportType){
		return Task{}, ErrInvalidReportType
	}

	item, err := s.repository.Post(ctx, input)
	if err!=nil{
		return Task{}, fmt.Errorf("post task error: %w", err)
	}
	return item, nil
}

func (s *Service) Get(ctx context.Context, ID string) (Task, error){
	if _, err := uuid.Parse(ID); err!=nil{
		return Task{}, ErrInvalidID
	}

	item, err := s.repository.Get(ctx, ID)
	if err!=nil{
		return Task{}, fmt.Errorf("get task error: %w", err)
	}
	return item, nil
}

func (s *Service) HealthStatus(ctx context.Context) (HealthStatus, error){
	responseTime := time.Now()

	return HealthStatus{
		Status: "ok",
		Response_Time_Ms: int32(time.Since(responseTime)),
	}, nil
}