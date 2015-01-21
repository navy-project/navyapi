package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/navy-project/navyapi/api"
)

type Client struct {
	BaseAddr string
}

func NewClient(address string) *Client {
	client := &Client{address}
	return client
}

func (c *Client) delete(path string) error {
	url := c.BaseAddr + path

	//fmt.Println("Requesting: DELETE " + url)

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return errors.New("Error - " + resp.Status)
	}

	return nil
}

func (c *Client) postJSON(path string, object interface{}) error {
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}

	body := strings.NewReader(string(b))
	url := c.BaseAddr + path

	//fmt.Println("Requesting: POST " + url)

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return errors.New("Error - " + resp.Status)
	}

	return nil
}

func (c *Client) CreateConvoy(name, manifest string) error {
	command := api.ConvoyRequest{name, manifest}
	return c.postJSON("/convoys", command)
}

func (c *Client) DeleteConvoy(name string) error {
	return c.delete("/convoys/" + name)
}
