package repository

import "github.com/ArtemZar/Go-level-3/lesson-5/internal/pkg/model"

type MemRepo struct{}

func NewMemRepo() *MemRepo {
	return &MemRepo{}
}

func (mr *MemRepo) Get(getReq string) (string, error) {
	// TODO: impl
	return "", nil
}

func (mr *MemRepo) Put(putReq *model.PutValue) error {
	// TODO: impl
	return nil
}
