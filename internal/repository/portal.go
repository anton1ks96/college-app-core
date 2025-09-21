package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/anton1ks96/college-app-core/internal/domain"
)

type PortalRepository struct {
	client  *http.Client
	baseURL string
}

func NewPortalRepository(baseURL string) *PortalRepository {
	return &PortalRepository{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}

func (r *PortalRepository) FetchSchedule(req domain.ScheduleRequest) ([]domain.ScheduleEvent, error) {
	body, _ := json.Marshal(req)

	resp, err := r.client.Post(fmt.Sprintf("%s/schedule25.php", r.baseURL), "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var events []domain.ScheduleEvent
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *PortalRepository) FetchClassDetails(clid string) (map[string]any, error) {
	body, _ := json.Marshal(map[string]string{"clid": clid})

	resp, err := r.client.Post(fmt.Sprintf("%s/classdetails25.php", r.baseURL), "application/json", bytes.NewReader(body))
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
