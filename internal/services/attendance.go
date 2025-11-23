package services

import (
	"fmt"

	"github.com/anton1ks96/college-app-core/internal/domain"
	"github.com/anton1ks96/college-app-core/internal/repository"
)

type AttendanceService struct {
	portal *repository.PortalRepository
}

func NewAttendanceService(portal *repository.PortalRepository) *AttendanceService {
	return &AttendanceService{
		portal: portal,
	}
}

func (s *AttendanceService) GetAttendance(login, start, end string) ([]domain.AttendanceRecord, error) {
	req := domain.AttendanceRequest{
		DStart: start,
		DEnd:   end,
	}

	records, err := s.portal.FetchAttendance(login, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attendance: %w", err)
	}

	for i := range records {
		if len(records[i].SubGroup) == 1 {
			sg := records[i].SubGroup[0]
			records[i].Title = sg.STitle

			if records[i].Topic == "" {
				records[i].Topic = sg.STopic
			}
			if records[i].Room == "" {
				records[i].Room = sg.SCaID
			}
			records[i].SubGroup = nil
		}
	}

	return records, nil
}
