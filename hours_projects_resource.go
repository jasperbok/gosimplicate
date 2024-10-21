package gosimplicate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"
)

// HoursProjectsResourceAPIResponse represents a response from the /hours/projects resource.
//
// The actual API response also contains the 'debug' fields, but I've never seen
// it with any other value than `null`, so I cannot define it here yet.
type HoursProjectsResourceAPIResponse struct {
	Data   []Project
	Errors []SimplicateAPIError `json:"errors"`
}

// Projects returns all projects that are available between start and end.
func (c *HoursClient) Projects(start, end time.Time) ([]Project, error) {
	var projects []Project

	pageLimit := 100
	pageOffset := 100

	for {
		result, err := c.projectsPage(start, end, pageLimit, pageOffset)
		if err != nil {
			return projects, err
		}
		projects = append(projects, result...)

		// The API response does not explicitly tell us whether there are more
		// results, so we have to guess by comparing the number of results to
		// the page size limit.
		if len(result) < pageLimit {
			break
		}

		pageOffset += 100
	}

	return projects, nil
}

// projectsPage retrieves a single 'page' of projects.
//
// Manual testing seems to indicate that the `limit` parameter has a maximum
// value of 100. Anything above that will cap the number of results to 100.
func (c *HoursClient) projectsPage(start, end time.Time, pageLimit, pageOffset int) ([]Project, error) {
	apiResponse := HoursProjectsResourceAPIResponse{}

	uri := fmt.Sprintf("https://%s.simplicate.nl/api/v2/hours/projects", c.Client.Domain)

	query := url.Values{}
	query.Add("q[employee_id]", c.Client.Employee.EmployeeID)
	query.Add("q[start_date]", start.Format("2006-01-02 15:04:05"))
	query.Add("q[end_date]", end.Format("2006-01-02 15:04:05"))
	query.Add("limit", fmt.Sprint(pageLimit))
	query.Add("offset", fmt.Sprint(pageOffset))
	query.Add("sort", "project_name")

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
