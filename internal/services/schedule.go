package services

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/anton1ks96/college-app-core/internal/domain"
	"github.com/anton1ks96/college-app-core/internal/repository"
)

type ScheduleService struct {
	portal *repository.PortalRepository
}

func NewScheduleService(portal *repository.PortalRepository) *ScheduleService {
	return &ScheduleService{
		portal: portal,
	}
}

var englishRe = regexp.MustCompile(`^(A0|A1|A2|B1)\.\d{2}$`)

func (s *ScheduleService) GetSchedule(group, subgroup, englishGroup, profileSubgroup, start, end string) ([]domain.ScheduleEvent, error) {
	req := domain.ScheduleRequest{
		DStart: start, DEnd: end, Group: group, Subgroup: "*",
	}
	events, err := s.portal.FetchSchedule(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedule: %w", err)
	}

	result := filterEventsForSelection(events, subgroup, englishGroup, profileSubgroup)

	if subgroup != "" && subgroup != "*" {
		for i := range result {
			if len(result[i].SubGroup) == 1 {
				sg := result[i].SubGroup[0]
				result[i].Title = sg.STitle

				if result[i].Topic == "" {
					result[i].Topic = sg.STopic
				}
				if result[i].Room == "" {
					result[i].Room = sg.SGCaID
				}
				result[i].SubGroup = nil
			}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Day != result[j].Day {
			return result[i].Day < result[j].Day
		}
		return result[i].Start < result[j].Start
	})

	return result, nil
}

func filterEventsForSelection(events []domain.ScheduleEvent, subgroup, englishGroup, profileSubgroup string) []domain.ScheduleEvent {
	if subgroup == "" || subgroup == "*" {
		return events
	}

	out := make([]domain.ScheduleEvent, 0, len(events))

	for _, ev := range events {
		if len(ev.SubGroup) == 0 {
			out = append(out, ev)
			continue
		}

		filtered := ev.SubGroup[:0]
		for _, sg := range ev.SubGroup {
			if strings.EqualFold(sg.SGrID, "ФизраКол") || strings.EqualFold(sg.SGrID, "БрайтФит") || strings.EqualFold(sg.SGrID, "БаскетКол") {
				filtered = append(filtered, sg)
				continue
			}
			if strings.EqualFold(sg.SGrID, subgroup) {
				filtered = append(filtered, sg)
				continue
			}
			if englishRe.MatchString(sg.SGrID) {
				if englishGroup == "*" || englishGroup == "" {
					filtered = append(filtered, sg)
				} else {
					if strings.EqualFold(sg.SGrID, englishGroup) {
						filtered = append(filtered, sg)
					}
				}
				continue
			}
			if strings.HasPrefix(sg.SGrID, "Подгр") {
				var mainSubgroup string
				if strings.HasPrefix(subgroup, "Подгр") {
					mainSubgroup = subgroup
				} else {
					mainSubgroup = profileSubgroup
				}

				if mainSubgroup == "" || mainSubgroup == "*" || strings.EqualFold(mainSubgroup, "Все") {
					filtered = append(filtered, sg)
				} else {
					if strings.EqualFold(sg.SGrID, mainSubgroup) {
						filtered = append(filtered, sg)
					}
				}
				continue
			}
		}

		if len(filtered) == 0 {
			continue
		}

		ev.SubGroup = filtered
		out = append(out, ev)
	}

	return out
}

func (s *ScheduleService) GetClassDetails(clid string) (map[string]any, error) {
	return s.portal.FetchClassDetails(clid)
}
