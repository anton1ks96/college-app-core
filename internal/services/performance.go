package services

import (
	"fmt"

	"github.com/anton1ks96/college-app-core/internal/domain"
	"github.com/anton1ks96/college-app-core/internal/repository"
)

type PerformanceService struct {
	portal *repository.PortalRepository
}

func NewPerformanceService(portal *repository.PortalRepository) *PerformanceService {
	return &PerformanceService{
		portal: portal,
	}
}

func (s *PerformanceService) GetSubjects(login string) ([]domain.PerformanceSubject, error) {
	subjects, err := s.portal.FetchPerformanceSubjects(login)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch performance subjects: %w", err)
	}

	return subjects, nil
}

func (s *PerformanceService) GetScore(login, suID, start, end string) (map[string]map[string][]domain.PerformanceScore, error) {
	req := domain.PerformanceScoreRequest{
		SuID:      suID,
		Datastart: start,
		Dataend:   end,
	}

	scores, err := s.portal.FetchPerformanceScore(login, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch performance score: %w", err)
	}

	return scores, nil
}
