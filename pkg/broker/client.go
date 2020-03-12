package broker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	Client struct {
		baseURL  string
		client   *http.Client
		token    string
		username string
		password string
	}

	Option func(*Client)
)

func WithBasicAuth(username, password string) Option {
	return func(broker *Client) {
		broker.username = username
		broker.password = password
	}
}

func WithClient(client *http.Client) Option {
	return func(broker *Client) {
		broker.client = client
	}
}

func NewClient(brokerURL string, opts ...Option) *Client {
	broker := Client{baseURL: brokerURL}
	for _, o := range opts {
		o(&broker)
	}

	if broker.client == nil {
		broker.client = &http.Client{Timeout: time.Second * 5}
	}

	return &broker
}

func (c *Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if len(c.username) > 0 && len(c.password) > 0 {
		auth := base64.StdEncoding.EncodeToString([]byte(c.username + ":" + c.password))
		req.Header.Add("Authorization", "Basic " + auth)
	}

	return c.client.Do(req)
}

func (c *Client) GetVersions(provider, consumer, tag string) ([]PactVersion, error) {
	url := fmt.Sprintf("%s/pacts/provider/%s/consumer/%s", c.baseURL, provider, consumer)
	if tag != "" {
		url = fmt.Sprintf("%s/tag/%s", url, tag)
	}

	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return nil, fmt.Errorf("error while requesting information: %d", resp.StatusCode)
	}

	var halPacts halPact
	err = json.NewDecoder(resp.Body).Decode(&halPacts)
	if err != nil {
		return nil, err
	}

	return halPacts.ToVersions(), nil
}

func (c *Client) GetDetails(provider, consumer, version string) (Pact, error) {
	url := fmt.Sprintf("%s/pacts/provider/%s/consumer/%s/version/%s", c.baseURL, provider, consumer, version)

	resp, err := c.Get(url)
	if err != nil {
		return Pact{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return Pact{}, fmt.Errorf("error while requesting information: %d", resp.StatusCode)
	}

	var pact pactDetails
	err = json.NewDecoder(resp.Body).Decode(&pact)
	if err != nil {
		return Pact{}, err
	}

	return Pact{
		PactVersion: PactVersion{
			Provider:        provider,
			Consumer:        consumer,
			ConsumerVersion: version,
		},
		UpdatedAt: pact.CreatedAt,
	}, nil
}

func (c *Client) GetDetailsRaw(provider, consumer, version string) ([]byte, error) {
	url := fmt.Sprintf("%s/pacts/provider/%s/consumer/%s/version/%s", c.baseURL, provider, consumer, version)

	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return nil, fmt.Errorf("error while requesting information: %d", resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
