package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/nenad/pact-resource/broker"
)

type (
	Version struct {
		Consumer  string    `json:"consumer"`
		UpdatedAt time.Time `json:"updated_at"`
		Version   string    `json:"version"`
	}

	Source struct {
		BrokerURL string   `json:"broker_url"`
		Provider  string   `json:"provider"`
		Consumers []string `json:"consumers"`
		Tag       *string  `json:"tag"`
		Username  *string  `json:"username"`
		Password  *string  `json:"password"`
	}

	CheckRequest struct {
		Source  Source  `json:"source"`
		Version Version `json:"version"`
	}
)

func main() {
	var request CheckRequest
	populateRequest(&request)

	client := broker.NewClient(request.Source.BrokerURL)

	if request.Source.Username != nil && request.Source.Password != nil {
		broker.WithBasicAuth(*request.Source.Username, *request.Source.Password)(client)
	}

	if request.Source.Tag == nil {
		empty := ""
		request.Source.Tag = &empty
	}

	var consumerUpdates []Version
	for _, c := range request.Source.Consumers {
		versions, err := client.GetVersions(request.Source.Provider, c, *request.Source.Tag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not get pacts: %s", err)
			os.Exit(1)
		}
		for _, p := range versions {
			// TODO Add verification check as a filter
			pact, err := client.GetDetails(p.Provider, p.Consumer, p.ConsumerVersion)
			if err != nil {
				fmt.Printf("could not get details: %s", err)
				os.Exit(1)
			}

			consumerUpdates = append(consumerUpdates, Version{
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
		fmt.Printf("error while encoding response: %s", err)
		os.Exit(1)
	}
}

func populateRequest(req *CheckRequest) {
	if err := json.NewDecoder(os.Stdin).Decode(req); err != nil {
		log.Fatalf("Could not decode request: %s", err)
	}
}
