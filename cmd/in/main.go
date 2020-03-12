package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nenad/pact-resource/pkg/broker"
	"github.com/nenad/pact-resource/pkg/concourse"
)

func main() {
	var request concourse.InRequest
	populateRequest(&request)

	client := broker.NewClient(request.Source.BrokerURL)

	if request.Source.Username != nil && request.Source.Password != nil {
		broker.WithBasicAuth(*request.Source.Username, *request.Source.Password)(client)
	}

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "first argument must be a directory")
		os.Exit(1)
	}

	dir := os.Args[1]

	if request.Source.Tag == nil {
		empty := ""
		request.Source.Tag = &empty
	}

	bytes, err := client.GetDetailsRaw(request.Source.Provider, request.Version.Consumer, request.Version.Version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read bytes: %s", err)
		os.Exit(1)
	}

	pactPath := fmt.Sprintf("%s-%s-%s.json", request.Source.Provider, request.Version.Consumer, request.Version.Version)
	pactPath = strings.ReplaceAll(pactPath, " ", "-")

	file, err := os.Create(filepath.Join(dir, pactPath))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open file: %s", err)
		os.Exit(1)
	}

	_, err = file.Write(bytes)
	defer file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not write to file: %s", err)
		os.Exit(1)
	}

	err = json.NewEncoder(os.Stdout).Encode(concourse.InResponse{
		Version: request.Version,
		Metadata: concourse.Metadata{
			{Name: "pact", Value: pactPath},
		},
	})
	if err != nil {
		fmt.Printf("error while encoding response: %s", err)
		os.Exit(1)
	}
}

func populateRequest(req *concourse.InRequest) {
	if err := json.NewDecoder(os.Stdin).Decode(req); err != nil {
		log.Fatalf("Could not decode request: %s", err)
	}
}
