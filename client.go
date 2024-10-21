package gosimplicate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Client struct {
	Username string
	Domain   string
	password string
	client   *http.Client
	Users    *UsersClient
	Hours    *HoursClient
}

func NewClient(username, password, domain string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	httpClient := http.Client{Jar: jar}
	c := Client{Username: username, Domain: domain, password: password, client: &httpClient}
	c.Users = &UsersClient{Client: &c}
	c.Hours = &HoursClient{Client: &c}

	err = c.Authenticate()

	return &c, err
}

func (c *Client) Authenticate() error {
	uri := fmt.Sprintf("https://%s.simplicate.nl/site/login", c.Domain)
	values := url.Values{}
	values.Add("LoginForm[username]", c.Username)
	values.Add("LoginForm[password]", c.password)

	resp, err := c.client.PostForm(uri, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) GetRegisteredHours(employeeId string, start, end time.Time) ([]Registration, error) {
	registrations := []Registration{}

	uri := fmt.Sprintf("https://%s.simplicate.nl/api/v2/hours/hours", c.Domain)

	query := url.Values{}
	query.Add("q[start_date][ge]", start.Format("2006-01-02 15:04:05"))
	query.Add("q[start_date][lt]", end.Format("2006-01-02 15:04:05"))
	query.Add("q[employee.id]", employeeId)
	query.Add(
		"select",
		"id,start_date,end_date,project.,organization.,person.,projectservice.,type.,hours,note,is_time_defined,is_recurring,recurrence,recurrence.id,recurrence.rrule,locked,corrections,leave_id,leave_status.,absence_id,assignment_id,address.id,should_sync_to_cronofy,custom_fields.external_url",
	)
	query.Add("limit", "100")
	query.Add("metadata", "count")

	uri = fmt.Sprintf("%s?%s", uri, query.Encode())

	resp, err := c.client.Get(uri)
	if err != nil {
		return registrations, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return registrations, err
	}

	responseStruct := struct {
		Data []Registration `json:"data"`
	}{}

	if err = json.Unmarshal(body, &responseStruct); err != nil {
		return registrations, err
	}

	return responseStruct.Data, nil
}

// GetActiveProjects retrieves all projects that were active between start and end.
func (c *Client) GetActiveProjects(employeeId string, start, end time.Time) ([]Project, error) {
	projects := []Project{}
	pageLimit := 100
	pageOffset := 0

	for {
		result, err := c.GetActiveProjectsPage(employeeId, start, end, pageLimit, pageOffset)
		if err != nil {
			return projects, err
		}
		projects = append(projects, result.Data...)

		// The API response does not explicitely tell us whether there are
		// more results, so we have to guess by comparing the number of
		// results to the page size limit.
		if len(result.Data) < pageLimit {
			break
		}

		pageOffset += 100
	}

	return projects, nil
}

// GetActiveProjectsPage retrieves a single 'page' of projects.
//
// Manual testing seems to indicate that the `limit` parameter has a
// maximum value of 100. Anything above that will cap the number of
// results to 100.
func (c *Client) GetActiveProjectsPage(employeeID string, start, end time.Time, pageLimit, pageOffset int) (ProjectsResponse, error) {
	response := ProjectsResponse{}

	uri := fmt.Sprintf("https://%s.simplicate.nl/api/v2/hours/projects", c.Domain)

	query := url.Values{}
	query.Add("q[employee_id]", employeeID)
	query.Add("q[start_date]", start.Format("2006-01-02 15:04:05"))
	query.Add("q[end_date]", end.Format("2006-01-02 15:04:05"))
	query.Add("limit", fmt.Sprint(pageLimit))
	query.Add("offset", fmt.Sprint(pageOffset))
	query.Add("sort", "project_name")

	uri = fmt.Sprintf("%s?%s", uri, query.Encode())

	fmt.Printf("%s\n", uri)

	resp, err := c.client.Get(uri)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}

	return response, nil
}
