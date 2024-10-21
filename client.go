package gosimplicate

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Client struct {
	Username string
	Domain   string
	password string
	Employee Person
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

	employee, err := c.Users.User()
	if err != nil {
		return nil, err
	}
	c.Employee = employee

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
