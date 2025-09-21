package domain

type SubGroup struct {
	SClID  string `json:"SClID"`
	SGrID  string `json:"SGrID"`
	SGCaID string `json:"SGCaID"`
	STopic string `json:"STopic"`
	STitle string `json:"STitle"`
}

type ScheduleEvent struct {
	ClID     string     `json:"ClID"`
	Type     string     `json:"type,omitempty"`
	Day      string     `json:"Day"`
	Group    string     `json:"group"`
	Topic    string     `json:"topic"`
	Start    string     `json:"start"`
	End      string     `json:"end"`
	Room     string     `json:"room"`
	Color    string     `json:"color"`
	Title    string     `json:"title"`
	SubGroup []SubGroup `json:"SubGroup,omitempty"`
}

type ScheduleRequest struct {
	DStart   string `json:"d_start"`
	DEnd     string `json:"d_end"`
	Group    string `json:"group"`
	Subgroup string `json:"subgroup"`
}

type ScheduleResponse struct {
	Events []ScheduleEvent `json:"events"`
}
