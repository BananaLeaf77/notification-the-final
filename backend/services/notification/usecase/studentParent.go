package usecase

import (
	"context"
	"notification/domain"
	"time"
)

type studentParentUseCase struct {
	repo    domain.StudentParentRepo
	TimeOut time.Duration
}

func NewStudentParentUseCase(repo domain.StudentParentRepo, to time.Duration) domain.StudentParentUseCase {
	return &studentParentUseCase{
		repo:    repo,
		TimeOut: to,
	}
}

func (spu *studentParentUseCase) CreateStudentAndParentUC(ctx context.Context, req *domain.StudentAndParent) error {
	ctx, cancel := context.WithTimeout(ctx, spu.TimeOut)
	defer cancel()

	err := spu.repo.CreateStudentAndParent(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (spu *studentParentUseCase) GetStudentAndParent(ctx context.Context, studentID string) (*domain.StudentAndParent, error) {
	ctx, cancel := context.WithTimeout(ctx, spu.TimeOut)
	defer cancel()

	data, err := spu.repo.GetStudentAndParent(ctx, studentID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (spu *studentParentUseCase) ImportCSV(ctx context.Context, payload *[]domain.StudentAndParent) (*[]string, error) {
	// ctx, cancel := context.WithTimeout(ctx, spu.TimeOut)
	// defer cancel()

	data, err := spu.repo.ImportCSV(ctx, payload)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (spu *studentParentUseCase) UpdateStudentandParent(ctx context.Context, id int64, student domain.Student, parent domain.Parent) error {

	// ctx, cancel := context.WithTimeout(ctx, spu.TimeOut)
	// defer cancel()

	err := spu.repo.UpdateStudentandParent(ctx, id, &student, &parent)
	if err != nil {
		return err
	}
	return nil
}
