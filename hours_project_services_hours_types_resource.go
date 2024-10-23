package gosimplicate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"
)

// HoursProjectServiceHoursTypesAPIResponse represents a response from the /hours/projectserviceshourstypes resource.
//
// The actual API response also contains the 'debug' fields, but I've never seen
// it with any other value than `null`, so I cannot define it here yet.
type HoursProjectServiceHoursTypesAPIResponse struct {
	Data   []HoursType          `json:"data"`
	Errors []SimplicateAPIError `json:"errors"`
}

func (c *HoursClient) ProjectServiceHoursTypes(projectId, projectServiceId string, start, end time.Time) ([]HoursType, error) {
	apiResponse := HoursProjectServiceHoursTypesAPIResponse{}

	uri := fmt.Sprintf(
		"https://%s.simplicate.nl/api/v2/hours/projectservicehourstypes", c.Client.Domain,
	)

	query := url.Values{}
	query.Add("q[employee_id]", c.Client.Employee.EmployeeID)
	query.Add("q[project_id]", projectId)
	query.Add("q[projectservice_id]", projectServiceId)
	query.Add("q[start_date]", start.Format("2006-01-02"))
	query.Add("q[end_date]", end.Format("2006-01-02"))

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
