# Pact Resource

Tracks the [pacts](https://docs.pact.io/) published to a [pact-broker](https://docs.pact.io/pact_broker). 


## Source Configuration

* `broker_url`: *Required.* The path of the hosted pact-broker.

* `provider`: *Required.* The provider to track consumer pacts against.

* `consumers`: *Required.* List of consumers, along with `provider`, to track pacts for.

* `tag`: *Optional.* If specified, pulls back only pacts with that [tag](https://docs.pact.io/pact_broker/advanced_topics/using_tags).

* `username`: *Optional.* If specified, along with `password`, used to provide basic auth to the pact-broker.

* `password`: *Optional.* The password to access pact-broker.

### Example

``` yaml
resource_types:
- name: pact-resource
  type: registry-image
  source:
    repository: nenaddev/pact-resource
    tag: latest

resources:
- name: pact
  type: pact-resource
  source:
    broker_url: https://path-to.your.pact-broker.io
    provider: provider-name
    consumers:
      - consumer-1
      - consumer-2
    tag: dev
    username: ((you-are.using-a))
    password: ((secret-manager.right))
```
