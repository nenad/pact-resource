package broker

import (
	"time"
)

type (
	halPact struct {
		Embedded struct {
			Pacts []pact `json:"pacts"`
		} `json:"_embedded"`
		Links struct {
			ConsumerOldMeta struct {
				Name string `json:"name"`
			} `json:"pb:consumer"`
			ProviderOldMeta struct {
				Name string `json:"name"`
			} `json:"pb:provider"`
			ConsumerMeta struct {
				Name string `json:"name"`
			} `json:"consumer"`
			ProviderMeta struct {
				Name string `json:"name"`
			} `json:"provider"`
		} `json:"_links"`
	}

	pact struct {
		Embedded struct {
			ConsumerVersion struct {
				Number string `json:"number"`
			} `json:"consumerVersion"`
		} `json:"_embedded"`
	}

	pactDetails struct {
		CreatedAt time.Time `json:"createdAt"`
	}

	PactVersion struct {
		Provider        string
		Consumer        string
		ConsumerVersion string
	}

	Pact struct {
		PactVersion
		UpdatedAt time.Time
	}
)

func (h *halPact) ToVersions() []PactVersion {
	var pacts []PactVersion
	consumer := h.Links.ConsumerMeta.Name
	if consumer == "" {
		consumer = h.Links.ConsumerOldMeta.Name
	}
	provider := h.Links.ProviderMeta.Name
	if provider == "" {
		provider = h.Links.ProviderOldMeta.Name
	}
	for _, v := range h.Embedded.Pacts {
		var p PactVersion
		p.ConsumerVersion = v.Embedded.ConsumerVersion.Number
		p.Consumer = consumer
		p.Provider = provider
		pacts = append(pacts, p)
	}

	return pacts
}
