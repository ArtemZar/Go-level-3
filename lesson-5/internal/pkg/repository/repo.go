package repository

import "github.com/ArtemZar/Go-level-3/lesson-5/internal/pkg/model"

type repo interface {
	Get(getReq string) (string, error)
	Put(putReq *model.PutValue) error
}