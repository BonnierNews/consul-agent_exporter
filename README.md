# Consul agent exporter

Prometheus exporter for [Consul](https://consul.io) agents. 

## Rationale
This exporter gets data from a single Consul agent. It differs from the [prometheus/consul_exporter](https://github.com/prometheus/consul_exporter) which gets metrics from the cluster catalog.

The agent metrics can both be more finely grained than cluster level metrics, as well as showing potential divergences between nodes. See [here](metrics.md) for a list of the metrics that are exported.

## Usage

    go get -u github.com/BonnierNews/consul-agent_exporter
    cd $env:GOPATH/src/github.com/BonnierNews/consul-agent_exporter
    .\consul-agent_exporter

The prometheus metrics will be exposed on [localhost:9217](http://localhost:9217)

## Development

This package uses [`dep`](https://github.com/golang/dep) for dependency management.

## License

Under [MIT](LICENSE)
