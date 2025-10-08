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
var profileSet = map[string]struct{}{
	"BE": {}, "FE": {}, "GD": {}, "PM": {}, "SA": {}, "CD": {},
}

func (s *ScheduleService) GetSchedule(group, subgroup, englishGroup, start, end string) ([]domain.ScheduleEvent, error) {
	commonReq := domain.ScheduleRequest{
		DStart: start, DEnd: end, Group: group, Subgroup: "*",
	}
	commonEvents, err := s.portal.FetchSchedule(commonReq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch common schedule: %w", err)
	}

	var subgroupEvents []domain.ScheduleEvent
	if subgroup != "" && subgroup != "*" {
		subReq := domain.ScheduleRequest{
			DStart: start, DEnd: end, Group: group, Subgroup: subgroup,
		}
		subgroupEvents, err = s.portal.FetchSchedule(subReq)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch subgroup schedule: %w", err)
		}
	}

	commonFiltered := filterEventsForSelection(commonEvents, subgroup, englishGroup)
	subFiltered := filterEventsForSelection(subgroupEvents, subgroup, englishGroup)

	result := mergeByClIDPreferCommon(commonFiltered, subFiltered)

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

func filterEventsForSelection(events []domain.ScheduleEvent, subgroup, englishGroup string) []domain.ScheduleEvent {
	if subgroup == "" || subgroup == "*" {
		return events
	}

	_, subgroupIsProfile := profileSet[strings.ToUpper(subgroup)]
	out := make([]domain.ScheduleEvent, 0, len(events))

	for _, ev := range events {
		if len(ev.SubGroup) == 0 {
			out = append(out, ev)
			continue
		}

		filtered := ev.SubGroup[:0]
		for _, sg := range ev.SubGroup {
			if strings.EqualFold(sg.SGrID, subgroup) {
				filtered = append(filtered, sg)
				continue
			}
			if englishRe.MatchString(sg.SGrID) {
				if englishGroup != "" && englishGroup != "*" {
					if strings.EqualFold(sg.SGrID, englishGroup) {
						filtered = append(filtered, sg)
					}
				} else if subgroupIsProfile {
					filtered = append(filtered, sg)
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

func mergeByClIDPreferCommon(common, sub []domain.ScheduleEvent) []domain.ScheduleEvent {
	byClID := make(map[string]domain.ScheduleEvent, len(common))
	result := make([]domain.ScheduleEvent, 0, len(common)+len(sub))

	for _, ev := range common {
		if _, ok := byClID[ev.ClID]; !ok {
			byClID[ev.ClID] = ev
			result = append(result, ev)
		}
	}
	for _, ev := range sub {
		if _, ok := byClID[ev.ClID]; ok {
			continue
		}
		byClID[ev.ClID] = ev
		result = append(result, ev)
	}
	return result
}

func (s *ScheduleService) GetClassDetails(clid string) (map[string]any, error) {
	return s.portal.FetchClassDetails(clid)
}
