package repository

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/anton1ks96/college-app-core/internal/domain"
)

type PortalRepository struct {
	client                 *http.Client
	baseURL                string
	attendanceURL          string
	performanceSubjectsURL string
	performanceScoreURL    string
}

func NewPortalRepository(baseURL, attendanceURL, performanceSubjectsURL, performanceScoreURL string) *PortalRepository {
	return &PortalRepository{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		baseURL:                baseURL,
		attendanceURL:          attendanceURL,
		performanceSubjectsURL: performanceSubjectsURL,
		performanceScoreURL:    performanceScoreURL,
	}
}

func (r *PortalRepository) FetchSchedule(req domain.ScheduleRequest) ([]domain.ScheduleEvent, error) {
	body, _ := json.Marshal(req)

	resp, err := r.client.Post(fmt.Sprintf("%s/Services/schedule25.php", r.baseURL), "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var events []domain.ScheduleEvent
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, err
	}

	for i := 0; i < len(events); i++ {
		if strings.Contains(events[i].Start, " ") {
			parts := strings.Split(events[i].Start, " ")
			events[i].Start = parts[len(parts)-1]
		}
		if strings.Contains(events[i].End, " ") {
			parts := strings.Split(events[i].End, " ")
			events[i].End = parts[len(parts)-1]
		}
	}

	return events, nil
}

func (r *PortalRepository) FetchClassDetails(clid string) (map[string]any, error) {
	body, _ := json.Marshal(map[string]string{"clid": clid})

	resp, err := r.client.Post(fmt.Sprintf("%s/Services/classdetails25.php", r.baseURL), "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var details map[string]any
	if err := json.Unmarshal(data, &details); err != nil {
		return nil, err
	}
	return details, nil
}

func (r *PortalRepository) FetchAttendance(login string, req domain.AttendanceRequest) ([]domain.AttendanceRecord, error) {

	body, _ := json.Marshal(req)

	httpReq, err := http.NewRequest("POST", r.attendanceURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	cookieValue := fmt.Sprintf("STDNT-login-user=%s", login)
	httpReq.Header.Set("Cookie", fmt.Sprintf("session=%s", cookieValue))

	resp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var records []domain.AttendanceRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}

	return records, nil
}

func (r *PortalRepository) FetchPerformanceSubjects(login string) ([]domain.PerformanceSubject, error) {
	httpReq, err := http.NewRequest("GET", r.performanceSubjectsURL, nil)
	if err != nil {
		return nil, err
	}

	cookieValue := fmt.Sprintf("STDNT-login-user=%s", login)
	httpReq.Header.Set("Cookie", fmt.Sprintf("session=%s", cookieValue))

	resp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var subjects []domain.PerformanceSubject
	if err := json.Unmarshal(data, &subjects); err != nil {
		return nil, err
	}

	return subjects, nil
}

func (r *PortalRepository) FetchPerformanceScore(login string, req domain.PerformanceScoreRequest) (map[string]map[string][]domain.PerformanceScore, error) {
	body, _ := json.Marshal(req)

	httpReq, err := http.NewRequest("POST", r.performanceScoreURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	cookieValue := fmt.Sprintf("STDNT-login-user=%s", login)
	httpReq.Header.Set("Cookie", fmt.Sprintf("session=%s", cookieValue))

	resp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var scores map[string]map[string][]domain.PerformanceScore
	if err := json.Unmarshal(data, &scores); err != nil {
		return nil, err
	}

	return scores, nil
}
