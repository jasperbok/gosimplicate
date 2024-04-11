package gosimplicate

import (
	"strings"
	"time"
)

type Employee struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Organization struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	RelationNumber string `json:"relation_number"`
}

type Project struct {
	Id            string       `json:"id"`
	Name          string       `json:"name"`
	ProjectNumber string       `json:"project_number"`
	Organization  Organization `json:"organization"`
}

type ProjectService struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	DefaultServiceId string `json:"default_service_id"`
	RevenueGroupId   string `json:"revenue_group_id"`
}

type ServiceType struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Color string `json:"color"`
	Type  string `json:"type"`
}

type Registration struct {
	Id                string         `json:"id"`
	Source            string         `json:"source"`
	StartDate         SimplicateTime `json:"start_date"`
	EndDate           SimplicateTime `json:"end_date"`
	Hours             float32        `json:"hours"`
	DurationInMinutes int            `json:"duration_in_minutes"`
	IsRecurring       bool           `json:"is_recurring"`
	IsBillable        bool           `json:"billable"`
	IsLocked          bool           `json:"locked"`
	IsProductive      bool           `json:"is_productive"`
	Employee          Employee       `json:"employee"`
	Project           Project        `json:"project"`
	ProjectService    ProjectService `json:"projectservice"`
	Type              ServiceType    `json:"type"`
}

type SimplicateTime struct {
	time.Time
}

const SimplicateTimeLayout = "2006-01-02 15:04:05"

func (st *SimplicateTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		st.Time = time.Time{}
		return
	}
	st.Time, err = time.Parse(SimplicateTimeLayout, s)
	return
}
