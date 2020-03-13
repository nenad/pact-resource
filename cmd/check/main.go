package main

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/nenad/pact-resource/pkg/broker"
	"github.com/nenad/pact-resource/pkg/concourse"
)

func main() {
	var request concourse.CheckRequest
	populateRequest(&request)

	client := broker.NewClient(request.Source.BrokerURL)

	if request.Source.Username != nil && request.Source.Password != nil {
		broker.WithBasicAuth(*request.Source.Username, *request.Source.Password)(client)
	}

	var consumerUpdates []concourse.Version
	for _, c := range request.Source.Consumers {
		var versions []broker.PactVersion
		var err error
		if request.Source.Tag == nil || *request.Source.Tag == "" {
			versions, err = client.GetVersions(request.Source.Provider, c)
		} else {
			versions, err = client.GetTaggedVersions(request.Source.Provider, c, *request.Source.Tag)
		}
		if err != nil {
			concourse.FailTask("could not get pacts: %s", err)
		}
		for _, p := range versions {
			// TODO Add verification check as a filter
			pact, err := client.GetDetails(p.Provider, p.Consumer, p.ConsumerVersion)
			if err != nil {
				concourse.FailTask("could not get details: %s", err)
			}

			consumerUpdates = append(consumerUpdates, concourse.Version{
				Consumer:  c,
				UpdatedAt: pact.UpdatedAt,
				Version:   pact.PactVersion.ConsumerVersion,
			})
		}
	}

	sort.SliceStable(consumerUpdates, func(i, j int) bool {
		return consumerUpdates[i].UpdatedAt.Before(consumerUpdates[j].UpdatedAt)
	})

	if err := json.NewEncoder(os.Stdout).Encode(consumerUpdates); err != nil {
		concourse.FailTask("error while encoding response: %s", err)
	}
}

func populateRequest(req *concourse.CheckRequest) {
	if err := json.NewDecoder(os.Stdin).Decode(req); err != nil {
		concourse.FailTask("could not decode request: %s", err)
	}
}
