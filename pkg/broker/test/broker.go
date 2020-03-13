package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func NewTestBroker(url string) Broker {
	if url == "" {
		panic("Test broker URL not set")
	}

	return Broker{
		URL:    url,
		client: &http.Client{Timeout: time.Second},
	}
}

type Broker struct {
	URL    string
	client *http.Client
}

func (b *Broker) Reset() error {
	pacticipants, err := b.GetPacticipants()
	if err != nil {
		return err
	}

	for _, p := range pacticipants {
		resp, err := b.do("DELETE", fmt.Sprintf("/pacticipants/%s", p), nil)
		if err != nil {
			return fmt.Errorf("could not delete pacticipant: %w", err)
		}
		if resp.StatusCode != 204 {
			return fmt.Errorf("could not delete pacticipant, got status code of %d", resp.StatusCode)
		}
	}

	return nil
}

// Returns a list of all pacticipants
func (b *Broker) GetPacticipants() (pacticipants []string, err error) {
	resp, err := b.do("GET", "/pacticipants", nil)
	if err != nil {
		return nil, fmt.Errorf("could not get pacticipants: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("could not get pacticipants, status code %d", resp.StatusCode)
	}

	type response struct {
		Embedded struct {
			Pacticipants []struct {
				Name string `json:"name"`
			} `json:"pacticipants"`
		} `json:"_embedded"`
	}
	var jsonResponse response
	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}

	for _, p := range jsonResponse.Embedded.Pacticipants {
		pacticipants = append(pacticipants, p.Name)
	}

	return pacticipants, nil
}

// Creates a dummy pact between a consumer and provider
func (b *Broker) CreatePact(provider, consumer, version string) error {
	resp, err := b.do(
		"PUT",
		fmt.Sprintf("/pacts/provider/%s/consumer/%s/version/%s", provider, consumer, version),
		generatePact(consumer, provider),
	)
	if err != nil {
		return fmt.Errorf("could not create a pact: %w", err)
	}

	if resp.StatusCode != 201 {
		return fmt.Errorf("failed to create a pact, got error code %d", resp.StatusCode)
	}
	return nil
}

func (b *Broker) do(method string, path string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(method, fmt.Sprintf("%s%s", b.URL, path), body)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")

	return b.client.Do(r)
}

func generatePact(consumer, provider string) *bytes.Buffer {
	raw := `{ "consumer": { "name": "%s" }, "provider": { "name": "%s" },
"interactions": [ { "description" : "a request for something", "provider_state": "something exists",
"request": { "method": "get", "path" : "/something/%d" }, "response": { "status": 200, "body" : "something" } } ] }`
	return bytes.NewBufferString(fmt.Sprintf(raw, consumer, provider, rand.Int()))
}
