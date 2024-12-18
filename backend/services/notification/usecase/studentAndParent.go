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

func (spu *studentParentUseCase) CreateStudentAndParentUC(ctx context.Context, req *domain.StudentAndParent) *[]string {
	ctx, cancel := context.WithTimeout(ctx, spu.TimeOut)
	defer cancel()

	err := spu.repo.CreateStudentAndParent(ctx, req)
	if err != nil {
		return err
	}
	return nil
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

func (spu *studentParentUseCase) UpdateStudentAndParent(ctx context.Context, id int, payload *domain.StudentAndParent) *[]string {

	// ctx, cancel := context.WithTimeout(ctx, spu.TimeOut)
	// defer cancel()

	err := spu.repo.UpdateStudentAndParent(ctx, id, payload)
	if err != nil {
		return err
	}
	return nil
}

func (spu *studentParentUseCase) GetStudentDetailsByID(ctx context.Context, id int) (*domain.StudentAndParent, error) {

	// ctx, cancel := context.WithTimeout(ctx, spu.TimeOut)
	// defer cancel()

	v, err := spu.repo.GetStudentDetailsByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (spu *studentParentUseCase) DeleteStudentAndParent(ctx context.Context, id int) error {

	// ctx, cancel := context.WithTimeout(ctx, spu.TimeOut)
	// defer cancel()

	err := spu.repo.DeleteStudentAndParent(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (spu *studentParentUseCase) DataChangeRequest(ctx context.Context, datas domain.DataChangeRequest) error {

	// ctx, cancel := context.WithTimeout(ctx, spu.TimeOut)
	// defer cancel()

	err := spu.repo.DataChangeRequest(ctx, datas)
	if err != nil {
		return err
	}
	return nil
}

func (spu *studentParentUseCase) GetAllDataChangeRequest(ctx context.Context) (*[]domain.DataChangeRequest, error) {

	// ctx, cancel := context.WithTimeout(ctx, spu.TimeOut)
	// defer cancel()

	v, err := spu.repo.GetAllDataChangeRequest(ctx)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (spu *studentParentUseCase) GetAllDataChangeRequestByID(ctx context.Context, dcrID int) (*domain.DataChangeRequest, error) {

	// ctx, cancel := context.WithTimeout(ctx, spu.TimeOut)
	// defer cancel()

	v, err := spu.repo.GetAllDataChangeRequestByID(ctx, dcrID)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (spu *studentParentUseCase) SPMassDelete(ctx context.Context, studentIDS *[]int) error {
	err := spu.repo.SPMassDelete(ctx, studentIDS)
	if err != nil {
		return err
	}
	return nil
}

func (spu *studentParentUseCase) ReviewDCR(ctx context.Context, dcrID int) error {
	err := spu.repo.ReviewDCR(ctx, dcrID)
	if err != nil {
		return err
	}
	return nil
}

func (spu *studentParentUseCase) ApproveDCR(ctx context.Context, req map[string]interface{}) (*string, error) {
	b, err := spu.repo.ApproveDCR(ctx, req)
	if err != nil {
		return b, err
	}
	return b, nil
}

// func (spu *studentParentUseCase) GetClassIDByName(className string) (*int, error) {

// 	// ctx, cancel := context.WithTimeout(ctx, spu.TimeOut)
// 	// defer cancel()

// 	v, err := spu.repo.GetClassIDByName(className)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return v, nil
// }
