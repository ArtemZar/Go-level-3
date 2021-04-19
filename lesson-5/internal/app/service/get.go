package service

import (
	"log"

	web_broker "github.com/ArtemZar/Go-level-3/lesson-5/pkg/web-broker"
)

func (s *Service) Get(req *web_broker.GetValueReq) (*web_broker.GetValueResp, error) {
	value, err := s.repo.Get(req.Key)
	if err != nil {
		log.Printf("service/Get: get from repo err: %v", err)
		return nil, err
	}

	return &web_broker.GetValueResp{Value: value}, nil
}
