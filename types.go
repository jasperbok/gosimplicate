package gosimplicate

import (
	"strings"
	"time"
)

type CustomField struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Label        string `json:"label"`
	RenderType   string `json:"render_type"`
	Position     int    `json:"position"`
	IsFilterable bool   `json:"filterable"`
	IsSearchable bool   `json:"searchable"`
	IsMandatory  bool   `json:"mandatory"`
	ValueType    string `json:"value_type"`
	// Couldn't find an example of Options yet, so don't know what
	// the implementation should be like.
	//Options [] `json:"options"`
}

type Correction struct {
	Amount             int    `json:"amount"`
	Value              int    `json:"value"`
	LastCorrectionDate string `json:"last_correction_date"`
}

type WrappedBool struct {
	Value bool `json:"value"`
}

type Employee struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Organization struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	RelationNumber string `json:"relation_number"`
}

// Person represents a person in the Simplicate system.
//
// The actual API response also contains the 'timezone', 'country' and 'rights'
// fields, but these are not implemented in this struct yet.
type Person struct {
	Username                     string         `json:"username"`
	Email                        string         `json:"email"`
	Gender                       string         `json:"gender"`
	BirthDate                    SimplicateDate `json:"birth_date"`
	Initials                     string         `json:"initials"`
	FirstName                    string         `json:"first_name"`
	FamilyNamePrefix             string         `json:"family_name_prefix"`
	FamilyName                   string         `json:"family_name"`
	IsAuthyEnabled               bool           `json:"is_authy_enabled"`
	IsEmployee                   bool           `json:"is_employee"`
	IsLightUser                  bool           `json:"is_light_user"`
	EmployeeID                   string         `json:"employee_id"`
	PersonID                     string         `json:"person_id"`
	IsBlocked                    bool           `json:"is_blocked"`
	IsLockNav                    bool           `json:"is_lock_nav"`
	KeyIdentifier                string         `json:"key_identifier"`
	HoursViewMode                string         `json:"hours_view_mode"`
	IsAccountOwner               bool           `json:"is_account_owner"`
	HasExternalAgendaIntegration bool           `json:"has_external_agenda_integration"`
}

type Project struct {
	Id                        string       `json:"id"`
	Name                      string       `json:"project_name"`
	ProjectNumber             string       `json:"project_number"`
	Organization              Organization `json:"organization"`
	HasRegisterMileageEnabled bool         `json:"has_register_mileage_enabled"`
}

type ProjectService struct {
	Id               string         `json:"id"`
	Name             string         `json:"name"`
	StartDate        SimplicateDate `json:"start_date"`
	DefaultServiceId string         `json:"default_service_id"`
	RevenueGroupId   string         `json:"revenue_group_id"`
}

type ServiceType struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Tariff string `json:"tariff"`
	Color  string `json:"color"`
	Type   string `json:"type"`
}

// Hours represents a single registration of work.
type Hours struct {
	Id                  string         `json:"id"`
	Note                string         `json:"note"`
	Source              string         `json:"source"`
	StartDate           SimplicateTime `json:"start_date"`
	EndDate             SimplicateTime `json:"end_date"`
	Hours               float32        `json:"hours"`
	DurationInMinutes   int            `json:"duration_in_minutes"`
	IsRecurring         bool           `json:"is_recurring"`
	IsTimeDefined       bool           `json:"is_time_defined"`
	IsBillable          bool           `json:"billable"`
	IsLocked            bool           `json:"locked"`
	IsEditable          WrappedBool    `json:"is_editable"`
	IsDeletable         WrappedBool    `json:"is_deletable"`
	IsProductive        bool           `json:"is_productive"`
	ShouldSyncToCronofy bool           `json:"should_sync_to_cronofy"`
	Employee            Employee       `json:"employee"`
	Project             Project        `json:"project"`
	ProjectService      ProjectService `json:"projectservice"`
	Type                ServiceType    `json:"type"`
	CustomFields        []CustomField  `json:"custom_fields"`
	Corrections         Correction     `json:"corrections,omitempty"`
}

type HoursType struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Color string `json:"color"`
}

// ProjectsResponse contains response data from the /projects resource.
//
// The response also contains `errors` and `debug` keys, but I've
// never seen them containing anything but `null`, so I don't know
// what data type they hold.
type ProjectsResponse struct {
	Data []Project `json:"data"`
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

type SimplicateDate struct {
	time.Time
}

const SimplicateDateLayout = "2006-01-02"

func (st *SimplicateDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		st.Time = time.Time{}
		return
	}
	st.Time, err = time.Parse(SimplicateDateLayout, s)
	return
}
