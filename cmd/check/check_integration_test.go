package main_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/nenad/pact-resource/pkg/broker/test"
	"github.com/nenad/pact-resource/pkg/concourse"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_CheckReturnsOutputSuccessfully(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	url := os.Getenv("TEST_PACT_BROKER_URL")
	tb := test.NewTestBroker(url)

	if err := tb.Reset(); err != nil {
		t.Fatalf("could not reset the broker: %s", err)
	}
	if err := tb.CreatePact("PROVIDER", "CONSUMER", "VERSION"); err != nil {
		t.Fatalf("could not create a pact: %s", err)
	}

	cmd := exec.Command("../../bin/check")
	input, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("could not get input pipe: %s", err)
	}
	output := &bytes.Buffer{}
	errors := &bytes.Buffer{}
	cmd.Stdout = output
	cmd.Stderr = errors
	if err := cmd.Start(); err != nil {
		t.Fatalf("could not start binary, make sure it's built first: %s", err)
	}


	req := concourse.CheckRequest{
		Source: concourse.Source{
			BrokerURL: url,
			Provider:  "PROVIDER",
			Consumers: []string{"CONSUMER"},
		},
	}

	if err := json.NewEncoder(input).Encode(req); err != nil {
		t.Fatalf("could not encode request: %s", err)
	}

	if err := cmd.Wait(); err != nil {
		t.Fatalf("could not stop binary: %s", err)
	}

	want := []concourse.Version{{
		Consumer:  "CONSUMER",
		UpdatedAt: time.Now().Round(time.Hour).UTC(),
		Version:   "VERSION",
	}}

	var got []concourse.Version
	if err := json.NewDecoder(output).Decode(&got); err != nil {
		t.Fatalf("could not decode response: %s", err)
	}

	assert.Len(t, got, 1)
	got[0].UpdatedAt = got[0].UpdatedAt.Round(time.Hour).UTC()
	assert.Equal(t, want, got)
}
