package auth

import (
	"github.com/google/uuid"
)

type Service struct {
	sessionIDGenerator func() string
}

func NewService(sessionIDGenerator func() string) *Service {
	if sessionIDGenerator == nil {
		sessionIDGenerator = func() string {
			return uuid.New().String()
		}
	}
	return &Service{
		sessionIDGenerator: sessionIDGenerator,
	}
}

func (service *Service) NewSessionID() string {
	return service.sessionIDGenerator()
}
