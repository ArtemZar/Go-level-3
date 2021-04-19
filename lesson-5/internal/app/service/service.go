package service

import (
	"github.com/ArtemZar/Go-level-3/lesson-5/internal/pkg/model"
)

type (
	repo interface {
		Get(key string) (string, error)
		Put(putReq *model.PutValue) error
	}

	Service struct {
		repo repo
	}
)

func New(repo repo) *Service {
	return &Service{repo: repo}
}
