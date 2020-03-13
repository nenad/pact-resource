package broker_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/nenad/pact-resource/pkg/broker"
	"github.com/nenad/pact-resource/pkg/broker/test"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_ClientMethods(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	url := os.Getenv("TEST_PACT_BROKER_URL")
	tb := test.NewTestBroker(url)

	tests := []struct {
		scenario string
		test     func(test.Broker) func(*testing.T)
	}{
		{
			scenario: "Test GetDetails returns single details",
			test:     GetDetails,
		},
		{
			scenario: "Test GetVersions returns all versions",
			test:     GetVersions,
		},
	}

	for _, tt := range tests {
		if err := tb.Reset(); err != nil {
			t.Fatalf("could not reset broker: %s", err)
		}
		t.Run(tt.scenario, tt.test(tb))
	}
}

func GetDetails(tb test.Broker) func(*testing.T) {
	return func(t *testing.T) {
		if err := tb.CreatePact("PROVIDER", "CONSUMER", "VERSION"); err != nil {
			t.Fatalf("could not create pact: %s", err)
		}

		b := broker.NewClient(tb.URL)
		got, err := b.GetDetails("PROVIDER", "CONSUMER", "VERSION")
		if err != nil {
			t.Fatalf("could not get details: %s", err)
		}

		want := broker.PactVersion{
			Provider:        "PROVIDER",
			Consumer:        "CONSUMER",
			ConsumerVersion: "VERSION",
		}

		assert.Equal(t, want, got.PactVersion)
	}
}

func GetVersions(tb test.Broker) func(*testing.T) {
	return func(t *testing.T) {
		for i := 0; i < 10; i++ {
			if err := tb.CreatePact("PROVIDER", "CONSUMER", fmt.Sprintf("%d", i)); err != nil {
				t.Fatalf("could not create pact: %s", err)
			}
		}

		b := broker.NewClient(tb.URL)
		got, err := b.GetVersions("PROVIDER", "CONSUMER")
		if err != nil {
			t.Fatalf("could not get details: %s", err)
		}

		var want []broker.PactVersion
		// Pacts are sorted by creation time
		for i := 9; i >= 0; i-- {
			want = append(want, broker.PactVersion{
				Provider:        "PROVIDER",
				Consumer:        "CONSUMER",
				ConsumerVersion: fmt.Sprintf("%d", i),
			})
		}

		assert.Equal(t, want, got)
	}
}
