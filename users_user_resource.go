package gosimplicate

import (
	"encoding/json"
	"fmt"
	"io"
)

// UsersUserAPIResponse represents a response from the /users/user resource.
//
// The actual API response also contains the 'errors' and 'debug' fields, but
// I've never seen them with any other value than `null`, so I cannot define
// them here yet.
type UsersUserAPIResponse struct {
	Data Person `json:"data"`
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
	return apiResponse.Data, err
}
