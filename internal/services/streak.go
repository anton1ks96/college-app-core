package services

import (
	"fmt"
	"sort"
	"time"

	"github.com/anton1ks96/college-app-core/internal/domain"
)

func (s *AttendanceService) GetAttendanceStreak(login string) (*domain.StreakResponse, error) {
	startDate := getAcademicYearStart()
	endDate := getToday()

	req := domain.AttendanceRequest{
		DStart: startDate,
		DEnd:   endDate,
	}

	records, err := s.portal.FetchAttendance(login, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attendance for streak: %w", err)
	}

	return s.calculateStreak(records, startDate, endDate), nil
}

func (s *AttendanceService) calculateStreak(records []domain.AttendanceRecord, periodStart, periodEnd string) *domain.StreakResponse {
	if len(records) == 0 {
		return &domain.StreakResponse{
			PeriodStart: periodStart,
			PeriodEnd:   periodEnd,
		}
	}

	dayStatus := s.groupByDayAndDetermineStatus(records)
	dates := s.getSortedDatesDesc(dayStatus)

	currentStreak := s.calcCurrentStreak(dates, dayStatus)
	longestStreak := s.calcLongestStreak(dates, dayStatus)
	totalAttended := s.countAttended(dayStatus)
	lastAttended := s.findLastAttended(dates, dayStatus)

	totalDays := len(dayStatus)
	var rate float64
	if totalDays > 0 {
		rate = float64(totalAttended) / float64(totalDays)
	}

	return &domain.StreakResponse{
		CurrentStreak:     currentStreak,
		LongestStreak:     longestStreak,
		TotalDaysAttended: totalAttended,
		TotalSchoolDays:   totalDays,
		AttendanceRate:    rate,
		LastAttendedDate:  lastAttended,
		PeriodStart:       periodStart,
		PeriodEnd:         periodEnd,
	}
}

func (s *AttendanceService) groupByDayAndDetermineStatus(records []domain.AttendanceRecord) map[string]bool {
	dayStatus := make(map[string]bool)

	for _, r := range records {
		day := r.Day
		if dayStatus[day] {
			continue
		}
		if r.Status == 2 {
			dayStatus[day] = true
		} else if _, exists := dayStatus[day]; !exists {
			dayStatus[day] = false
		}
	}

	return dayStatus
}

func (s *AttendanceService) getSortedDatesDesc(dayStatus map[string]bool) []string {
	dates := make([]string, 0, len(dayStatus))
	for date := range dayStatus {
		dates = append(dates, date)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))
	return dates
}

func (s *AttendanceService) calcCurrentStreak(sortedDatesDesc []string, dayStatus map[string]bool) int {
	streak := 0
	for _, date := range sortedDatesDesc {
		if dayStatus[date] {
			streak++
		} else {
			break
		}
	}
	return streak
}

func (s *AttendanceService) calcLongestStreak(sortedDatesDesc []string, dayStatus map[string]bool) int {
	longest := 0
	current := 0

	for i := len(sortedDatesDesc) - 1; i >= 0; i-- {
		date := sortedDatesDesc[i]
		if dayStatus[date] {
			current++
			if current > longest {
				longest = current
			}
		} else {
			current = 0
		}
	}

	return longest
}

func (s *AttendanceService) countAttended(dayStatus map[string]bool) int {
	count := 0
	for _, attended := range dayStatus {
		if attended {
			count++
		}
	}
	return count
}

func (s *AttendanceService) findLastAttended(sortedDatesDesc []string, dayStatus map[string]bool) string {
	for _, date := range sortedDatesDesc {
		if dayStatus[date] {
			return date
		}
	}
	return ""
}

func getAcademicYearStart() string {
	now := time.Now()
	year := now.Year()

	if now.Month() < time.September {
		year--
	}

	return time.Date(year, time.September, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
}

func getToday() string {
	return time.Now().Format("2006-01-02")
}
