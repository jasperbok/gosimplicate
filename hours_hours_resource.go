package gosimplicate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"
)

// HoursHoursResourceAPIResponse represents a response from the /hours/hours resource.
//
// The actual API response also contains the 'debug' fields, but I've never seen
// it with any other value than `null`, so I cannot define it here yet.
type HoursHoursResourceAPIResponse struct {
	Data   []Hours
	Errors []SimplicateAPIError `json:"errors"`
}

// Hours returns all projects that are available between start and end.
func (c *HoursClient) Hours(start, end time.Time) ([]Hours, error) {
	apiResponse := HoursHoursResourceAPIResponse{}

	uri := fmt.Sprintf("https://%s.simplicate.nl/api/v2/hours/hours", c.Client.Domain)

	query := url.Values{}
	query.Add("q[employee.id]", c.Client.Employee.EmployeeID)
	query.Add("q[start_date][ge]", start.Format("2006-01-02 15:04:05"))
	query.Add("q[start_date][lt]", end.Format("2006-01-02 15:04:05"))
	query.Add(
		"select",
		"id,start_date,end_date,project.,organization.,person.,projectservice.,type.,hours,note,is_time_defined,is_recurring,recurrence,recurrence.id,recurrence.rrule,locked,corrections,leave_id,leave_status.,absence_id,assignment_id,address.id,should_sync_to_cronofy,custom_fields.external_url",
	)
	query.Add("limit", "100")
	query.Add("metadata", "count")

	uri = fmt.Sprintf("%s?%s", uri, query.Encode())

	resp, err := c.Client.client.Get(uri)
	if err != nil {
		return apiResponse.Data, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return apiResponse.Data, err
	}

	err = json.Unmarshal(body, &apiResponse)

	if len(apiResponse.Errors) > 0 {
		return apiResponse.Data, apiResponse.Errors[0]
	}

	return apiResponse.Data, err
}
