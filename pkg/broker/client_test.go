package broker

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func loadFixture(t *testing.T, name string) []byte {
	bytes, err := ioutil.ReadFile("./testdata/" + name)
	if err != nil {
		t.Fatalf("Failed to load fixture %s", name)
	}
	return bytes
}

func TestGetVersionsWithoutTag(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/pacts/provider/PROVIDER/consumer/CONSUMER/versions", r.URL.Path)
		assert.Equal(t, r.Method, "GET")
		w.Header().Add("Content-Type", "application/json")
		w.Write(loadFixture(t, "versions.json"))
	}))

	client := NewClient(ts.URL)
	got, err := client.GetVersions("PROVIDER", "CONSUMER")
	assert.NoError(t, err)

	want := []PactVersion{
		{
			Provider:        "PROVIDER",
			Consumer:        "CONSUMER",
			ConsumerVersion: "5556b8149bf8bac76bc30f50a8a2dd4c22c85f30",
		},
	}

	assert.Equal(t, want, got)

}

func TestGetVersionsWithTag(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/pacts/provider/PROVIDER/consumer/CONSUMER/tag/TAG", r.URL.Path)
		assert.Equal(t, r.Method, "GET")
		w.Header().Add("Content-Type", "application/json")
		w.Write(loadFixture(t, "versions_with_tag.json"))
	}))

	client := NewClient(ts.URL)
	got, err := client.GetTaggedVersions("PROVIDER", "CONSUMER", "TAG")
	assert.NoError(t, err)

	want := []PactVersion{
		{
			Provider:        "PROVIDER",
			Consumer:        "CONSUMER",
			ConsumerVersion: "e15da45d3943bf10793a6d04cfb9f5dabe430fe2",
		},
	}

	assert.Equal(t, want, got)

}

func TestGetDetails(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/pacts/provider/PROVIDER/consumer/CONSUMER/version/VERSION", r.URL.Path)
		assert.Equal(t, r.Method, "GET")
		w.Header().Add("Content-Type", "application/json")
		w.Write(loadFixture(t, "pact_details.json"))
	}))

	client := NewClient(ts.URL)
	got, err := client.GetDetails("PROVIDER", "CONSUMER", "VERSION")
	assert.NoError(t, err)

	wantTime, err := time.Parse(time.RFC3339, "2020-03-12T07:03:04+00:00")
	if err != nil {
		t.Fatalf("failed to parse time: %s", err)
	}

	want := Pact{
		PactVersion: PactVersion{
			Provider:        "PROVIDER",
			Consumer:        "CONSUMER",
			ConsumerVersion: "VERSION",
		},
		UpdatedAt: wantTime,
	}

	assert.Equal(t, want, got)
}
