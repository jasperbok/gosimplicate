package gosimplicate

import (
	"encoding/json"
	"fmt"
	"io"
)

// UsersUserAPIResponse represents a response from the /users/user resource.
//
// The actual API response also contains the 'debug' fields, but I've never seen
// it with any other value than `null`, so I cannot define it here yet.
type UsersUserAPIResponse struct {
	Data   Person               `json:"data"`
	Errors []SimplicateAPIError `json:"errors"`
}

func (c *UsersClient) User() (Person, error) {
	apiResponse := UsersUserAPIResponse{}

	uri := fmt.Sprintf("https://%s.simplicate.nl/api/v2/users/user", c.Client.Domain)

	resp, err := c.Client.client.Get(uri)
	if err != nil {
		return Person{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Person{}, err
	}

	err = json.Unmarshal(body, &apiResponse)

	if len(apiResponse.Errors) > 0 {
		return apiResponse.Data, apiResponse.Errors[0]
	}

	return apiResponse.Data, err
}
