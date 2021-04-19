package service

import (
	"log"

	"github.com/ArtemZar/Go-level-3/lesson-5/internal/pkg/model"
	web_broker "github.com/ArtemZar/Go-level-3/lesson-5/pkg/web-broker"
)

func (s *Service) Put(req *web_broker.PutValueReq) error {
	if err := s.repo.Put(&model.PutValue{
		Key:   req.Key,
		Value: req.Value,
	}); err != nil {
		log.Printf("service/Put: put repo err: %v", err)
		return err
	}

	return nil
}
